# Protocol Inspection Data Source Example
# Retrieves information about an existing Protocol Inspection

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Protocol Inspection by name
data "f5xc_protocol_inspection" "example" {
  name      = "example-protocol-inspection"
  namespace = "staging"
}

output "protocol_inspection_id" {
  value = data.f5xc_protocol_inspection.example.id
}
