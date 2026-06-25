# Role Data Source Example
# Retrieves information about an existing Role

# Look up an existing Role by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_role" "example" {
  name      = "example-role"
  namespace = "system"
}

output "role_id" {
  value = data.xcsh_role.example.id
}
