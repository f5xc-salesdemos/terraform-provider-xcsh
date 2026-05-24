# K8s Cluster Data Source Example
# Retrieves information about an existing K8s Cluster

# Look up an existing K8s Cluster by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_k8s_cluster" "example" {
  name      = "example-k8s-cluster"
  namespace = "system"
}

# Example: Use the data source in another resource
# output "k8s_cluster_id" {
#   value = data.f5xc_k8s_cluster.example.id
# }
