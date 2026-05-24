# Data Group Data Source Example
# Retrieves information about an existing Data Group

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Data Group by name
data "f5xc_data_group" "example" {
  name      = "example-data-group"
  namespace = "shared"
}

output "data_group_id" {
  value = data.f5xc_data_group.example.id
}
