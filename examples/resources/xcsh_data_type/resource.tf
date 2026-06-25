# Data Type Resource Example
# Manages data_type creates a new object in the storage backend for metadata.namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Basic Data Type configuration
resource "xcsh_data_type" "example" {
  name      = "example-data-type"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Configure key/value or regex match rules to enable the pl...
  rules {
    # Configure rules settings
  }
  # Configuration parameter for key pattern.
  key_pattern {
    # Configure key_pattern settings
  }
  # Configuration parameter for exact values.
  exact_values {
    # Configure exact_values settings
  }
}
