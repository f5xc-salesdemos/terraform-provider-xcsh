# Certificate Resource Example
# Manages a Certificate resource in F5 Distributed Cloud for certificate. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Certificate configuration
resource "xcsh_certificate" "example" {
  name      = "example-certificate"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Certificate configuration
  certificate_url = "string:///LS0tLS1CRUdJTi..."
  private_key {
    clear_secret_info {
      url = "string:///LS0tLS1CRUdJTi..."
    }
  }
}
