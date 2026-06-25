# Container Registry Resource Example
# Manages a Container Registry resource in F5 Distributed Cloud for container image registry configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Container Registry configuration
resource "xcsh_container_registry" "example" {
  name      = "example-container-registry"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # SecretType is used in an object to indicate a sensitive/c...
  password {
    # Configure password settings
  }
  # BlindfoldSecretInfoType specifies information about the S...
  blindfold_secret_info {
    # Configure blindfold_secret_info settings
  }
  # ClearSecretInfoType specifies information about the Secre...
  clear_secret_info {
    # Configure clear_secret_info settings
  }
}
