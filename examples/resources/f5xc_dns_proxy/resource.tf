# DNS Proxy Resource Example
# Manages DNS Proxy in a given namespace. If one already exists it will give an error. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source = "f5xc-salesdemos/f5xc"
    }
  }
}

# Basic DNS Proxy configuration
resource "f5xc_dns_proxy" "example" {
  name      = "example-dns-proxy"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # DNS Cache specifies cache configuration.
  cache_profile {
    # Configure cache_profile settings
  }
  # Configuration parameter for disable cache profile.
  disable_cache_profile {
    # Configure disable_cache_profile settings
  }
  # Configuration parameter for ddos profile.
  ddos_profile {
    # Configure ddos_profile settings
  }
}
