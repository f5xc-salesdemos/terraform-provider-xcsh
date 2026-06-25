# Service Policy Resource Example
# Manages service_policy creates a new object in the storage backend for metadata.namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Service Policy configuration
resource "xcsh_service_policy" "example" {
  name      = "example-service-policy"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Service Policy configuration
  // One of the arguments from this list "allow_all_requests allow_list deny_all_requests deny_list rule_list" must be set

  rule_list {
    rules {
      metadata {
        name = "allow-api"
      }
      spec {
        action = "ALLOW"
        any_client {}
        any_ip {}
        path {
          prefix_values = ["/api/"]
        }
      }
    }
  }

  // One of the arguments from this list "any_server server_name server_name_matcher server_selector" must be set

  any_server {}
}

# The following optional fields have server-applied defaults and can be omitted:
# - port_matcher
# - any_server
