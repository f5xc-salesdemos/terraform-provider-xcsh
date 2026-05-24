# Voltstack Site Data Source Example
# Retrieves information about an existing Voltstack Site

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Voltstack Site by name
data "f5xc_voltstack_site" "example" {
  name      = "example-voltstack-site"
  namespace = "staging"
}

output "voltstack_site_id" {
  value = data.f5xc_voltstack_site.example.id
}
