# DNS Compliance Checks Resource Example
# Manages DNS Compliance Checks Specification in a given namespace. If one already exists it will give an error. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic DNS Compliance Checks configuration
resource "f5xc_dns_compliance_checks" "example" {
  name      = "example-dns-compliance-checks"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }
}
