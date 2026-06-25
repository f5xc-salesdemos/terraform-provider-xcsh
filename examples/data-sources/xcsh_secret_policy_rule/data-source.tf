# Secret Policy Rule Data Source Example
# Retrieves information about an existing Secret Policy Rule

# Look up an existing Secret Policy Rule by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_secret_policy_rule" "example" {
  name      = "example-secret-policy-rule"
  namespace = "system"
}

output "secret_policy_rule_id" {
  value = data.xcsh_secret_policy_rule.example.id
}
