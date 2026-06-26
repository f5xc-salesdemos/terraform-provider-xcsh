# APM Resource Example
# Manages new APM as a service with configured parameters. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic APM configuration
resource "xcsh_apm" "example" {
  name      = "example-apm"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # [OneOf: aws_site_type_choice, baremetal_site_type_choice]...
  aws_site_type_choice {
    # Configure aws_site_type_choice settings
  }
  # Virtual F5 BIG-IP configuration for AWS TGW Site using BI...
  apm_aws_site {
    # Configure apm_aws_site settings
  }
  # SecretType is used in an object to indicate a sensitive/c...
  admin_password {
    # Configure admin_password settings
  }
}
