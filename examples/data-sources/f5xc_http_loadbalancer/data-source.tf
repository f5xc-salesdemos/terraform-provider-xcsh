# HTTP Loadbalancer Data Source Example
# Retrieves information about an existing HTTP Loadbalancer

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing HTTP Loadbalancer by name
data "f5xc_http_loadbalancer" "example" {
  name      = "example-http-loadbalancer"
  namespace = "staging"
}

output "http_loadbalancer_id" {
  value = data.f5xc_http_loadbalancer.example.id
}

# Example: Reference in another load balancer configuration
# resource "f5xc_service_policy" "example" {
#   name      = "policy-for-lb"
#   namespace = "staging"
#
#   # Use the load balancer's domains
#   # domain = data.f5xc_http_loadbalancer.example.domains[0]
# }
