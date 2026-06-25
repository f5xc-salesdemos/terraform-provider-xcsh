# Authentication Data Source Example
# Retrieves information about an existing Authentication

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Authentication by name
data "xcsh_authentication" "example" {
  name      = "example-authentication"
  namespace = "staging"
}

output "authentication_id" {
  value = data.xcsh_authentication.example.id
}
