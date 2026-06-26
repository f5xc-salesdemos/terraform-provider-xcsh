# Tpm API Key Data Source Example
# Retrieves information about an existing Tpm API Key

# Look up an existing Tpm API Key by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_tpm_api_key" "example" {
  name      = "example-tpm-api-key"
  namespace = "system"
}

output "tpm_api_key_id" {
  value = data.xcsh_tpm_api_key.example.id
}
