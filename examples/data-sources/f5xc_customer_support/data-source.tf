# Customer Support Data Source Example
# Retrieves information about an existing Customer Support

# Look up an existing Customer Support by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_customer_support" "example" {
  name      = "example-customer-support"
  namespace = "system"
}

# Example: Use the data source in another resource
# output "customer_support_id" {
#   value = data.f5xc_customer_support.example.id
# }
