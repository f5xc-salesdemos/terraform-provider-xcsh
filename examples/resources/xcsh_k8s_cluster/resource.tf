# K8S Cluster Resource Example
# Manages k8s_cluster will create the object in the storage backend for namespace metadata.namespace. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic K8S Cluster configuration
resource "xcsh_k8s_cluster" "example" {
  name      = "example-k8s-cluster"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Kubernetes Cluster configuration
  # Use custom local domain
  use_custom_cluster_role_bindings {
    cluster_role_bindings {
      name      = "admin-binding"
      namespace = "staging"
    }
  }

  cluster_wide_app_list {
    cluster_wide_apps {
      name      = "nginx-ingress"
      namespace = "staging"
    }
  }

  local_access_config {
    local_domain = "cluster.local"
    default_port {}
  }

  global_access_enable {}
}

# The following optional fields have server-applied defaults and can be omitted:
# - cluster_scoped_access_deny
# - no_cluster_wide_apps
# - no_global_access
# - no_insecure_registries
# - no_local_access
# - use_default_cluster_role_bindings
# - use_default_cluster_roles
# - use_default_psp
# - vk8s_namespace_access_deny
