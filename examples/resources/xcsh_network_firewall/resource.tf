# Network Firewall Resource Example
# Manages a Network Firewall resource in F5 Distributed Cloud for network firewall is created by users in system namespace. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Basic Network Firewall configuration
resource "xcsh_network_firewall" "example" {
  name      = "example-network-firewall"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # [OneOf: active_enhanced_firewall_policies, active_network...
  active_enhanced_firewall_policies {
    # Configure active_enhanced_firewall_policies settings
  }
  # Ordered List of Enhanced Firewall Policies active .
  enhanced_firewall_policies {
    # Configure enhanced_firewall_policies settings
  }
  # [OneOf: active_fast_acls, disable_fast_acl; Default: disa...
  active_fast_acls {
    # Configure active_fast_acls settings
  }
}

# The following optional fields have server-applied defaults and can be omitted:
# - disable_fast_acl
# - disable_forward_proxy_policy
# - disable_network_policy
