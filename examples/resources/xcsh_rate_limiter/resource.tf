# Rate Limiter Resource Example
# Manages rate_limiter creates a new object in the storage backend for metadata.namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Rate Limiter configuration
resource "xcsh_rate_limiter" "example" {
  name      = "example-rate-limiter"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Rate Limiter configuration
  limits {
    total_number     = 100
    unit             = "MINUTE"
    burst_multiplier = 10
  }
}

# The following optional fields have server-applied defaults and can be omitted:
# - user_identification
