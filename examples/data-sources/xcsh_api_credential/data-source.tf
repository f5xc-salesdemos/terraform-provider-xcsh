# API Credential Data Source Example
# Retrieves information about an existing API Credential

# Look up an existing API Credential by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_api_credential" "example" {
  name      = "example-api-credential"
  namespace = "system"
}

output "api_credential_id" {
  value = data.xcsh_api_credential.example.id
}
