# Tunnel Resource Example
# Manages tunnel in a given namespace. If one already exist it will give a error. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Tunnel configuration
resource "f5xc_tunnel" "example" {
  name      = "example-tunnel"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  tunnel_type = "IPSEC_PSK"
}
