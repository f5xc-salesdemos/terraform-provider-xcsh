# Rate Limiter Data Source Example
# Retrieves information about an existing Rate Limiter

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Rate Limiter by name
data "f5xc_rate_limiter" "example" {
  name      = "example-rate-limiter"
  namespace = "staging"
}

output "rate_limiter_id" {
  value = data.f5xc_rate_limiter.example.id
}

# Example: Reference rate limiter in HTTP load balancer
# resource "f5xc_http_loadbalancer" "example" {
#   name      = "rate-limited-lb"
#   namespace = "staging"
#
#   rate_limit {
#     rate_limiter {
#       name      = data.f5xc_rate_limiter.example.name
#       namespace = data.f5xc_rate_limiter.example.namespace
#     }
#   }
# }
