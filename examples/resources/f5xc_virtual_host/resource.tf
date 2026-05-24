# Virtual Host Resource Example
# Manages virtual host in a given namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Virtual Host configuration
resource "f5xc_virtual_host" "example" {
  name      = "example-virtual-host"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Advertise Policy allows you to define networks or sites w...
  advertise_policies {
    # Configure advertise_policies settings
  }
  # [OneOf: authentication, no_authentication; Default: no_au...
  authentication {
    # Configure authentication settings
  }
  # Reference to Authentication Config Object .
  auth_config {
    # Configure auth_config settings
  }
}
