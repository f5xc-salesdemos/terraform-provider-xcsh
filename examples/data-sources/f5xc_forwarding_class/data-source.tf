# Forwarding Class Data Source Example
# Retrieves information about an existing Forwarding Class

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Forwarding Class by name
data "f5xc_forwarding_class" "example" {
  name      = "example-forwarding-class"
  namespace = "staging"
}

output "forwarding_class_id" {
  value = data.f5xc_forwarding_class.example.id
}
