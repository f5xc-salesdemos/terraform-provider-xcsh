# Certificate Chain Data Source Example
# Retrieves information about an existing Certificate Chain

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Certificate Chain by name
data "f5xc_certificate_chain" "example" {
  name      = "example-certificate-chain"
  namespace = "staging"
}

output "certificate_chain_id" {
  value = data.f5xc_certificate_chain.example.id
}
