# DNS Zone Data Source Example
# Retrieves information about an existing DNS Zone

# Look up an existing DNS Zone by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_dns_zone" "example" {
  name      = "example-dns-zone"
  namespace = "system"
}


# Example: Reference DNS zone in DNS load balancer
# resource "xcsh_dns_load_balancer" "example" {
#   name      = "example-dns-lb"
#   namespace = "system"
#
#   dns_zone {
#     name      = data.xcsh_dns_zone.example.name
#     namespace = data.xcsh_dns_zone.example.namespace
#   }
# }

output "dns_zone_id" {
  value = data.xcsh_dns_zone.example.id
}
