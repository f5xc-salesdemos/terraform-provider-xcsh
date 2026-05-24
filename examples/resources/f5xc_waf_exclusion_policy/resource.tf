# WAF Exclusion Policy Resource Example
# Manages WAF exclusion policy. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic WAF Exclusion Policy configuration
resource "f5xc_waf_exclusion_policy" "example" {
  name      = "example-waf-exclusion-policy"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # WAF Exclusion Rules. An ordered list of rules.
  waf_exclusion_rules {
    # Configure waf_exclusion_rules settings
  }
  # Enable this option
  any_domain {
    # Configure any_domain settings
  }
  # Enable this option
  any_path {
    # Configure any_path settings
  }
}
