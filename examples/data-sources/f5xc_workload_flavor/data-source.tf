# Workload Flavor Data Source Example
# Retrieves information about an existing Workload Flavor

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Workload Flavor by name
data "f5xc_workload_flavor" "example" {
  name      = "example-workload-flavor"
  namespace = "staging"
}

output "workload_flavor_id" {
  value = data.f5xc_workload_flavor.example.id
}
