# Policer Resource Example
# Manages new policer with traffic rate limits. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source = "f5xc-salesdemos/f5xc"
    }
  }
}

# Basic Policer configuration
resource "f5xc_policer" "example" {
  name      = "example-policer"
  namespace = "shared"

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
