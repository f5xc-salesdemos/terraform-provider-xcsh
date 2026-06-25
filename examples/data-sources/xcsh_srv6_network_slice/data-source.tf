# Srv6 Network Slice Data Source Example
# Retrieves information about an existing Srv6 Network Slice

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Srv6 Network Slice by name
data "xcsh_srv6_network_slice" "example" {
  name      = "example-srv6-network-slice"
  namespace = "staging"
}

output "srv6_network_slice_id" {
  value = data.xcsh_srv6_network_slice.example.id
}
