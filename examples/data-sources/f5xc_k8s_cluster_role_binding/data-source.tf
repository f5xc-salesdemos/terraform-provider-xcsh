# K8S Cluster Role Binding Data Source Example
# Retrieves information about an existing K8S Cluster Role Binding

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing K8S Cluster Role Binding by name
data "f5xc_k8s_cluster_role_binding" "example" {
  name      = "example-k8s-cluster-role-binding"
  namespace = "staging"
}

output "k8s_cluster_role_binding_id" {
  value = data.f5xc_k8s_cluster_role_binding.example.id
}
