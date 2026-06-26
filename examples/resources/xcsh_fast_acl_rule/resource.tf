# Fast ACL Rule Resource Example
# Manages new Fast ACL rule, has specification to match source IP, source port and action to apply. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Fast ACL Rule configuration
resource "xcsh_fast_acl_rule" "example" {
  name      = "example-fast-acl-rule"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Source Ports. L4 port numbers to match .
  port {
    # Configure port settings
  }
  # Enable this option
  all {
    # Configure all settings
  }
  # Enable this option
  dns {
    # Configure dns settings
  }
}
