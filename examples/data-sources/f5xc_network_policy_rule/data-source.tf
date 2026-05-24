# Network Policy Rule Data Source Example
# Retrieves information about an existing Network Policy Rule

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Network Policy Rule by name
data "f5xc_network_policy_rule" "example" {
  name      = "example-network-policy-rule"
  namespace = "system"
}

output "network_policy_rule_id" {
  value = data.f5xc_network_policy_rule.example.id
}
