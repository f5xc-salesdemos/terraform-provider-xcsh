# Network Policy Resource Example
# Manages new network policy with configured parameters in specified namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Network Policy configuration
resource "xcsh_network_policy" "example" {
  name      = "example-network-policy"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Network Policy configuration
  endpoint {
    any {}
  }

  ingress_rules {
    metadata {
      name = "allow-http"
    }
    spec {
      action = "ALLOW"
      any    = {}
    }
  }

  egress_rules {
    metadata {
      name = "allow-all-egress"
    }
    spec {
      action = "ALLOW"
      any    = {}
    }
  }
}
