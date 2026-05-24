# Virtual Site Data Source Example
# Retrieves information about an existing Virtual Site

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Virtual Site by name
data "f5xc_virtual_site" "example" {
  name      = "example-virtual-site"
  namespace = "system"
}

output "virtual_site_id" {
  value = data.f5xc_virtual_site.example.id
}

# Example: Reference virtual site for site selection
# resource "f5xc_http_loadbalancer" "example" {
#   name      = "vs-advertised-lb"
#   namespace = "staging"
#
#   advertise_custom {
#     advertise_where {
#       virtual_site {
#         virtual_site {
#           name      = data.f5xc_virtual_site.example.name
#           namespace = data.f5xc_virtual_site.example.namespace
#         }
#       }
#     }
#   }
# }
