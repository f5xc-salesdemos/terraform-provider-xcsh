# Virtual K8S Resource Example
# Manages virtual_k8s will create the object in the storage backend for namespace metadata.namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Virtual K8S configuration
resource "xcsh_virtual_k8s" "example" {
  name      = "example-virtual-k8s"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Type establishes a direct reference from one object(the r...
  default_flavor_ref {
    # Configure default_flavor_ref settings
  }
  # [OneOf: disabled, isolated] Enable this option
  disabled {
    # Configure disabled settings
  }
  # Enable this option
  isolated {
    # Configure isolated settings
  }
}
