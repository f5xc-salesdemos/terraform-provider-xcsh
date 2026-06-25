# Blocking — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "xcsh_app_firewall" "test" {
  name      = "example"
  namespace = "system"

  # Use default detection settings
  default_detection_settings {}

  # Allow all response codes
  allow_all_response_codes {}

  # Blocking mode - actively block malicious requests
  blocking {}

  # Use default blocking page
  use_default_blocking_page {}

  # Use default bot settings
  default_bot_setting {}

  # Use default anonymization
  default_anonymization {}
}