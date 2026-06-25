# Fleet Resource Example
# Manages fleet will create a fleet object in 'system' namespace of the user. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Basic Fleet configuration
resource "xcsh_fleet" "example" {
  name      = "example-fleet"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Fleet configuration
  fleet_label = "env=production"

  # Network connectors
  inside_virtual_network {
    name      = "inside-network"
    namespace = "staging"
  }

  outside_virtual_network {
    name      = "outside-network"
    namespace = "staging"
  }

  # Default config
  default_config {}
}
