# Rate Limiter Policy Resource Example
# Manages a Rate Limiter Policy resource in F5 Distributed Cloud for rate limiter policy create specification. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Rate Limiter Policy configuration
resource "f5xc_rate_limiter_policy" "example" {
  name      = "example-rate-limiter-policy"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # [OneOf: any_server, server_name, server_name_matcher, ser...
  any_server {
    # Configure any_server settings
  }
  # Matcher specifies multiple criteria for matching an input...
  server_name_matcher {
    # Configure server_name_matcher settings
  }
  # Type can be used to establish a 'selector reference' from...
  server_selector {
    # Configure server_selector settings
  }
}

# The following optional fields have server-applied defaults and can be omitted:
# - rules
