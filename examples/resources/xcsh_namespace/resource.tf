# Namespace Resource Example
# Manages new namespace. Name of the object is name of the name space. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Namespace configuration
resource "xcsh_namespace" "example" {
  name      = "example-namespace"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Namespace configuration
  description = "Example namespace for application workloads"
}
