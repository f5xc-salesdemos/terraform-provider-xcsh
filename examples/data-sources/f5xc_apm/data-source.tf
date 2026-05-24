# APM Data Source Example
# Retrieves information about an existing APM

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing APM by name
data "f5xc_apm" "example" {
  name      = "example-apm"
  namespace = "staging"
}

output "apm_id" {
  value = data.f5xc_apm.example.id
}
