# Tenant Profile Data Source Example
# Retrieves information about an existing Tenant Profile

# Look up an existing Tenant Profile by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_tenant_profile" "example" {
  name      = "example-tenant-profile"
  namespace = "system"
}

output "tenant_profile_id" {
  value = data.xcsh_tenant_profile.example.id
}
