# IKE Phase1 Profile Data Source Example
# Retrieves information about an existing IKE Phase1 Profile

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing IKE Phase1 Profile by name
data "f5xc_ike_phase1_profile" "example" {
  name      = "example-ike-phase1-profile"
  namespace = "staging"
}

output "ike_phase1_profile_id" {
  value = data.f5xc_ike_phase1_profile.example.id
}
