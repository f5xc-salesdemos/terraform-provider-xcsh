# Application Profiles Resource Example
# Manages Application Profiles in a given namespace. If one already exists it will give an error. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Application Profiles configuration
resource "f5xc_application_profiles" "example" {
  name      = "example-application-profiles"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Configuration parameter for advanced tcp profile.
  advanced_tcp_profile {
    # Configure advanced_tcp_profile settings
  }
  # Configuration parameter for disable tcp advanced profile.
  disable_tcp_advanced_profile {
    # Configure disable_tcp_advanced_profile settings
  }
  # Configuration parameter for enable tcp advanced profile.
  enable_tcp_advanced_profile {
    # Configure enable_tcp_advanced_profile settings
  }
}
