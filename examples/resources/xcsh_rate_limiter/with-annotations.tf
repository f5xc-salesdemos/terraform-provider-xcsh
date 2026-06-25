# WithAnnotations — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "xcsh_rate_limiter" "test" {
  name      = "example"
  namespace = "system"

  annotations = {
    example-key = "example-value"
  }
}