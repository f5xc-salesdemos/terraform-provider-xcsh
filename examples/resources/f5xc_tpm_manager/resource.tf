# Tpm Manager Resource Example
# Manages a Tpm Manager resource in F5 Distributed Cloud for create a tpm manager configuration.

# Basic Tpm Manager configuration
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


resource "f5xc_tpm_manager" "example" {
  name      = "example-tpm-manager"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }
}
