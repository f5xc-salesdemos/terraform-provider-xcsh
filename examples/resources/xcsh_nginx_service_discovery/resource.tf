# Nginx Service Discovery Resource Example
# Manages a Nginx Service Discovery resource in F5 Distributed Cloud for api to create nginx service discovery object for a site or virtual site in system namespace. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Nginx Service Discovery configuration
resource "xcsh_nginx_service_discovery" "example" {
  name      = "example-nginx-service-discovery"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Configuration parameter for discovery target.
  discovery_target {
    # Configure discovery_target settings
  }
  # Configuration parameter for config sync group.
  config_sync_group {
    # Configure config_sync_group settings
  }
  # NGINXInstance Reference. Select new NGINX Instance.
  nginx_instance {
    # Configure nginx_instance settings
  }
}
