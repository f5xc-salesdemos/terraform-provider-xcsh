# K8S Cluster Role Binding Resource Example
# Manages k8s_cluster_role_binding will create the object in the storage backend for namespace metadata.namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Basic K8S Cluster Role Binding configuration
resource "xcsh_k8s_cluster_role_binding" "example" {
  name      = "example-k8s-cluster-role-binding"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # List of subjects (user, group or service account) to which...
  subjects {
    # Configure subjects settings
  }
  # ServiceAccountType.
  service_account {
    # Configure service_account settings
  }
  # Type establishes a direct reference from one object(the r...
  k8s_cluster_role {
    # Configure k8s_cluster_role settings
  }
}
