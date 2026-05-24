# Protocol Inspection Resource Example
# Manages Protocol Inspection Specification in a given namespace. If one already exists it will give an error. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Protocol Inspection configuration
resource "f5xc_protocol_inspection" "example" {
  name      = "example-protocol-inspection"
  namespace = "shared"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Enable Disable Compliance Checks Choice.
  enable_disable_compliance_checks {
    # Configure enable_disable_compliance_checks settings
  }
  # Configuration parameter for disable compliance checks.
  disable_compliance_checks {
    # Configure disable_compliance_checks settings
  }
  # Type establishes a direct reference from one object(the r...
  enable_compliance_checks {
    # Configure enable_compliance_checks settings
  }
}

# The following optional fields have server-applied defaults and can be omitted:
# - action
