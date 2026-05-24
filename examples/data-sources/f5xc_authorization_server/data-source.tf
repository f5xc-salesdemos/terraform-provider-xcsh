# Authorization Server Data Source Example
# Retrieves information about an existing Authorization Server

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Authorization Server by name
data "f5xc_authorization_server" "example" {
  name      = "example-authorization-server"
  namespace = "staging"
}

output "authorization_server_id" {
  value = data.f5xc_authorization_server.example.id
}
