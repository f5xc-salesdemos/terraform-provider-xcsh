# WithMitigationType — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

terraform {
  required_providers {
    time = {
      source  = "hashicorp/time"
      version = ">= 0.9.0"
    }
  }
}

resource "f5xc_namespace" "test" {
  name = "example"
}

resource "time_sleep" "wait_for_namespace" {
  depends_on      = [f5xc_namespace.test]
  create_duration = "5s"
}

resource "f5xc_malicious_user_mitigation" "test" {
  depends_on  = [time_sleep.wait_for_namespace]
  name        = "example-value"
  namespace   = f5xc_namespace.test.name
  description = "Malicious user mitigation with mitigation type configuration"

  mitigation_type {
    rules {
      threat_level {
        high {}
      }
      mitigation_action {
        block_temporarily {}
      }
    }
  }
}