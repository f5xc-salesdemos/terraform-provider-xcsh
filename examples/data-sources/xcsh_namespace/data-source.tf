# Namespace Data Source Example
# Retrieves information about an existing Namespace

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Namespace by name
data "xcsh_namespace" "example" {
  name      = "example-namespace"
  namespace = "staging"
}

output "namespace_id" {
  value = data.xcsh_namespace.example.id
}

# Example: Create resources in a namespace discovered via data source
# resource "xcsh_origin_pool" "example" {
#   name      = "example-pool"
#   namespace = data.xcsh_namespace.example.name
#   # ... other configuration
# }
