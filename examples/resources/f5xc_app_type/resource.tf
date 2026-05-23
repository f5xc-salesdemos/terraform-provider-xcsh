# App Type Resource Example
# Manages App type will create the configuration in namespace metadata.namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source = "f5xc-salesdemos/f5xc"
    }
  }
}

# Basic App Type configuration
resource "f5xc_app_type" "example" {
  name      = "example-app-type"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Settings specifying how API Discovery will be performed.
  business_logic_markup_setting {
    # Configure business_logic_markup_setting settings
  }
  # Enable this option
  disable_spec {
    # Configure disable_spec settings
  }
  # Discovered API Settings. Configure Discovered API Settings.
  discovered_api_settings {
    # Configure discovered_api_settings settings
  }
}
