# Registration Data Source Example
# Retrieves information about an existing Registration

# Look up an existing Registration by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_registration" "example" {
  name      = "example-registration"
  namespace = "system"
}

output "registration_id" {
  value = data.xcsh_registration.example.id
}
