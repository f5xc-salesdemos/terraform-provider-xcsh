# Tpm API Key Data Source Example
# Retrieves information about an existing Tpm API Key

# Look up an existing Tpm API Key by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_tpm_api_key" "example" {
  name      = "example-tpm-api-key"
  namespace = "system"
}

output "tpm_api_key_id" {
  value = data.f5xc_tpm_api_key.example.id
}
