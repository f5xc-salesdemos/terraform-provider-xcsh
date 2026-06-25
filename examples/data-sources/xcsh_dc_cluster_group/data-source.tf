# Dc Cluster Group Data Source Example
# Retrieves information about an existing Dc Cluster Group

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Dc Cluster Group by name
data "xcsh_dc_cluster_group" "example" {
  name      = "example-dc-cluster-group"
  namespace = "staging"
}

output "dc_cluster_group_id" {
  value = data.xcsh_dc_cluster_group.example.id
}
