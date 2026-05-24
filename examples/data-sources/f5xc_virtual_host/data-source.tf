# Virtual Host Data Source Example
# Retrieves information about an existing Virtual Host

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Virtual Host by name
data "f5xc_virtual_host" "example" {
  name      = "example-virtual-host"
  namespace = "staging"
}

output "virtual_host_id" {
  value = data.f5xc_virtual_host.example.id
}
