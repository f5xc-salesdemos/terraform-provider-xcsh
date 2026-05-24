# Bigip Irule Resource Example
# Manages a BIG-IP Irule resource in F5 Distributed Cloud for desired state for big-ip irule service configuration.

# Basic Bigip Irule configuration
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


resource "f5xc_bigip_irule" "example" {
  name      = "example-bigip-irule"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }
}
