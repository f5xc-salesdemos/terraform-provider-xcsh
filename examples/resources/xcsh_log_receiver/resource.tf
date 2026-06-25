# Log Receiver Resource Example
# Manages new Log Receiver object. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Basic Log Receiver configuration
resource "xcsh_log_receiver" "example" {
  name      = "example-log-receiver"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Log Receiver configuration
  # HTTP receiver example
  http_receiver {
    uri = "https://logs.example.com/ingest"
    batch {
      max_bytes       = 1048576
      max_events      = 100
      timeout_seconds = 5
    }
    no_tls_verify_hostname {}
    no_compression {}
  }
}
