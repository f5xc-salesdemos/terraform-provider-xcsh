# Tunnel Data Source Example
# Retrieves information about an existing Tunnel

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Tunnel by name
data "f5xc_tunnel" "example" {
  name      = "example-tunnel"
  namespace = "system"
}

output "tunnel_id" {
  value = data.f5xc_tunnel.example.id
}
