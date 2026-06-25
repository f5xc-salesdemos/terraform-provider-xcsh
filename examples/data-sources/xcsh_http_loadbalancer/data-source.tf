# HTTP Loadbalancer Data Source Example
# Retrieves information about an existing HTTP Loadbalancer

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing HTTP Loadbalancer by name
data "xcsh_http_loadbalancer" "example" {
  name      = "example-http-loadbalancer"
  namespace = "staging"
}

output "http_loadbalancer_id" {
  value = data.xcsh_http_loadbalancer.example.id
}

# Example: Reference in another load balancer configuration
# resource "xcsh_service_policy" "example" {
#   name      = "policy-for-lb"
#   namespace = "staging"
#
#   # Use the load balancer's domains
#   # domain = data.xcsh_http_loadbalancer.example.domains[0]
# }
