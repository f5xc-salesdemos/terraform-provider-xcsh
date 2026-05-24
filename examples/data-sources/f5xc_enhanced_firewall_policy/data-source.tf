# Enhanced Firewall Policy Data Source Example
# Retrieves information about an existing Enhanced Firewall Policy

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Enhanced Firewall Policy by name
data "f5xc_enhanced_firewall_policy" "example" {
  name      = "example-enhanced-firewall-policy"
  namespace = "staging"
}

output "enhanced_firewall_policy_id" {
  value = data.f5xc_enhanced_firewall_policy.example.id
}
