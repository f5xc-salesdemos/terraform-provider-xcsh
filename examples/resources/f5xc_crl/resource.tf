# CRL Resource Example
# Manages a CRL resource in F5 Distributed Cloud for api to create crl object. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source = "f5xc-salesdemos/f5xc"
    }
  }
}

# Basic CRL configuration
resource "f5xc_crl" "example" {
  name      = "example-crl"
  namespace = "shared"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Configuration parameter for http access.
  http_access {
    # Configure http_access settings
  }
}
