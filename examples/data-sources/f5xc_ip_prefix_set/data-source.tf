# IP Prefix Set Data Source Example
# Retrieves information about an existing IP Prefix Set

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing IP Prefix Set by name
data "f5xc_ip_prefix_set" "example" {
  name      = "example-ip-prefix-set"
  namespace = "shared"
}

output "ip_prefix_set_id" {
  value = data.f5xc_ip_prefix_set.example.id
}
