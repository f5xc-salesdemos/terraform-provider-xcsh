# Policer Resource Example
# Manages new policer with traffic rate limits. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Basic Policer configuration
resource "xcsh_policer" "example" {
  name      = "example-policer"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # API-discovered default values (shown for reference)
  # These values are applied by the API if not specified
  # policer_mode = "POLICER_MODE_NOT_SHARED"  # API default
  # policer_type = "POLICER_SINGLE_RATE_TWO_COLOR"  # API default
}

# The following optional fields have server-applied defaults and can be omitted:
# - policer_mode
# - policer_type
