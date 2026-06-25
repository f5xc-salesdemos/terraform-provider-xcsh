# Child Tenant Manager Data Source Example
# Retrieves information about an existing Child Tenant Manager

# Look up an existing Child Tenant Manager by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_child_tenant_manager" "example" {
  name      = "example-child-tenant-manager"
  namespace = "system"
}

output "child_tenant_manager_id" {
  value = data.xcsh_child_tenant_manager.example.id
}
