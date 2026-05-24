# DNS Domain Resource Example
# Manages DNS Domain in a given namespace. If one already exist it will give a error. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic DNS Domain configuration
resource "f5xc_dns_domain" "example" {
  name      = "example-dns-domain"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Configuration parameter for volterra managed.
  volterra_managed {
    # Configure volterra_managed settings
  }
}
