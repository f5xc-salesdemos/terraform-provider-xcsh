# K8S Cluster Data Source Example
# Retrieves information about an existing K8S Cluster

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing K8S Cluster by name
data "xcsh_k8s_cluster" "example" {
  name      = "example-k8s-cluster"
  namespace = "staging"
}

output "k8s_cluster_id" {
  value = data.xcsh_k8s_cluster.example.id
}
