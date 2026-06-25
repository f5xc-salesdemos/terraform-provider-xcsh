# Alert Receiver Data Source Example
# Retrieves information about an existing Alert Receiver

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Alert Receiver by name
data "xcsh_alert_receiver" "example" {
  name      = "example-alert-receiver"
  namespace = "staging"
}

output "alert_receiver_id" {
  value = data.xcsh_alert_receiver.example.id
}

# Example: Reference alert receiver in alert policy
# resource "xcsh_alert_policy" "example" {
#   name      = "example-policy"
#   namespace = "staging"
#
#   receivers {
#     name      = data.xcsh_alert_receiver.example.name
#     namespace = data.xcsh_alert_receiver.example.namespace
#   }
# }
