# Secret Policy Data Source Example
# Retrieves information about an existing Secret Policy

# Look up an existing Secret Policy by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_secret_policy" "example" {
  name      = "example-secret-policy"
  namespace = "system"
}

output "secret_policy_id" {
  value = data.xcsh_secret_policy.example.id
}
