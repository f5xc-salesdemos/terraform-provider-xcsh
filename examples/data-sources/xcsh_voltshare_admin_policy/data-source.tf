# Voltshare Admin Policy Data Source Example
# Retrieves information about an existing Voltshare Admin Policy

# Look up an existing Voltshare Admin Policy by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_voltshare_admin_policy" "example" {
  name      = "example-voltshare-admin-policy"
  namespace = "system"
}

output "voltshare_admin_policy_id" {
  value = data.xcsh_voltshare_admin_policy.example.id
}
