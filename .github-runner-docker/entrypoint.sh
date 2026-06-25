#!/bin/bash
# Entrypoint for GitHub Actions Runner Container
#
# Required environment variables:
#   GITHUB_REPOSITORY  - owner/repo (e.g., robinmordasiewicz/terraform-provider-xcsh)
#   GITHUB_TOKEN       - Personal Access Token with repo scope
#
# Optional environment variables:
#   RUNNER_NAME        - Custom runner name (default: hostname)
#   RUNNER_LABELS      - Additional labels (default: self-hosted,Linux,X64)
#   RUNNER_WORKDIR     - Work directory (default: _work)

set -euo pipefail

# Validate required environment variables
if [[ -z "${GITHUB_REPOSITORY:-}" ]]; then
  echo "ERROR: GITHUB_REPOSITORY is required (e.g., owner/repo)"
  exit 1
fi

if [[ -z "${GITHUB_TOKEN:-}" ]]; then
  echo "ERROR: GITHUB_TOKEN is required (Personal Access Token with repo scope)"
  exit 1
fi

# Set defaults
RUNNER_NAME="${RUNNER_NAME:-$(hostname)}"
RUNNER_LABELS="${RUNNER_LABELS:-self-hosted,Linux,X64,container}"
RUNNER_WORKDIR="${RUNNER_WORKDIR:-_work}"

cd /home/runner/actions-runner

echo "=============================================="
echo "  GitHub Actions Runner (Containerized)"
echo "=============================================="
echo "Repository:  ${GITHUB_REPOSITORY}"
echo "Runner Name: ${RUNNER_NAME}"
echo "Labels:      ${RUNNER_LABELS}"
echo "=============================================="

# Get registration token from GitHub API
echo "Getting registration token..."
REGISTRATION_TOKEN=$(curl -sX POST \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  -H "Accept: application/vnd.github+json" \
  "https://api.github.com/repos/${GITHUB_REPOSITORY}/actions/runners/registration-token" |
  jq -r '.token')

if [[ -z "$REGISTRATION_TOKEN" || "$REGISTRATION_TOKEN" == "null" ]]; then
  echo "ERROR: Failed to get registration token"
  echo "Make sure GITHUB_TOKEN has 'repo' scope and you have admin access"
  exit 1
fi

echo "Registration token obtained"

# Remove existing runner configuration if present
if [[ -f ".runner" ]]; then
  echo "Removing existing runner configuration..."
  ./config.sh remove --token "$REGISTRATION_TOKEN" 2>/dev/null || true
fi

# Configure the runner
echo "Configuring runner..."
./config.sh \
  --url "https://github.com/${GITHUB_REPOSITORY}" \
  --token "$REGISTRATION_TOKEN" \
  --name "$RUNNER_NAME" \
  --labels "$RUNNER_LABELS" \
  --work "$RUNNER_WORKDIR" \
  --unattended \
  --replace

# Cleanup function for graceful shutdown
cleanup() {
  echo ""
  echo "Received shutdown signal, removing runner..."

  # Get removal token
  REMOVAL_TOKEN=$(curl -sX POST \
    -H "Authorization: token ${GITHUB_TOKEN}" \
    -H "Accept: application/vnd.github+json" \
    "https://api.github.com/repos/${GITHUB_REPOSITORY}/actions/runners/remove-token" |
    jq -r '.token')

  if [[ -n "$REMOVAL_TOKEN" && "$REMOVAL_TOKEN" != "null" ]]; then
    ./config.sh remove --token "$REMOVAL_TOKEN" 2>/dev/null || true
    echo "Runner removed from repository"
  fi

  exit 0
}

# Trap signals for graceful shutdown
trap cleanup SIGTERM SIGINT SIGQUIT

echo "Starting runner..."
./run.sh &

# Wait for runner process
wait $!
