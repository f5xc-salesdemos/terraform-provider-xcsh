# DNS Proxy Data Source Example
# Retrieves information about an existing DNS Proxy

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing DNS Proxy by name
data "xcsh_dns_proxy" "example" {
  name      = "example-dns-proxy"
  namespace = "staging"
}

output "dns_proxy_id" {
  value = data.xcsh_dns_proxy.example.id
}
