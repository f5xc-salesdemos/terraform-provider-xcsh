# Namespace Data Source Example
# Retrieves information about an existing Namespace

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Namespace by name
data "f5xc_namespace" "example" {
  name      = "example-namespace"
  namespace = "staging"
}

output "namespace_id" {
  value = data.f5xc_namespace.example.id
}

# Example: Create resources in a namespace discovered via data source
# resource "f5xc_origin_pool" "example" {
#   name      = "example-pool"
#   namespace = data.f5xc_namespace.example.name
#   # ... other configuration
# }
