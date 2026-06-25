# Network Firewall Data Source Example
# Retrieves information about an existing Network Firewall

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Network Firewall by name
data "xcsh_network_firewall" "example" {
  name      = "example-network-firewall"
  namespace = "staging"
}

output "network_firewall_id" {
  value = data.xcsh_network_firewall.example.id
}
