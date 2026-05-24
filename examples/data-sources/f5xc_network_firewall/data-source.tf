# Network Firewall Data Source Example
# Retrieves information about an existing Network Firewall

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Network Firewall by name
data "f5xc_network_firewall" "example" {
  name      = "example-network-firewall"
  namespace = "staging"
}

output "network_firewall_id" {
  value = data.f5xc_network_firewall.example.id
}
