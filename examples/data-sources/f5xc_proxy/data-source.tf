# Proxy Data Source Example
# Retrieves information about an existing Proxy

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Proxy by name
data "f5xc_proxy" "example" {
  name      = "example-proxy"
  namespace = "staging"
}

output "proxy_id" {
  value = data.f5xc_proxy.example.id
}
