# CRL Data Source Example
# Retrieves information about an existing CRL

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing CRL by name
data "f5xc_crl" "example" {
  name      = "example-crl"
  namespace = "shared"
}

output "crl_id" {
  value = data.f5xc_crl.example.id
}
