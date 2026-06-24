# WithRateLimit — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "xcsh_http_loadbalancer" "test" {
  name      = "example"
  namespace = "system"
  domains   = ["test.example.com"]

  http {
    port = 80
  }

  rate_limit {
    rate_limiter {
      total_number     = 100
      unit             = "MINUTE"
      burst_multiplier = 10
    }
    no_ip_allowed_list {}
  }

  advertise_on_public_default_vip {}
}