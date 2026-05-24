# Site Mesh Group Resource Example
# Manages Site Mesh Group in system namespace of user. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Site Mesh Group configuration
resource "f5xc_site_mesh_group" "example" {
  name      = "example-site-mesh-group"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Site Mesh Group configuration
  type = "SITE_MESH_GROUP_TYPE_FULL_MESH"

  # Control and data plane settings
  full_mesh {
    control_and_data_plane_mesh {}
  }

  # Hub status
  hub {}

  # Virtual site reference
  virtual_site {
    name      = "example-virtual-site"
    namespace = "staging"
  }
}
