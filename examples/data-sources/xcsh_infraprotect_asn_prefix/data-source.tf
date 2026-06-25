# Infraprotect Asn Prefix Data Source Example
# Retrieves information about an existing Infraprotect Asn Prefix

# Look up an existing Infraprotect Asn Prefix by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_infraprotect_asn_prefix" "example" {
  name      = "example-infraprotect-asn-prefix"
  namespace = "system"
}

output "infraprotect_asn_prefix_id" {
  value = data.xcsh_infraprotect_asn_prefix.example.id
}
