# Log Receiver Data Source Example
# Retrieves information about an existing Log Receiver

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Log Receiver by name
data "xcsh_log_receiver" "example" {
  name      = "example-log-receiver"
  namespace = "staging"
}

output "log_receiver_id" {
  value = data.xcsh_log_receiver.example.id
}

# Example: Reference log receiver in site configuration
# resource "xcsh_securemesh_site_v2" "example" {
#   name      = "example-site"
#   namespace = "staging"
#
#   log_receiver {
#     name      = data.xcsh_log_receiver.example.name
#     namespace = data.xcsh_log_receiver.example.namespace
#   }
# }
