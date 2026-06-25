# Segment Resource Example
# Manages a Segment resource in F5 Distributed Cloud for segment. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Basic Segment configuration
resource "xcsh_segment" "example" {
  name      = "example-segment"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # [OneOf: disable, enable] Enable this option
  disable_spec {
    # Configure disable_spec settings
  }
  # Enable this option
  enable {
    # Configure enable settings
  }
}
