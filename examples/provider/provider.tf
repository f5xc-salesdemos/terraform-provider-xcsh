terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Configure the F5XC Provider with API Token Authentication
provider "xcsh" {
  api_url   = "https://your-tenant.console.ves.volterra.io"
  api_token = var.xcsh_api_token
}

# Alternatively, use environment variables:
# export XCSH_API_URL="https://your-tenant.console.ves.volterra.io"
# export XCSH_API_TOKEN="your-api-token"

variable "xcsh_api_token" {
  description = "F5 Distributed Cloud API Token"
  type        = string
  sensitive   = true
}

# Or use P12 Certificate Authentication:
# provider "xcsh" {
#   api_url      = "https://your-tenant.console.ves.volterra.io"
#   api_p12_file = "/path/to/certificate.p12"
#   p12_password = var.xcsh_p12_password
# }
#
# Environment variables for P12 authentication:
# export XCSH_API_URL="https://your-tenant.console.ves.volterra.io"
# export XCSH_P12_FILE="/path/to/certificate.p12"
# export XCSH_P12_PASSWORD="your-p12-password"  # gitleaks:allow
