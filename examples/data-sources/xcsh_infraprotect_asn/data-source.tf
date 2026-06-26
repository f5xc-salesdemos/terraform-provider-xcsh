# Infraprotect Asn Data Source Example
# Retrieves information about an existing Infraprotect Asn

# Look up an existing Infraprotect Asn by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_infraprotect_asn" "example" {
  name      = "example-infraprotect-asn"
  namespace = "system"
}

output "infraprotect_asn_id" {
  value = data.xcsh_infraprotect_asn.example.id
}
