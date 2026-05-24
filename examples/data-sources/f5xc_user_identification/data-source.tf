# User Identification Data Source Example
# Retrieves information about an existing User Identification

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing User Identification by name
data "f5xc_user_identification" "example" {
  name      = "example-user-identification"
  namespace = "shared"
}

output "user_identification_id" {
  value = data.f5xc_user_identification.example.id
}
