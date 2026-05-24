# Infraprotect Asn Prefix Data Source Example
# Retrieves information about an existing Infraprotect Asn Prefix

# Look up an existing Infraprotect Asn Prefix by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_infraprotect_asn_prefix" "example" {
  name      = "example-infraprotect-asn-prefix"
  namespace = "system"
}

# Example: Use the data source in another resource
# output "infraprotect_asn_prefix_id" {
#   value = data.f5xc_infraprotect_asn_prefix.example.id
# }
