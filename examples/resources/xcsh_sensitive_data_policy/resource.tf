# Sensitive Data Policy Resource Example
# Manages sensitive_data_policy creates a new object in the storage backend for metadata.namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Basic Sensitive Data Policy configuration
resource "xcsh_sensitive_data_policy" "example" {
  name      = "example-sensitive-data-policy"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Select your custom data types to be monitored in the API ...
  custom_data_types {
    # Configure custom_data_types settings
  }
  # Type establishes a direct reference from one object(the r...
  custom_data_type_ref {
    # Configure custom_data_type_ref settings
  }
}

# The following optional fields have server-applied defaults and can be omitted:
# - compliances
# - disabled_predefined_data_types
# - custom_data_types
