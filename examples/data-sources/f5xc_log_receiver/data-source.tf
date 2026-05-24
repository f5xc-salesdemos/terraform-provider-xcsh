# Log Receiver Data Source Example
# Retrieves information about an existing Log Receiver

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Log Receiver by name
data "f5xc_log_receiver" "example" {
  name      = "example-log-receiver"
  namespace = "staging"
}

output "log_receiver_id" {
  value = data.f5xc_log_receiver.example.id
}

# Example: Reference log receiver in site configuration
# resource "f5xc_securemesh_site_v2" "example" {
#   name      = "example-site"
#   namespace = "staging"
#
#   log_receiver {
#     name      = data.f5xc_log_receiver.example.name
#     namespace = data.f5xc_log_receiver.example.namespace
#   }
# }
