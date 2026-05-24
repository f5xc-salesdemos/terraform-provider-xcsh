# API Discovery Data Source Example
# Retrieves information about an existing API Discovery

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing API Discovery by name
data "f5xc_api_discovery" "example" {
  name      = "example-api-discovery"
  namespace = "staging"
}

output "api_discovery_id" {
  value = data.f5xc_api_discovery.example.id
}
