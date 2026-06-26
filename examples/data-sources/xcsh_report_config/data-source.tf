# Report Config Data Source Example
# Retrieves information about an existing Report Config

# Look up an existing Report Config by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_report_config" "example" {
  name      = "example-report-config"
  namespace = "system"
}

output "report_config_id" {
  value = data.xcsh_report_config.example.id
}
