# Tenant Configuration Resource Example
# Manages a Tenant Configuration resource in F5 Distributed Cloud for tenant configuration specification. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source = "f5xc-salesdemos/f5xc"
    }
  }
}

# Basic Tenant Configuration configuration
resource "f5xc_tenant_configuration" "example" {
  name      = "example-tenant-configuration"
  namespace = "system"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Configuration parameter for basic configuration.
  basic_configuration {
    # Configure basic_configuration settings
  }
  # Configuration parameter for brute force detection.
  brute_force_detection {
    # Configure brute_force_detection settings
  }
  # Configuration parameter for brute force detection settings.
  brute_force_detection_settings {
    # Configure brute_force_detection_settings settings
  }
}
