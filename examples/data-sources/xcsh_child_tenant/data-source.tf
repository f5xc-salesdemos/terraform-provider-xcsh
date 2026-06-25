# Child Tenant Data Source Example
# Retrieves information about an existing Child Tenant

# Look up an existing Child Tenant by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_child_tenant" "example" {
  name      = "example-child-tenant"
  namespace = "system"
}

output "child_tenant_id" {
  value = data.xcsh_child_tenant.example.id
}
