# Trusted CA List Resource Example
# Manages a Trusted CA List resource in F5 Distributed Cloud for trusted certificate authority list management.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source = "f5xc-salesdemos/f5xc"
    }
  }
}

# Basic Trusted CA List configuration
resource "f5xc_trusted_ca_list" "example" {
  name      = "example-trusted-ca-list"
  namespace = "shared"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Trusted CA List configuration
  trusted_ca_url = "string:///LS0tLS1CRUdJTi..."
}
