# API Crawler Resource Example
# Manages a API Crawler resource in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic API Crawler configuration
resource "f5xc_api_crawler" "example" {
  name      = "example-api-crawler"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # API Crawler Configuration .
  domains {
    # Configure domains settings
  }
  # Configuration parameter for simple login.
  simple_login {
    # Configure simple_login settings
  }
  # SecretType is used in an object to indicate a sensitive/c...
  password {
    # Configure password settings
  }
}
