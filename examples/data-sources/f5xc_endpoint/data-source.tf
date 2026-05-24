# Endpoint Data Source Example
# Retrieves information about an existing Endpoint

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Endpoint by name
data "f5xc_endpoint" "example" {
  name      = "example-endpoint"
  namespace = "staging"
}

output "endpoint_id" {
  value = data.f5xc_endpoint.example.id
}
