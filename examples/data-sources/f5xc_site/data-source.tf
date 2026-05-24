# Site Data Source Example
# Retrieves information about an existing Site

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Site by name
data "f5xc_site" "example" {
  name      = "example-site"
  namespace = "staging"
}

output "site_id" {
  value = data.f5xc_site.example.id
}
