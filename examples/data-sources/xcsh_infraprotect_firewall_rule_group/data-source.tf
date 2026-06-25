# Infraprotect Firewall Rule Group Data Source Example
# Retrieves information about an existing Infraprotect Firewall Rule Group

# Look up an existing Infraprotect Firewall Rule Group by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_infraprotect_firewall_rule_group" "example" {
  name      = "example-infraprotect-firewall-rule-group"
  namespace = "system"
}

output "infraprotect_firewall_rule_group_id" {
  value = data.xcsh_infraprotect_firewall_rule_group.example.id
}
