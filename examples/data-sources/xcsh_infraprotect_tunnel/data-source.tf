# Infraprotect Tunnel Data Source Example
# Retrieves information about an existing Infraprotect Tunnel

# Look up an existing Infraprotect Tunnel by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_infraprotect_tunnel" "example" {
  name      = "example-infraprotect-tunnel"
  namespace = "system"
}

output "infraprotect_tunnel_id" {
  value = data.xcsh_infraprotect_tunnel.example.id
}
