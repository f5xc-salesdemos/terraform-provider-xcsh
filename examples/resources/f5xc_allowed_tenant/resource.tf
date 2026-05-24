# Allowed Tenant Resource Example
# Manages allowed_tenant config instance. Name of the object is name of the tenant that is allowed to manage. in F5 Distributed Cloud.

# Basic Allowed Tenant configuration
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


resource "f5xc_allowed_tenant" "example" {
  name      = "example-allowed-tenant"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Allowed Groups. List of references to allowed user_group ...
  allowed_groups {
    # Configure allowed_groups settings
  }
}
