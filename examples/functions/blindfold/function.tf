terraform {
  required_version = ">= 1.8"
  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Encrypt a secret string using F5XC blindfold
#
# The blindfold function encrypts base64-encoded plaintext using F5 Distributed
# Cloud Secret Management. The encryption happens locally - your secret is never
# transmitted to F5XC during encryption.

# Example: Encrypt a password for use in origin pool authentication
locals {
  encrypted_password = provider::xcsh::blindfold(
    base64encode("example-secret-password"),
    "production-secrets-policy",
    "shared"
  )
}

# Example: Encrypt a TLS private key from a file
locals {
  encrypted_key = provider::xcsh::blindfold(
    base64encode(file("${path.module}/certs/private.key")),
    "tls-secrets-policy",
    "shared"
  )
}

# Example: Using the encrypted secrets in a resource
resource "xcsh_origin_pool" "example" {
  name      = "secure-pool"
  namespace = "production"

  origin_servers {
    private_ip {
      ip = "10.0.0.1"
    }
  }

  port = 443

  # Use the encrypted password from locals
  custom_hash_algorithms {
    hash_algorithms = [local.encrypted_password]
  }
}

resource "xcsh_http_loadbalancer" "example" {
  name      = "secure-lb"
  namespace = "production"

  domains = ["example.com"]

  https_auto_cert {
    tls_config {
      custom_security {
        private_key {
          blindfold_secret_info {
            location = local.encrypted_key
          }
        }
      }
    }
  }
}
