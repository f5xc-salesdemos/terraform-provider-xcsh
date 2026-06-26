# Tpm Category Data Source Example
# Retrieves information about an existing Tpm Category

# Look up an existing Tpm Category by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_tpm_category" "example" {
  name      = "example-tpm-category"
  namespace = "system"
}

output "tpm_category_id" {
  value = data.xcsh_tpm_category.example.id
}
