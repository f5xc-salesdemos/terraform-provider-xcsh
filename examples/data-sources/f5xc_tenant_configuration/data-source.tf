# Tenant Configuration Data Source Example
# Retrieves information about an existing Tenant Configuration

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Tenant Configuration by name
data "f5xc_tenant_configuration" "example" {
  name      = "example-tenant-configuration"
  namespace = "staging"
}

output "tenant_configuration_id" {
  value = data.f5xc_tenant_configuration.example.id
}
