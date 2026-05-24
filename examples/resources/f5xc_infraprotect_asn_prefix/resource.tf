# Infraprotect Asn Prefix Resource Example
# Manages DDoS transit Prefix in F5 Distributed Cloud.

# Basic Infraprotect Asn Prefix configuration
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


resource "f5xc_infraprotect_asn_prefix" "example" {
  name      = "example-infraprotect-asn-prefix"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Object reference. This type establishes a direct reference...
  asn {
    # Configure asn settings
  }
}
