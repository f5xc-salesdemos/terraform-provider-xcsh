# Address Allocator Resource Example
# Manages Address Allocator will create an address allocator object in 'system' namespace of the user. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Address Allocator configuration
resource "xcsh_address_allocator" "example" {
  name      = "example-address-allocator"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Decides the scheme to be used to allocate addresses from ...
  address_allocation_scheme {
    # Configure address_allocation_scheme settings
  }
}
