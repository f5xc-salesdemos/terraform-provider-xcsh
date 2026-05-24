# Registration Data Source Example
# Retrieves information about an existing Registration

# Look up an existing Registration by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_registration" "example" {
  name      = "example-registration"
  namespace = "system"
}

# Example: Use the data source in another resource
# output "registration_id" {
#   value = data.f5xc_registration.example.id
# }
