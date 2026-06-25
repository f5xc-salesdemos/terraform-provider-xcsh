# Protocol Policer Data Source Example
# Retrieves information about an existing Protocol Policer

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Protocol Policer by name
data "xcsh_protocol_policer" "example" {
  name      = "example-protocol-policer"
  namespace = "staging"
}

output "protocol_policer_id" {
  value = data.xcsh_protocol_policer.example.id
}
