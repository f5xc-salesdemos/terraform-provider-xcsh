# Infraprotect Asn Data Source Example
# Retrieves information about an existing Infraprotect Asn

# Look up an existing Infraprotect Asn by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_infraprotect_asn" "example" {
  name      = "example-infraprotect-asn"
  namespace = "system"
}

# Example: Use the data source in another resource
# output "infraprotect_asn_id" {
#   value = data.f5xc_infraprotect_asn.example.id
# }
