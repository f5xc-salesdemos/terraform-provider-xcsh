# Cloud Elastic IP Resource Example
# Manages Cloud Elastic IP creates Cloud Elastic IP object Object is attached to a site. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Basic Cloud Elastic IP configuration
resource "xcsh_cloud_elastic_ip" "example" {
  name      = "example-cloud-elastic-ip"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Site to which this cloud elastic IP object is attached .
  site_ref {
    # Configure site_ref settings
  }
}
