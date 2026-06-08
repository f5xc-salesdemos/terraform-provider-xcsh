# K8S Pod Security Admission Data Source Example
# Retrieves information about an existing K8S Pod Security Admission

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing K8S Pod Security Admission by name
data "f5xc_k8s_pod_security_admission" "example" {
  name      = "example-k8s-pod-security-admission"
  namespace = "staging"
}

output "k8s_pod_security_admission_id" {
  value = data.f5xc_k8s_pod_security_admission.example.id
}
