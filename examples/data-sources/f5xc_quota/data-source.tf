# Quota Data Source Example
# Retrieves information about an existing Quota

# Look up an existing Quota by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_quota" "example" {
  name      = "example-quota"
  namespace = "system"
}

output "quota_id" {
  value = data.f5xc_quota.example.id
}
