# Cloud Connect Data Source Example
# Retrieves information about an existing Cloud Connect

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Cloud Connect by name
data "f5xc_cloud_connect" "example" {
  name      = "example-cloud-connect"
  namespace = "staging"
}

output "cloud_connect_id" {
  value = data.f5xc_cloud_connect.example.id
}
