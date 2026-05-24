# Malicious User Mitigation Data Source Example
# Retrieves information about an existing Malicious User Mitigation

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Malicious User Mitigation by name
data "f5xc_malicious_user_mitigation" "example" {
  name      = "example-malicious-user-mitigation"
  namespace = "staging"
}

output "malicious_user_mitigation_id" {
  value = data.f5xc_malicious_user_mitigation.example.id
}
