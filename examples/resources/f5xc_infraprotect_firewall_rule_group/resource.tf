# Infraprotect Firewall Rule Group Resource Example
# Manages a Infraprotect Firewall Rule Group resource in F5 Distributed Cloud for amends a ddos transit firewall rule group configuration.

# Basic Infraprotect Firewall Rule Group configuration
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


resource "f5xc_infraprotect_firewall_rule_group" "example" {
  name      = "example-infraprotect-firewall-rule-group"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }
}
