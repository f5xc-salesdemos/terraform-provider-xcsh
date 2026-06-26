# BGP Routing Policy Resource Example
# Manages a BGP Routing Policy resource in F5 Distributed Cloud for bgp routing policy is a list of rules containing match criteria and action to be applied. these rules help control routes which are imported or exported to bgp peers. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic BGP Routing Policy configuration
resource "xcsh_bgp_routing_policy" "example" {
  name      = "example-bgp-routing-policy"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # BGP Routing policy is composed of one or more rules. Note...
  rules {
    # Configure rules settings
  }
  # Action to be enforced if the BGP route matches the rule.
  action {
    # Configure action settings
  }
  # Enable this option
  aggregate {
    # Configure aggregate settings
  }
}
