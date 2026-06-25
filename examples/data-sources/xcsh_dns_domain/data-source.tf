# DNS Domain Data Source Example
# Retrieves information about an existing DNS Domain

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing DNS Domain by name
data "xcsh_dns_domain" "example" {
  name      = "example-dns-domain"
  namespace = "staging"
}

output "dns_domain_id" {
  value = data.xcsh_dns_domain.example.id
}
