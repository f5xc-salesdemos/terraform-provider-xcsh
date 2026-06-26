# Fast ACL Data Source Example
# Retrieves information about an existing Fast ACL

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Fast ACL by name
data "xcsh_fast_acl" "example" {
  name      = "example-fast-acl"
  namespace = "staging"
}

output "fast_acl_id" {
  value = data.xcsh_fast_acl.example.id
}
