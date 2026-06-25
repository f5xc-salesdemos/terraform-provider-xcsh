# Network Connector Resource Example
# Manages a Network Connector resource in F5 Distributed Cloud for network connector is created by users in system namespace. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Network Connector configuration
resource "xcsh_network_connector" "example" {
  name      = "example-network-connector"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Network Connector configuration
  # Direct connection
  sli_to_global_dr {
    global_vn {
      name      = "global-network"
      namespace = "staging"
    }
  }

  # Disable forward proxy
  disable_forward_proxy {}
}
