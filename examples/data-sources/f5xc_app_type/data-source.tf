# App Type Data Source Example
# Retrieves information about an existing App Type

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing App Type by name
data "f5xc_app_type" "example" {
  name      = "example-app-type"
  namespace = "staging"
}

output "app_type_id" {
  value = data.f5xc_app_type.example.id
}
