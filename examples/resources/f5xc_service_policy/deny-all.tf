# DenyAll — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "f5xc_service_policy" "test" {
  name       = "example"
  namespace  = "system"

  # Deny all requests
  deny_all_requests {}

  # Apply to any server
  any_server {}
}