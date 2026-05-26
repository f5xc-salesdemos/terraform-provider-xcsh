#!/usr/bin/env bash
set -euo pipefail

# UAT Test Harness for AI-generated Terraform plans
# Validates golden files and compares AI-generated output against them.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROMPTS_DIR="${SCRIPT_DIR}/prompts"
GOLDEN_DIR="${SCRIPT_DIR}/golden"
RESULTS_DIR="${SCRIPT_DIR}/results"

PASS=0
FAIL=0
SKIP=0

# Ensure results dir exists
mkdir -p "${RESULTS_DIR}"

# Validate a Terraform file in a temporary directory.
# Returns 0 on success, non-zero on failure.
validate_tf() {
  local tf_file="$1"
  local tmpdir
  tmpdir="$(mktemp -d)"
  cp "${tf_file}" "${tmpdir}/main.tf"
  (
    cd "${tmpdir}"
    terraform init -backend=false -no-color >/dev/null 2>&1 &&
      terraform validate -no-color 2>&1
  )
  local rc=$?
  rm -rf "${tmpdir}"
  return $rc
}

echo "========================================================================"
echo "UAT Test Harness"
echo "========================================================================"
echo ""

for prompt_file in "${PROMPTS_DIR}"/*.txt; do
  # Extract test name from filename (strip path and .txt extension)
  test_name="$(basename "${prompt_file}" .txt)"
  golden_file="${GOLDEN_DIR}/${test_name}.golden.tf"
  result_file="${RESULTS_DIR}/${test_name}.result.tf"

  echo "------------------------------------------------------------------------"
  echo "Test: ${test_name}"

  # ── Step 1: verify golden file exists ──────────────────────────────────────
  if [[ ! -f "${golden_file}" ]]; then
    echo "  ERROR: Golden file not found: ${golden_file}"
    FAIL=$((FAIL + 1))
    continue
  fi

  # ── Step 2: validate the golden file itself (sanity check) ─────────────────
  echo "  Validating golden file..."
  if validate_tf "${golden_file}"; then
    echo "  Golden file: OK"
  else
    echo "  FAIL: Golden file failed terraform validate"
    FAIL=$((FAIL + 1))
    continue
  fi

  # ── Step 3: compare result against golden (if result exists) ───────────────
  if [[ ! -f "${result_file}" ]]; then
    echo "  SKIP: No result file found at ${result_file}"
    echo "        (Run the AI against the prompt to generate a result)"
    SKIP=$((SKIP + 1))
    continue
  fi

  echo "  Validating result file..."
  if ! validate_tf "${result_file}"; then
    echo "  FAIL: Result file failed terraform validate"
    FAIL=$((FAIL + 1))
    continue
  fi
  echo "  Result file: OK"

  echo "  Diffing result against golden..."
  if diff --unified "${golden_file}" "${result_file}"; then
    echo "  PASS: Result matches golden"
    PASS=$((PASS + 1))
  else
    echo "  FAIL: Result differs from golden (see diff above)"
    FAIL=$((FAIL + 1))
  fi
done

echo ""
echo "========================================================================"
echo "Results: PASS=${PASS}  FAIL=${FAIL}  SKIP=${SKIP}"
echo "========================================================================"

# ── Generate coverage report ──────────────────────────────────────────────────
TOTAL=$((PASS + FAIL + SKIP))
COVERAGE_REPORT="${RESULTS_DIR}/coverage-report.json"
cat >"${COVERAGE_REPORT}" <<EOF
{
  "total": ${TOTAL},
  "pass": ${PASS},
  "fail": ${FAIL},
  "skip": ${SKIP},
  "generated_at": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
}
EOF
echo ""
echo "Coverage report written to: ${COVERAGE_REPORT}"

# Exit non-zero if any test failed
if [[ "${FAIL}" -gt 0 ]]; then
  exit 1
fi
exit 0
