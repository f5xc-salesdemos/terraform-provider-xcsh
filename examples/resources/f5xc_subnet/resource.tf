# Subnet Resource Example
# Manages a Subnet resource in F5 Distributed Cloud for subnet object contains configuration for an interface of a vm/pod. it is created in user or shared namespace. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source = "f5xc-salesdemos/f5xc"
    }
  }
}

# Basic Subnet configuration
resource "f5xc_subnet" "example" {
  name      = "example-subnet"
  namespace = "system"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Configure subnet parameters per site .
  site_subnet_params {
    # Configure site_subnet_params settings
  }
  # Enable this option
  dhcp {
    # Configure dhcp settings
  }
  # Type establishes a direct reference from one object(the r...
  site {
    # Configure site settings
  }
}
