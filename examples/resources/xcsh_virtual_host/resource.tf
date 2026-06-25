# Virtual Host Resource Example
# Manages virtual host in a given namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Virtual Host configuration
resource "xcsh_virtual_host" "example" {
  name      = "example-virtual-host"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  domains                     = ["app.example.com"]
  proxy                       = "DNS_PROXY"
  idle_timeout                = 30000
  connection_idle_timeout     = 120000
  max_request_header_size     = 32768
  add_location                = false
  disable_dns_resolve         = false
  disable_default_error_pages = false
  request_headers_to_remove   = []
  response_headers_to_remove  = []
  request_cookies_to_remove   = []
  response_cookies_to_remove  = []
}
