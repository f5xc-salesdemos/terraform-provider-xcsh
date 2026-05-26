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

  waf_exclusion_rules {
    // One of the arguments from this list "any_domain exact_value suffix_value" must be set

    any_domain {}

    // One of the arguments from this list "any_path path_prefix path_regex" must be set

    any_path {}
  }
}
