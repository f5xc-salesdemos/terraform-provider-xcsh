# Workload Flavor Data Source Example
# Retrieves information about an existing Workload Flavor

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Workload Flavor by name
data "xcsh_workload_flavor" "example" {
  name      = "example-workload-flavor"
  namespace = "staging"
}

output "workload_flavor_id" {
  value = data.xcsh_workload_flavor.example.id
}
