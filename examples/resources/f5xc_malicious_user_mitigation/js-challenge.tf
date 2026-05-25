# JsChallenge — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "f5xc_malicious_user_mitigation" "test" {
  name      = "example"
  namespace = "system"

  mitigation_type {
    rules {
      threat_level {
        medium {}
      }
      mitigation_action {
        javascript_challenge {}
      }
    }
  }
}