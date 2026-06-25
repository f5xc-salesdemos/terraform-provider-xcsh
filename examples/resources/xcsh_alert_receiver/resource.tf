# Alert Receiver Resource Example
# Manages new Alert Receiver object. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Alert Receiver configuration
resource "xcsh_alert_receiver" "example" {
  name      = "example-alert-receiver"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Alert Receiver configuration
  # Slack configuration
  slack {
    url = "https://your-slack-webhook-url"
  }
}
