# DetectionSettings — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "f5xc_app_firewall" "test" {
  name      = "example"
  namespace = "system"

  allow_all_response_codes {}
  blocking {}
  use_default_blocking_page {}
  default_bot_setting {}
  default_anonymization {}

  detection_settings {
    default_violation_settings {}
    default_bot_setting {}
    enable_suppression {}
    enable_threat_campaigns {}
    signature_selection_setting {
      high_medium_accuracy_signatures {}
      default_attack_type_settings {}
    }
  }
}