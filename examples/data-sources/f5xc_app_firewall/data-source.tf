# App Firewall Data Source Example
# Retrieves information about an existing App Firewall

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing App Firewall by name
data "f5xc_app_firewall" "example" {
  name      = "example-app-firewall"
  namespace = "staging"
}

output "app_firewall_id" {
  value = data.f5xc_app_firewall.example.id
}

# Example: Reference WAF in HTTP load balancer
# resource "f5xc_http_loadbalancer" "example" {
#   name      = "protected-lb"
#   namespace = "staging"
#
#   app_firewall {
#     name      = data.f5xc_app_firewall.example.name
#     namespace = data.f5xc_app_firewall.example.namespace
#   }
# }
