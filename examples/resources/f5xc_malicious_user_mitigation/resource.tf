# Malicious User Mitigation Resource Example
# Manages malicious_user_mitigation creates a new object in the storage backend for metadata.namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Malicious User Mitigation configuration
resource "f5xc_malicious_user_mitigation" "example" {
  name      = "example-malicious-user-mitigation"
  namespace = "shared"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Malicious User Mitigation configuration
  # Detection rules
  rules {
    threat_level = "HIGH"
    mitigation_action {
      block {
        body   = "Access denied"
        status = "403"
      }
    }
  }
}

# The following optional fields have server-applied defaults and can be omitted:
# - mitigation_type
