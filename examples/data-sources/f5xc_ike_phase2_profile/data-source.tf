# Ike Phase2 Profile Data Source Example
# Retrieves information about an existing Ike Phase2 Profile

# Look up an existing Ike Phase2 Profile by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_ike_phase2_profile" "example" {
  name      = "example-ike-phase2-profile"
  namespace = "system"
}

# Example: Use the data source in another resource
# output "ike_phase2_profile_id" {
#   value = data.f5xc_ike_phase2_profile.example.id
# }
