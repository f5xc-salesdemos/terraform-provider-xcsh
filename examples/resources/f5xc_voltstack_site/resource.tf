# Voltstack Site Resource Example
# Manages a Voltstack Site resource in F5 Distributed Cloud for deploying Volterra stack sites for edge computing.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Voltstack Site configuration
resource "f5xc_voltstack_site" "example" {
  name      = "example-voltstack-site"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  volterra_certified_hw = "kvm-voltstack-combo"
  worker_nodes          = []
  address               = "123 Main St, Example City, EX 12345"

  k8s_cluster {
    name      = "example-k8s-cluster"
    namespace = "staging"
  }
}
