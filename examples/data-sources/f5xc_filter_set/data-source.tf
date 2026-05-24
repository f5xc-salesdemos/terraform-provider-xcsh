# Filter Set Data Source Example
# Retrieves information about an existing Filter Set

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Filter Set by name
data "f5xc_filter_set" "example" {
  name      = "example-filter-set"
  namespace = "shared"
}

output "filter_set_id" {
  value = data.f5xc_filter_set.example.id
}
