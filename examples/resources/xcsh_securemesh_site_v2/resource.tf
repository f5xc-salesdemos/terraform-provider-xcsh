# Securemesh Site V2 Resource Example
# Manages a Securemesh Site V2 resource in F5 Distributed Cloud for deploying secure mesh edge sites with enhanced security and networking features.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Basic Securemesh Site V2 configuration
resource "xcsh_securemesh_site_v2" "example" {
  name      = "example-securemesh-site-v2"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Optimize the site for L3 or L7 traffic processing. L7 opt...
  performance_enhancement_mode {
    # Configure performance_enhancement_mode settings
  }
  # Configuration parameter for perf mode l3 enhanced.
  perf_mode_l3_enhanced {
    # Configure perf_mode_l3_enhanced settings
  }
  # Enable this option
  jumbo {
    # Configure jumbo settings
  }
}
