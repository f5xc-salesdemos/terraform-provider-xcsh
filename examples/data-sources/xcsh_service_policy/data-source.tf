# Service Policy Data Source Example
# Retrieves information about an existing Service Policy

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Service Policy by name
data "xcsh_service_policy" "example" {
  name      = "example-service-policy"
  namespace = "staging"
}

output "service_policy_id" {
  value = data.xcsh_service_policy.example.id
}

# Example: Reference service policy in HTTP load balancer
# resource "xcsh_http_loadbalancer" "example" {
#   name      = "policy-protected-lb"
#   namespace = "staging"
#
#   active_service_policies {
#     policies {
#       name      = data.xcsh_service_policy.example.name
#       namespace = data.xcsh_service_policy.example.namespace
#     }
#   }
# }
