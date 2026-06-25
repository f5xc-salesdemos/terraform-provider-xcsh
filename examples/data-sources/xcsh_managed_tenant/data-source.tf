# Managed Tenant Data Source Example
# Retrieves information about an existing Managed Tenant

# Look up an existing Managed Tenant by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_managed_tenant" "example" {
  name      = "example-managed-tenant"
  namespace = "system"
}

output "managed_tenant_id" {
  value = data.xcsh_managed_tenant.example.id
}
