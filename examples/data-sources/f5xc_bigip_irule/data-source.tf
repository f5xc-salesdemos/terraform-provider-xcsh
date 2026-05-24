# BigIP IRULE Data Source Example
# Retrieves information about an existing BigIP IRULE

# Look up an existing BigIP IRULE by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_bigip_irule" "example" {
  name      = "example-bigip-irule"
  namespace = "system"
}

# Example: Use the data source in another resource
# output "bigip_irule_id" {
#   value = data.f5xc_bigip_irule.example.id
# }
