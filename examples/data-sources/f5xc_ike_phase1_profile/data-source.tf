# Ike Phase1 Profile Data Source Example
# Retrieves information about an existing Ike Phase1 Profile

# Look up an existing Ike Phase1 Profile by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_ike_phase1_profile" "example" {
  name      = "example-ike-phase1-profile"
  namespace = "system"
}

output "ike_phase1_profile_id" {
  value = data.f5xc_ike_phase1_profile.example.id
}
