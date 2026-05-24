# Dc Cluster Group Data Source Example
# Retrieves information about an existing Dc Cluster Group

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Dc Cluster Group by name
data "f5xc_dc_cluster_group" "example" {
  name      = "example-dc-cluster-group"
  namespace = "staging"
}

output "dc_cluster_group_id" {
  value = data.f5xc_dc_cluster_group.example.id
}
