# Dc Cluster Group Resource Example
# Manages DC Cluster group in given namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Dc Cluster Group configuration
resource "xcsh_dc_cluster_group" "example" {
  name      = "example-dc-cluster-group"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # DC Cluster Group Mesh Type. Details of DC Cluster Group M...
  type {
    # Configure type settings
  }
  # Enable this option
  control_and_data_plane_mesh {
    # Configure control_and_data_plane_mesh settings
  }
  # Enable this option
  data_plane_mesh {
    # Configure data_plane_mesh settings
  }
}
