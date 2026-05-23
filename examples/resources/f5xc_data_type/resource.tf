# Data Type Resource Example
# Manages data_type creates a new object in the storage backend for metadata.namespace. in F5 Distributed Cloud.

# Basic Data Type configuration
resource "f5xc_data_type" "example" {
  name      = "example-data-type"
  namespace = "shared"

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
