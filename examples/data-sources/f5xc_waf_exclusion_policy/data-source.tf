# WAF Exclusion Policy Data Source Example
# Retrieves information about an existing WAF Exclusion Policy

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing WAF Exclusion Policy by name
data "f5xc_waf_exclusion_policy" "example" {
  name      = "example-waf-exclusion-policy"
  namespace = "shared"
}

output "waf_exclusion_policy_id" {
  value = data.f5xc_waf_exclusion_policy.example.id
}
