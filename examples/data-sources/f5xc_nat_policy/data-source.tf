# NAT Policy Data Source Example
# Retrieves information about an existing NAT Policy

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing NAT Policy by name
data "f5xc_nat_policy" "example" {
  name      = "example-nat-policy"
  namespace = "staging"
}

output "nat_policy_id" {
  value = data.f5xc_nat_policy.example.id
}
