# Infraprotect Deny List Rule Data Source Example
# Retrieves information about an existing Infraprotect Deny List Rule

# Look up an existing Infraprotect Deny List Rule by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_infraprotect_deny_list_rule" "example" {
  name      = "example-infraprotect-deny-list-rule"
  namespace = "system"
}

output "infraprotect_deny_list_rule_id" {
  value = data.xcsh_infraprotect_deny_list_rule.example.id
}
