# Securemesh Site Data Source Example
# Retrieves information about an existing Securemesh Site

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Securemesh Site by name
data "f5xc_securemesh_site" "example" {
  name      = "example-securemesh-site"
  namespace = "system"
}

output "securemesh_site_id" {
  value = data.f5xc_securemesh_site.example.id
}
