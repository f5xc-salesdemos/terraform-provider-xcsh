# Service Policy Data Source Example
# Retrieves information about an existing Service Policy

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Service Policy by name
data "f5xc_service_policy" "example" {
  name      = "example-service-policy"
  namespace = "staging"
}

output "service_policy_id" {
  value = data.f5xc_service_policy.example.id
}

# Example: Reference service policy in HTTP load balancer
# resource "f5xc_http_loadbalancer" "example" {
#   name      = "policy-protected-lb"
#   namespace = "staging"
#
#   active_service_policies {
#     policies {
#       name      = data.f5xc_service_policy.example.name
#       namespace = data.f5xc_service_policy.example.namespace
#     }
#   }
# }
