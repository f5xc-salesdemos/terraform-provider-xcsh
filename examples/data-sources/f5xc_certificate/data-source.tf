# Certificate Data Source Example
# Retrieves information about an existing Certificate

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Certificate by name
data "f5xc_certificate" "example" {
  name      = "example-certificate"
  namespace = "staging"
}

output "certificate_id" {
  value = data.f5xc_certificate.example.id
}

# Example: Reference certificate in HTTPS configuration
# resource "f5xc_http_loadbalancer" "example" {
#   name      = "https-lb"
#   namespace = "staging"
#
#   https {
#     tls_cert_params {
#       certificates {
#         name      = data.f5xc_certificate.example.name
#         namespace = data.f5xc_certificate.example.namespace
#       }
#     }
#   }
# }
