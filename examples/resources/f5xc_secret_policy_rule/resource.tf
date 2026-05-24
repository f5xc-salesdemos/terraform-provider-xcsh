# Secret Policy Rule Resource Example
# Manages secret_policy_rule creates a new object in storage backend for metadata.namespace. in F5 Distributed Cloud.

# Basic Secret Policy Rule configuration
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


resource "f5xc_secret_policy_rule" "example" {
  name      = "example-secret-policy-rule"
  namespace = "shared"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Matcher. A matcher specifies multiple criteria for matchi...
  client_name_matcher {
    # Configure client_name_matcher settings
  }
  # Label Selector. This type can be used to establish a 'sel...
  client_selector {
    # Configure client_selector settings
  }
}
