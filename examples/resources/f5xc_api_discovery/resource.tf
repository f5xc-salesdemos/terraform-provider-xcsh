# API Discovery Resource Example
# Manages API discovery creates a new object in the storage backend for metadata.namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic API Discovery configuration
resource "f5xc_api_discovery" "example" {
  name      = "example-api-discovery"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Select your custom authentication types to be detected in...
  custom_auth_types {
    # Configure custom_auth_types settings
  }
}

# The following optional fields have server-applied defaults and can be omitted:
# - custom_auth_types
