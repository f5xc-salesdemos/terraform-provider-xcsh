# Rate Limiter Policy Data Source Example
# Retrieves information about an existing Rate Limiter Policy

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Rate Limiter Policy by name
data "f5xc_rate_limiter_policy" "example" {
  name      = "example-rate-limiter-policy"
  namespace = "shared"
}

output "rate_limiter_policy_id" {
  value = data.f5xc_rate_limiter_policy.example.id
}
