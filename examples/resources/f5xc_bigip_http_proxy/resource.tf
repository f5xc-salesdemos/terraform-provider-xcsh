# Bigip HTTP Proxy Resource Example
# Manages BIG-IP HTTP Proxy in a given namespace. If one already exists, it will give an error. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Bigip HTTP Proxy configuration
resource "f5xc_bigip_http_proxy" "example" {
  name      = "example-bigip-http-proxy"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Defines various advanced Profile OPTIONS for a Loadbalancer.
  advanced_profile {
    # Configure advanced_profile settings
  }
  # Enable this option
  disable_spec {
    # Configure disable_spec settings
  }
  # Configuration parameter for enable default profile.
  enable_default_profile {
    # Configure enable_default_profile settings
  }
}
