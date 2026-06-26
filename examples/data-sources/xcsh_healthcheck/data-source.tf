# Healthcheck Data Source Example
# Retrieves information about an existing Healthcheck

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Healthcheck by name
data "xcsh_healthcheck" "example" {
  name      = "example-healthcheck"
  namespace = "staging"
}

output "healthcheck_id" {
  value = data.xcsh_healthcheck.example.id
}

# Example: Reference healthcheck in origin pool
# resource "xcsh_origin_pool" "example" {
#   name      = "example-pool"
#   namespace = "staging"
#
#   healthcheck {
#     name      = data.xcsh_healthcheck.example.name
#     namespace = data.xcsh_healthcheck.example.namespace
#   }
# }
