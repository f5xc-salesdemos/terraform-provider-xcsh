# Allowed Tenant Data Source Example
# Retrieves information about an existing Allowed Tenant

# Look up an existing Allowed Tenant by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_allowed_tenant" "example" {
  name      = "example-allowed-tenant"
  namespace = "system"
}

output "allowed_tenant_id" {
  value = data.xcsh_allowed_tenant.example.id
}
