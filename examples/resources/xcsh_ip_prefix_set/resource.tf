# IP Prefix Set Resource Example
# Manages ip_prefix_set creates a new object in the storage backend for metadata.namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic IP Prefix Set configuration
resource "xcsh_ip_prefix_set" "example" {
  name      = "example-ip-prefix-set"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # IP Prefix Set configuration
  prefix = ["192.168.1.0/24", "10.0.0.0/8"]
}
