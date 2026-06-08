# Securemesh Site V2 Data Source Example
# Retrieves information about an existing Securemesh Site V2

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Securemesh Site V2 by name
data "f5xc_securemesh_site_v2" "example" {
  name      = "example-securemesh-site-v2"
  namespace = "staging"
}

output "securemesh_site_v2_id" {
  value = data.f5xc_securemesh_site_v2.example.id
}
