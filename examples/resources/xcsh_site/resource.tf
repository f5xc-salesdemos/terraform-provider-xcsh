# Site Resource Example
# Manages virtual site object in given namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Site configuration
resource "xcsh_site" "example" {
  name      = "example-site"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Type can be used to establish a 'selector reference' from...
  site_selector {
    # Configure site_selector settings
  }
}
