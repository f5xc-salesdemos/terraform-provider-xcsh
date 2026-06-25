# AllowedResponseCodes — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "xcsh_app_firewall" "test" {
  name      = "example"
  namespace = "system"

  default_detection_settings {}
  blocking {}
  use_default_blocking_page {}
  default_bot_setting {}
  default_anonymization {}

  allowed_response_codes {
    response_code = [200, 204, 301, 302]
  }
}