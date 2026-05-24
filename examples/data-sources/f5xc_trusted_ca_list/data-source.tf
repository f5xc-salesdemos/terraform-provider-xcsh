# Trusted CA List Data Source Example
# Retrieves information about an existing Trusted CA List

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Trusted CA List by name
data "f5xc_trusted_ca_list" "example" {
  name      = "example-trusted-ca-list"
  namespace = "shared"
}

output "trusted_ca_list_id" {
  value = data.f5xc_trusted_ca_list.example.id
}
