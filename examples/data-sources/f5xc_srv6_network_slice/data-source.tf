# Srv6 Network Slice Data Source Example
# Retrieves information about an existing Srv6 Network Slice

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Srv6 Network Slice by name
data "f5xc_srv6_network_slice" "example" {
  name      = "example-srv6-network-slice"
  namespace = "staging"
}

output "srv6_network_slice_id" {
  value = data.f5xc_srv6_network_slice.example.id
}
