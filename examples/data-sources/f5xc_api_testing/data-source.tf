# API Testing Data Source Example
# Retrieves information about an existing API Testing

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing API Testing by name
data "f5xc_api_testing" "example" {
  name      = "example-api-testing"
  namespace = "staging"
}

output "api_testing_id" {
  value = data.f5xc_api_testing.example.id
}
