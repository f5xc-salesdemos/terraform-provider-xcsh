# Enhanced Firewall Policy Resource Example
# Manages a Enhanced Firewall Policy resource in F5 Distributed Cloud for enhanced firewall policy specification. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Enhanced Firewall Policy configuration
resource "xcsh_enhanced_firewall_policy" "example" {
  name      = "example-enhanced-firewall-policy"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Enhanced Firewall Policy configuration
  rule_list {
    rules {
      metadata {
        name = "allow-web-traffic"
      }
      allow {}
      advanced_action {
        action = "LOG"
      }
      source_prefix_list {
        ip_prefix_set {
          name      = "trusted-ips"
          namespace = "staging"
        }
      }
      all_traffic {}
    }
  }
}

# The following optional fields have server-applied defaults and can be omitted:
# - allow_all
