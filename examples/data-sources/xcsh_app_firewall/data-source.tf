# App Firewall Data Source Example
# Retrieves information about an existing App Firewall

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing App Firewall by name
data "xcsh_app_firewall" "example" {
  name      = "example-app-firewall"
  namespace = "staging"
}

output "app_firewall_id" {
  value = data.xcsh_app_firewall.example.id
}

# Example: Reference WAF in HTTP load balancer
# resource "xcsh_http_loadbalancer" "example" {
#   name      = "protected-lb"
#   namespace = "staging"
#
#   app_firewall {
#     name      = data.xcsh_app_firewall.example.name
#     namespace = data.xcsh_app_firewall.example.namespace
#   }
# }
