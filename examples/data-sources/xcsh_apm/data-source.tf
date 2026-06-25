# APM Data Source Example
# Retrieves information about an existing APM

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing APM by name
data "xcsh_apm" "example" {
  name      = "example-apm"
  namespace = "staging"
}

output "apm_id" {
  value = data.xcsh_apm.example.id
}
