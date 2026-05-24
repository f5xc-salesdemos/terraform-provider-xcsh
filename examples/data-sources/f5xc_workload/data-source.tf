# Workload Data Source Example
# Retrieves information about an existing Workload

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Workload by name
data "f5xc_workload" "example" {
  name      = "example-workload"
  namespace = "staging"
}

output "workload_id" {
  value = data.f5xc_workload.example.id
}
