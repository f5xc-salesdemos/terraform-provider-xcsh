# BigIP IRULE Data Source Example
# Retrieves information about an existing BigIP IRULE

# Look up an existing BigIP IRULE by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_bigip_irule" "example" {
  name      = "example-bigip-irule"
  namespace = "system"
}

output "bigip_irule_id" {
  value = data.xcsh_bigip_irule.example.id
}
