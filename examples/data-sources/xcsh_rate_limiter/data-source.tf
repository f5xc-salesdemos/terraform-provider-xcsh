# Rate Limiter Data Source Example
# Retrieves information about an existing Rate Limiter

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Rate Limiter by name
data "xcsh_rate_limiter" "example" {
  name      = "example-rate-limiter"
  namespace = "staging"
}

output "rate_limiter_id" {
  value = data.xcsh_rate_limiter.example.id
}

# Example: Reference rate limiter in HTTP load balancer
# resource "xcsh_http_loadbalancer" "example" {
#   name      = "rate-limited-lb"
#   namespace = "staging"
#
#   rate_limit {
#     rate_limiter {
#       name      = data.xcsh_rate_limiter.example.name
#       namespace = data.xcsh_rate_limiter.example.namespace
#     }
#   }
# }
