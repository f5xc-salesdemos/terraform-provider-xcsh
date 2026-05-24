# Role Data Source Example
# Retrieves information about an existing Role

# Look up an existing Role by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_role" "example" {
  name      = "example-role"
  namespace = "system"
}

output "role_id" {
  value = data.f5xc_role.example.id
}
