# K8S Cluster Role Binding Resource Example
# Manages k8s_cluster_role_binding will create the object in the storage backend for namespace metadata.namespace in F5 Distributed Cloud.

# Basic K8S Cluster Role Binding configuration
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


resource "f5xc_k8s_cluster_role_binding" "example" {
  name      = "example-k8s-cluster-role-binding"
  namespace = "system"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Object reference. This type establishes a direct reference...
  k8s_cluster_role {
    # Configure k8s_cluster_role settings
  }
  # Subjects. List of subjects (user, group or service account...
  subjects {
    # Configure subjects settings
  }
  # ServiceAccountType.
  service_account {
    # Configure service_account settings
  }
}
