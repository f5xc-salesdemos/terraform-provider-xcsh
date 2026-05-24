# Tpm Manager Data Source Example
# Retrieves information about an existing Tpm Manager

# Look up an existing Tpm Manager by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_tpm_manager" "example" {
  name      = "example-tpm-manager"
  namespace = "system"
}

output "tpm_manager_id" {
  value = data.f5xc_tpm_manager.example.id
}
