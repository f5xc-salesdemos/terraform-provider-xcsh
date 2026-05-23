# CDN Cache Rule Resource Example
# Manages a CDN Cache Rule resource in F5 Distributed Cloud for cdn loadbalancer specification. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source = "f5xc-salesdemos/f5xc"
    }
  }
}

# Basic CDN Cache Rule configuration
resource "f5xc_cdn_cache_rule" "example" {
  name      = "example-cdn-cache-rule"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Cache Rule. This defines a CDN Cache Rule.
  cache_rules {
    # Configure cache_rules settings
  }
  # Configuration parameter for cache bypass.
  cache_bypass {
    # Configure cache_bypass settings
  }
  # Configuration parameter for eligible for cache.
  eligible_for_cache {
    # Configure eligible_for_cache settings
  }
}
