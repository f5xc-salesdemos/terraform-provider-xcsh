# Monitoring — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "f5xc_app_firewall" "test" {
  name       = "example"
  namespace  = "system"

  default_detection_settings {}
  allow_all_response_codes {}
  monitoring {}
  use_default_blocking_page {}
  default_bot_setting {}
  default_anonymization {}
}