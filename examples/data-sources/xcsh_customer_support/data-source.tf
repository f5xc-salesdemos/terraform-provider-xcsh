# Customer Support Data Source Example
# Retrieves information about an existing Customer Support

# Look up an existing Customer Support by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_customer_support" "example" {
  name      = "example-customer-support"
  namespace = "system"
}

output "customer_support_id" {
  value = data.xcsh_customer_support.example.id
}
