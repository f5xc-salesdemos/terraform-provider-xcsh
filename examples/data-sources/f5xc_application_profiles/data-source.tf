# Application Profiles Data Source Example
# Retrieves information about an existing Application Profiles

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Application Profiles by name
data "f5xc_application_profiles" "example" {
  name      = "example-application-profiles"
  namespace = "staging"
}

output "application_profiles_id" {
  value = data.f5xc_application_profiles.example.id
}
