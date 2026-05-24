# Service Policy Resource Example
# Manages service_policy creates a new object in the storage backend for metadata.namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Service Policy configuration
resource "f5xc_service_policy" "example" {
  name      = "example-service-policy"
  namespace = "shared"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Service Policy configuration
  algo = "FIRST_MATCH"

  # Allow specific paths
  rules {
    metadata {
      name = "allow-api"
    }
    spec {
      action = "ALLOW"
      path {
        prefix = "/api/"
      }
    }
  }
}

# The following optional fields have server-applied defaults and can be omitted:
# - port_matcher
# - any_server
