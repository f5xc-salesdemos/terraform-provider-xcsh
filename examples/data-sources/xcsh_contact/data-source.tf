# Contact Data Source Example
# Retrieves information about an existing Contact

# Look up an existing Contact by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_contact" "example" {
  name      = "example-contact"
  namespace = "system"
}

output "contact_id" {
  value = data.xcsh_contact.example.id
}
