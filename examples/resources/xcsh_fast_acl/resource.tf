# Fast ACL Resource Example
# Manages object, object contains rules to protect site from denial of service It has destination{destination IP, destination port) and references to. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Fast ACL configuration
resource "xcsh_fast_acl" "example" {
  name      = "example-fast-acl"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Type establishes a direct reference from one object(the r...
  protocol_policer {
    # Configure protocol_policer settings
  }
  # [OneOf: re_acl, site_acl] Fast ACL for RE. Fast ACL define...
  re_acl {
    # Configure re_acl settings
  }
  # Enable this option
  all_public_vips {
    # Configure all_public_vips settings
  }
}
