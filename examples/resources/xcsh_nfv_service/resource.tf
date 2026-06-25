# Nfv Service Resource Example
# Manages new NFV service with configured parameters. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Basic Nfv Service configuration
resource "xcsh_nfv_service" "example" {
  name      = "example-nfv-service"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # [OneOf: disable_https_management, https_management; Defau...
  disable_https_management {
    # Configure disable_https_management settings
  }
  # [OneOf: disable_ssh_access, enabled_ssh_access; Default: ...
  disable_ssh_access {
    # Configure disable_ssh_access settings
  }
  # Configuration parameter for enabled ssh access.
  enabled_ssh_access {
    # Configure enabled_ssh_access settings
  }
}
