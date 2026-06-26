# GCP VPC Site Data Source Example
# Retrieves information about an existing GCP VPC Site

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing GCP VPC Site by name
data "xcsh_gcp_vpc_site" "example" {
  name      = "example-gcp-vpc-site"
  namespace = "staging"
}

output "gcp_vpc_site_id" {
  value = data.xcsh_gcp_vpc_site.example.id
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
#           name      = data.xcsh_gcp_vpc_site.example.name
#           namespace = data.xcsh_gcp_vpc_site.example.namespace
#         }
#       }
#     }
#   }
# }
