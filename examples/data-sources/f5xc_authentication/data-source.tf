# Authentication Data Source Example
# Retrieves information about an existing Authentication

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Authentication by name
data "f5xc_authentication" "example" {
  name      = "example-authentication"
  namespace = "staging"
}

output "authentication_id" {
  value = data.f5xc_authentication.example.id
}
