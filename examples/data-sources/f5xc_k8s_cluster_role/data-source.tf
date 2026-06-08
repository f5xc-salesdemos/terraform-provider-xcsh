# K8S Cluster Role Data Source Example
# Retrieves information about an existing K8S Cluster Role

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing K8S Cluster Role by name
data "f5xc_k8s_cluster_role" "example" {
  name      = "example-k8s-cluster-role"
  namespace = "staging"
}

output "k8s_cluster_role_id" {
  value = data.f5xc_k8s_cluster_role.example.id
}
