# Ike2 Data Source Example
# Retrieves information about an existing Ike2

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Ike2 by name
data "f5xc_ike2" "example" {
  name      = "example-ike2"
  namespace = "staging"
}

output "ike2_id" {
  value = data.f5xc_ike2.example.id
}
