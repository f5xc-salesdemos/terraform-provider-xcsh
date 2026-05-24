# Nginx Server Data Source Example
# Retrieves information about an existing Nginx Server

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Nginx Server by name
data "f5xc_nginx_server" "example" {
  name      = "example-nginx-server"
  namespace = "staging"
}

output "nginx_server_id" {
  value = data.f5xc_nginx_server.example.id
}
