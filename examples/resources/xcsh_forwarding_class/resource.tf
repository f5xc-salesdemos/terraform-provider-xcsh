# Forwarding Class Resource Example
# Manages a Forwarding Class resource in F5 Distributed Cloud for forwarding class is created by users in system namespace. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Forwarding Class configuration
resource "xcsh_forwarding_class" "example" {
  name      = "example-forwarding-class"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # [OneOf: dscp, no_marking, tos_value; Default: no_marking]...
  dscp {
    # Configure dscp settings
  }
  # [OneOf: dscp_based_queue, queue_id_to_use] Configuration ...
  dscp_based_queue {
    # Configure dscp_based_queue settings
  }
  # Enable this option
  no_marking {
    # Configure no_marking settings
  }
}
