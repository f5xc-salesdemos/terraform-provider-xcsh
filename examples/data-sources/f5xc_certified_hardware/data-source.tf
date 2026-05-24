# Certified Hardware Data Source Example
# Retrieves information about an existing Certified Hardware

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Certified Hardware by name
data "f5xc_certified_hardware" "example" {
  name      = "example-certified-hardware"
  namespace = "staging"
}

output "certified_hardware_id" {
  value = data.f5xc_certified_hardware.example.id
}
