# Tpm API Key Resource Example
# Manages a Tpm API Key resource in F5 Distributed Cloud for apikey object when successfully created returns actual apikey bytes which is used by the users to call in to tpm provisioning api. configuration.

# Basic Tpm API Key configuration
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


resource "f5xc_tpm_api_key" "example" {
  name      = "example-tpm-api-key"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # TPM Category. APIKey needs a reference to an existing TPM...
  category_ref {
    # Configure category_ref settings
  }
}
