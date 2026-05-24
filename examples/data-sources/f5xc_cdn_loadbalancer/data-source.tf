# CDN Loadbalancer Data Source Example
# Retrieves information about an existing CDN Loadbalancer

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing CDN Loadbalancer by name
data "f5xc_cdn_loadbalancer" "example" {
  name      = "example-cdn-loadbalancer"
  namespace = "staging"
}

output "cdn_loadbalancer_id" {
  value = data.f5xc_cdn_loadbalancer.example.id
}
