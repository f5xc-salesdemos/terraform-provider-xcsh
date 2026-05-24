# Route Data Source Example
# Retrieves information about an existing Route

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Route by name
data "f5xc_route" "example" {
  name      = "example-route"
  namespace = "system"
}

output "route_id" {
  value = data.f5xc_route.example.id
}
