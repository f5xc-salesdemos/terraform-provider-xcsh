# Quota Data Source Example
# Retrieves information about an existing Quota

# Look up an existing Quota by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_quota" "example" {
  name      = "example-quota"
  namespace = "system"
}

output "quota_id" {
  value = data.xcsh_quota.example.id
}
