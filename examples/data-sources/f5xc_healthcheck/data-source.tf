# Healthcheck Data Source Example
# Retrieves information about an existing Healthcheck

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Healthcheck by name
data "f5xc_healthcheck" "example" {
  name      = "example-healthcheck"
  namespace = "staging"
}

output "healthcheck_id" {
  value = data.f5xc_healthcheck.example.id
}

# Example: Reference healthcheck in origin pool
# resource "f5xc_origin_pool" "example" {
#   name      = "example-pool"
#   namespace = "staging"
#
#   healthcheck {
#     name      = data.f5xc_healthcheck.example.name
#     namespace = data.f5xc_healthcheck.example.namespace
#   }
# }
