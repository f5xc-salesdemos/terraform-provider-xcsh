# HttpStatusCodes — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "xcsh_healthcheck" "test" {
  name      = "example"
  namespace = "system"

  healthy_threshold   = 1
  unhealthy_threshold = 2
  timeout             = 3
  interval            = 5

  http_health_check {
    path                  = "example-value"
    host_header           = "example.com"
    expected_status_codes = ["200", "201", "204"]
  }
}