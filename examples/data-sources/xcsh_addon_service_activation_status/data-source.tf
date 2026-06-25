# Addon Service Activation Status Data Source Example
# Retrieves information about an existing Addon Service Activation Status

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Addon Service Activation Status by name
data "xcsh_addon_service_activation_status" "example" {
  name      = "example-addon-service-activation-status"
  namespace = "staging"
}

output "addon_service_activation_status_id" {
  value = data.xcsh_addon_service_activation_status.example.id
}
