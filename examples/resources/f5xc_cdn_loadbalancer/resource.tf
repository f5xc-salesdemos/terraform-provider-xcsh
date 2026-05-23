# CDN Loadbalancer Resource Example
# Manages a CDN Load Balancer resource in F5 Distributed Cloud for content delivery and edge caching with load balancing.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source = "f5xc-salesdemos/f5xc"
    }
  }
}

# Basic CDN Loadbalancer configuration
resource "f5xc_cdn_loadbalancer" "example" {
  name      = "example-cdn-loadbalancer"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # CDN Load Balancer configuration
  domains = ["cdn.example.com"]

  # Origin pool
  origin_pool {
    public_name {
      dns_name = "origin.example.com"
    }
    follow_origin_redirect = true
    no_tls {}
  }

  # Cache TTL settings
  cache_ttl_options {
    cache_ttl_default = "1h"
  }

  # HTTP protocol
  https_auto_cert {
    http_redirect = true
  }

  # Add location header
  add_location = true
}
