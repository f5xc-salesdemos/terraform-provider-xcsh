# CDN Purge Command Resource Example
# Manages a CDN Purge Command resource in F5 Distributed Cloud for cdn purge command specification. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source = "f5xc-salesdemos/f5xc"
    }
  }
}

# Basic CDN Purge Command configuration
resource "f5xc_cdn_purge_command" "example" {
  name      = "example-cdn-purge-command"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # [OneOf: hard_purge, soft_purge] Enable this option
  hard_purge {
    # Configure hard_purge settings
  }
  # Enable this option
  purge_all {
    # Configure purge_all settings
  }
  # Enable this option
  soft_purge {
    # Configure soft_purge settings
  }
}
