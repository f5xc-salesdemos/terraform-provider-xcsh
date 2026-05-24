# Cloud Link Resource Example
# Manages new CloudLink with configured parameters. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Cloud Link configuration
resource "f5xc_cloud_link" "example" {
  name      = "example-cloud-link"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # [OneOf: aws, gcp] Amazon Web Services(AWS) CloudLink Prov...
  aws {
    # Configure aws settings
  }
  # Type establishes a direct reference from one object(the r...
  aws_cred {
    # Configure aws_cred settings
  }
  # Bring Your Own Connections. List of Bring You Own Connect...
  byoc {
    # Configure byoc settings
  }
}
