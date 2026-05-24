# Managed Tenant Data Source Example
# Retrieves information about an existing Managed Tenant

# Look up an existing Managed Tenant by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_managed_tenant" "example" {
  name      = "example-managed-tenant"
  namespace = "system"
}

# Example: Use the data source in another resource
# output "managed_tenant_id" {
#   value = data.f5xc_managed_tenant.example.id
# }
