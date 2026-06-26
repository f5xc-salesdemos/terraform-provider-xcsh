terraform {
  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# F5 Distributed Cloud Provider - Authentication Outputs
# =======================================================
#
# These outputs help verify that authentication is working correctly
# and provide useful information about the connection.

output "api_url" {
  description = "The F5XC API URL used for authentication (from environment or configuration)"
  value       = var.xcsh_api_url != "" ? var.xcsh_api_url : "Set via XCSH_API_URL environment variable"
}

output "authentication_method" {
  description = "The authentication method being used based on configuration"
  value       = var.xcsh_api_token != "" ? "API Token" : (var.xcsh_p12_file != "" ? "P12 Certificate" : (var.xcsh_cert != "" ? "PEM Certificate" : "Environment Variables"))
}

output "system_namespace" {
  description = "Information about the system namespace (confirms authentication is working)"
  value = {
    name        = data.xcsh_namespace.system.name
    description = data.xcsh_namespace.system.description
  }
}
