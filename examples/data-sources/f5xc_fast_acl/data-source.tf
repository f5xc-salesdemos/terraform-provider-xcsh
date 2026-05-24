# Fast ACL Data Source Example
# Retrieves information about an existing Fast ACL

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Fast ACL by name
data "f5xc_fast_acl" "example" {
  name      = "example-fast-acl"
  namespace = "system"
}

output "fast_acl_id" {
  value = data.f5xc_fast_acl.example.id
}
