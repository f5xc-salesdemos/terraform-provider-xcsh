# Contact Data Source Example
# Retrieves information about an existing Contact

# Look up an existing Contact by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_contact" "example" {
  name      = "example-contact"
  namespace = "system"
}

output "contact_id" {
  value = data.f5xc_contact.example.id
}
