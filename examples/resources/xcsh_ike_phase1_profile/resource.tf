# IKE Phase1 Profile Resource Example
# Manages a IKE Phase1 Profile resource in F5 Distributed Cloud for ike phase1 profile specification. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic IKE Phase1 Profile configuration
resource "xcsh_ike_phase1_profile" "example" {
  name      = "example-ike-phase1-profile"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # [OneOf: ike_keylifetime_hours, ike_keylifetime_minutes, u...
  ike_keylifetime_hours {
    # Configure ike_keylifetime_hours settings
  }
  # Configuration parameter for ike keylifetime minutes.
  ike_keylifetime_minutes {
    # Configure ike_keylifetime_minutes settings
  }
  # [OneOf: reauth_disabled, reauth_timeout_days, reauth_time...
  reauth_disabled {
    # Configure reauth_disabled settings
  }
}
