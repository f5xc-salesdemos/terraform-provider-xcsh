# Infraprotect Deny List Rule Data Source Example
# Retrieves information about an existing Infraprotect Deny List Rule

# Look up an existing Infraprotect Deny List Rule by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_infraprotect_deny_list_rule" "example" {
  name      = "example-infraprotect-deny-list-rule"
  namespace = "system"
}

output "infraprotect_deny_list_rule_id" {
  value = data.f5xc_infraprotect_deny_list_rule.example.id
}
