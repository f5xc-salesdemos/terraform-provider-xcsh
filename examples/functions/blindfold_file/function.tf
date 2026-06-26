terraform {
  required_version = ">= 1.8"
  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Encrypt a file using F5XC blindfold
#
# The blindfold_file function reads a file and encrypts its contents using F5
# Distributed Cloud Secret Management. This is a convenience function equivalent
# to: provider::xcsh::blindfold(base64encode(file(path)), policy_name, namespace)
#
# The encryption happens locally - file contents are never transmitted to F5XC.

# Example: Encrypt a TLS private key file
resource "xcsh_http_loadbalancer" "secure" {
  name      = "secure-lb"
  namespace = "production"

  domains = ["secure.example.com"]

  https_auto_cert {
    tls_config {
      custom_security {
        private_key {
          blindfold_secret_info {
            location = provider::xcsh::blindfold_file(
              "${path.module}/certs/server.key",
              "tls-secrets-policy",
              "shared"
            )
          }
        }
        certificate {
          certificate_url = "string:///${base64encode(file("${path.module}/certs/server.crt"))}"
        }
      }
    }
  }
}

# Example: Encrypt multiple certificate files using for_each
locals {
  certificates = {
    "server" = "${path.module}/certs/server.key"
    "client" = "${path.module}/certs/client.key"
    "ca"     = "${path.module}/certs/ca.key"
  }
}

resource "xcsh_certificate" "certs" {
  for_each  = local.certificates
  name      = each.key
  namespace = "production"

  private_key {
    blindfold_secret_info {
      location = provider::xcsh::blindfold_file(
        each.value,
        "cert-secrets-policy",
        "shared"
      )
    }
  }
}
