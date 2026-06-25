# OIDC Provider Data Source Example
# Retrieves information about an existing OIDC Provider

# Look up an existing OIDC Provider by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_oidc_provider" "example" {
  name      = "example-oidc-provider"
  namespace = "system"
}

output "oidc_provider_id" {
  value = data.xcsh_oidc_provider.example.id
}
