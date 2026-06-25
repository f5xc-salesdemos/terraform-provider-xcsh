# Code Base Integration Resource Example
# Manages integration details. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Code Base Integration configuration
resource "xcsh_code_base_integration" "example" {
  name      = "example-code-base-integration"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Choose your code base (e.g. GitHub, GitLab, Bitbucket, Az...
  code_base_integration {
    # Configure code_base_integration settings
  }
  # Configuration parameter for azure repos.
  azure_repos {
    # Configure azure_repos settings
  }
  # SecretType is used in an object to indicate a sensitive/c...
  access_token {
    # Configure access_token settings
  }
}
