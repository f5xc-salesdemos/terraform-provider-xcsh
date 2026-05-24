# Network Policy View Resource Example
# Manages a Network Policy View resource in F5 Distributed Cloud for network policy view specification. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Network Policy View configuration
resource "f5xc_network_policy_view" "example" {
  name      = "example-network-policy-view"
  namespace = "system"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Ordered list of rules applied to connections from policy ...
  egress_rules {
    # Configure egress_rules settings
  }
  # Network Policy Rule Advanced Action provides additional O...
  adv_action {
    # Configure adv_action settings
  }
  # Configuration parameter for all tcp traffic.
  all_tcp_traffic {
    # Configure all_tcp_traffic settings
  }
}
