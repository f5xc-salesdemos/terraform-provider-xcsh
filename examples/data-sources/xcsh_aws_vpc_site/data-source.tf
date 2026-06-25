# AWS VPC Site Data Source Example
# Retrieves information about an existing AWS VPC Site

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing AWS VPC Site by name
data "xcsh_aws_vpc_site" "example" {
  name      = "example-aws-vpc-site"
  namespace = "staging"
}

output "aws_vpc_site_id" {
  value = data.xcsh_aws_vpc_site.example.id
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
#           name      = data.xcsh_aws_vpc_site.example.name
#           namespace = data.xcsh_aws_vpc_site.example.namespace
#         }
#       }
#     }
#   }
# }
