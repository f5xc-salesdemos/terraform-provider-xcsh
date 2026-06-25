# TlsFingerprint — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "xcsh_user_identification" "test" {
  name      = "example"
  namespace = "system"

  rules {
    tls_fingerprint {}
  }
}