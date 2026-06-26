# Virtual Site Resource Example
# Manages virtual site object in given namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Virtual Site configuration
resource "xcsh_virtual_site" "example" {
  name      = "example-virtual-site"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Virtual Site configuration
  site_type = "CUSTOMER_EDGE"

  # Site selector expression
  site_selector {
    expressions = ["region in (us-west-2, us-east-1)"]
  }
}
