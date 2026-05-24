# Voltshare Admin Policy Data Source Example
# Retrieves information about an existing Voltshare Admin Policy

# Look up an existing Voltshare Admin Policy by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_voltshare_admin_policy" "example" {
  name      = "example-voltshare-admin-policy"
  namespace = "system"
}

output "voltshare_admin_policy_id" {
  value = data.f5xc_voltshare_admin_policy.example.id
}
