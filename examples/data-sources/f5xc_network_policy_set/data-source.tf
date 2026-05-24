# Network Policy Set Data Source Example
# Retrieves information about an existing Network Policy Set

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Network Policy Set by name
data "f5xc_network_policy_set" "example" {
  name      = "example-network-policy-set"
  namespace = "staging"
}

output "network_policy_set_id" {
  value = data.f5xc_network_policy_set.example.id
}
