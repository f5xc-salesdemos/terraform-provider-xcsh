# Workload Resource Example
# Manages a Workload resource in F5 Distributed Cloud for workload. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Basic Workload configuration
resource "xcsh_workload" "example" {
  name      = "example-workload"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  // One of the arguments from this list "job service stateful_service" must be set

  service {
    containers {
      name = "web"
      image {
        name        = "nginx"
        pull_policy = "IMAGE_PULL_POLICY_ALWAYS"
        public {}
      }
    }
    deploy_options {
      default_virtual_sites {}
    }
  }
}
