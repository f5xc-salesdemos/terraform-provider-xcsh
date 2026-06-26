# Origin Pool Data Source Example
# Retrieves information about an existing Origin Pool

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Origin Pool by name
data "xcsh_origin_pool" "example" {
  name      = "example-origin-pool"
  namespace = "staging"
}

output "origin_pool_id" {
  value = data.xcsh_origin_pool.example.id
}

# Example: Use origin pool data in HTTP load balancer
# resource "xcsh_http_loadbalancer" "example" {
#   name      = "example-lb"
#   namespace = "staging"
#   domains   = ["app.example.com"]
#
#   default_route_pools {
#     pool {
#       name      = data.xcsh_origin_pool.example.name
#       namespace = data.xcsh_origin_pool.example.namespace
#     }
#   }
# }
