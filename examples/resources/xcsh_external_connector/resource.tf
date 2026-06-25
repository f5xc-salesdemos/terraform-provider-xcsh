# External Connector Resource Example
# Manages a External Connector resource in F5 Distributed Cloud for external_connector configuration specification. configuration.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Basic External Connector configuration
resource "xcsh_external_connector" "example" {
  name      = "example-external-connector"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Type establishes a direct reference from one object(the r...
  ce_site_reference {
    # Configure ce_site_reference settings
  }
  # GRE. External Connector with GRE tunnel.
  gre {
    # Configure gre settings
  }
  # X-displayName: 'GRE Tunnel Parameters' GRE configuration ...
  gre_parameters {
    # Configure gre_parameters settings
  }
}
