# Azure VNET Site Data Source Example
# Retrieves information about an existing Azure VNET Site

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Azure VNET Site by name
data "xcsh_azure_vnet_site" "example" {
  name      = "example-azure-vnet-site"
  namespace = "staging"
}

output "azure_vnet_site_id" {
  value = data.xcsh_azure_vnet_site.example.id
}

# Example: Reference cloud site for advertising load balancer
# resource "xcsh_http_loadbalancer" "example" {
#   name      = "site-advertised-lb"
#   namespace = "staging"
#
#   advertise_custom {
#     advertise_where {
#       site {
#         site {
#           name      = data.xcsh_azure_vnet_site.example.name
#           namespace = data.xcsh_azure_vnet_site.example.namespace
#         }
#       }
#     }
#   }
# }
