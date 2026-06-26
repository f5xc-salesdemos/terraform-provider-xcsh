# F5 Distributed Cloud Provider - Authentication Example
# ======================================================
#
# This example demonstrates the three authentication methods supported
# by the F5XC Terraform provider. Uncomment the method you want to use.
#
# IMPORTANT: Never commit credentials to version control!
#
# QUICK START:
# 1. Choose your authentication method (environment variables recommended)
# 2. Set the required environment variables
# 3. Run: terraform init && terraform plan

terraform {
  required_version = ">= 1.8"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# =============================================================================
# AUTHENTICATION METHOD 1: Environment Variables (Recommended)
# =============================================================================
#
# This is the recommended approach - credentials are passed via environment
# variables and the provider block remains empty.
#
# For API Token:
#   export XCSH_API_URL="https://your-tenant.console.ves.volterra.io"
#   export XCSH_API_TOKEN="your-api-token"
#
# For P12 Certificate:
#   export XCSH_API_URL="https://your-tenant.console.ves.volterra.io"
#   export XCSH_P12_FILE="/path/to/credentials.p12"
#   export XCSH_P12_PASSWORD="your-p12-password"  # gitleaks:allow
#
# For PEM Certificate:
#   export XCSH_API_URL="https://your-tenant.console.ves.volterra.io"
#   export XCSH_CERT="/path/to/certificate.pem"
#   export XCSH_KEY="/path/to/private-key.pem"

provider "xcsh" {
  # Authentication via environment variables
  # No explicit configuration needed
}

# =============================================================================
# AUTHENTICATION METHOD 2: Provider Configuration with Variables
# =============================================================================
#
# Uncomment this block and comment out the empty provider block above if you
# prefer explicit configuration. The actual values should come from
# variables (see variables.tf) populated via terraform.tfvars or TF_VAR_
# environment variables.
#
# provider "xcsh" {
#   api_url   = var.xcsh_api_url
#   api_token = var.xcsh_api_token
# }

# =============================================================================
# AUTHENTICATION METHOD 3: P12 Certificate via Variables
# =============================================================================
#
# For P12 certificate authentication with explicit configuration:
#
# provider "xcsh" {
#   api_url      = var.xcsh_api_url
#   api_p12_file = var.xcsh_api_p12_file
#   p12_password = var.xcsh_p12_password
# }

# =============================================================================
# AUTHENTICATION METHOD 4: PEM Certificate via Variables
# =============================================================================
#
# For PEM certificate authentication (extracted from P12):
#
# provider "xcsh" {
#   api_url  = var.xcsh_api_url
#   api_cert = var.xcsh_api_cert  # gitleaks:allow
#   api_key  = var.xcsh_api_key  # gitleaks:allow
# }

# =============================================================================
# Test Resource - Validates Authentication
# =============================================================================
#
# This data source validates that authentication is working correctly.
# It retrieves information about the "system" namespace which always exists.

data "xcsh_namespace" "system" {
  name = "system"
}

# Output the namespace to confirm authentication worked
output "authentication_test" {
  description = "Authentication successful - retrieved system namespace"
  value = {
    namespace   = data.xcsh_namespace.system.name
    description = data.xcsh_namespace.system.description
  }
}
