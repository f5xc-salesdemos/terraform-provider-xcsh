# Virtual Network Resource Example
# Manages virtual network in given namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Basic Virtual Network configuration
resource "xcsh_virtual_network" "example" {
  name      = "example-virtual-network"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  legacy_type = "VIRTUAL_NETWORK_SITE_LOCAL"
}
