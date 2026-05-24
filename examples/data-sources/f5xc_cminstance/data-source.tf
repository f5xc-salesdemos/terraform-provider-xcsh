# Cminstance Data Source Example
# Retrieves information about an existing Cminstance

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Cminstance by name
data "f5xc_cminstance" "example" {
  name      = "example-cminstance"
  namespace = "staging"
}

output "cminstance_id" {
  value = data.f5xc_cminstance.example.id
}
