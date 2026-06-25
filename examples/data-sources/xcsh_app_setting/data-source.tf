# App Setting Data Source Example
# Retrieves information about an existing App Setting

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing App Setting by name
data "xcsh_app_setting" "example" {
  name      = "example-app-setting"
  namespace = "staging"
}

output "app_setting_id" {
  value = data.xcsh_app_setting.example.id
}
