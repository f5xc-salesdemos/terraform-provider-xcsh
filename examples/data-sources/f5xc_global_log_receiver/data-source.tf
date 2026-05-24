# Global Log Receiver Data Source Example
# Retrieves information about an existing Global Log Receiver

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Global Log Receiver by name
data "f5xc_global_log_receiver" "example" {
  name      = "example-global-log-receiver"
  namespace = "staging"
}

output "global_log_receiver_id" {
  value = data.f5xc_global_log_receiver.example.id
}
