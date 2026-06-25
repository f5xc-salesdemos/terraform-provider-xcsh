# Infraprotect Firewall Rule Data Source Example
# Retrieves information about an existing Infraprotect Firewall Rule

# Look up an existing Infraprotect Firewall Rule by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_infraprotect_firewall_rule" "example" {
  name      = "example-infraprotect-firewall-rule"
  namespace = "system"
}

output "infraprotect_firewall_rule_id" {
  value = data.xcsh_infraprotect_firewall_rule.example.id
}
