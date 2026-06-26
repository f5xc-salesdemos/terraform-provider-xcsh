# Virtual Network Data Source Example
# Retrieves information about an existing Virtual Network

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Virtual Network by name
data "xcsh_virtual_network" "example" {
  name      = "example-virtual-network"
  namespace = "staging"
}

output "virtual_network_id" {
  value = data.xcsh_virtual_network.example.id
}
