# Virtual K8s Data Source Example
# Retrieves information about an existing Virtual K8s

# Look up an existing Virtual K8s by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_virtual_k8s" "example" {
  name      = "example-virtual-k8s"
  namespace = "system"
}

output "virtual_k8s_id" {
  value = data.f5xc_virtual_k8s.example.id
}
