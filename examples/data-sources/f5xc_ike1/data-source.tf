# Ike1 Data Source Example
# Retrieves information about an existing Ike1

# Look up an existing Ike1 by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_ike1" "example" {
  name      = "example-ike1"
  namespace = "system"
}

# Example: Use the data source in another resource
# output "ike1_id" {
#   value = data.f5xc_ike1.example.id
# }
