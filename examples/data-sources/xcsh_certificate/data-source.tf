# Certificate Data Source Example
# Retrieves information about an existing Certificate

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Certificate by name
data "xcsh_certificate" "example" {
  name      = "example-certificate"
  namespace = "staging"
}

output "certificate_id" {
  value = data.xcsh_certificate.example.id
}

# Example: Reference certificate in HTTPS configuration
# resource "xcsh_http_loadbalancer" "example" {
#   name      = "https-lb"
#   namespace = "staging"
#
#   https {
#     tls_cert_params {
#       certificates {
#         name      = data.xcsh_certificate.example.name
#         namespace = data.xcsh_certificate.example.namespace
#       }
#     }
#   }
# }
