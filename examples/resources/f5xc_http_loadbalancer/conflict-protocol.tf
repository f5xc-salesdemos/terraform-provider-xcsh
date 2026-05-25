# ConflictProtocol — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "f5xc_http_loadbalancer" "test" {
  name      = "example"
  namespace = "system"
  domains   = ["test.example.com"]

  http {
    port = 80
  }

  https_auto_cert {
    add_hsts = false
    no_mtls {}
    default_header {}
    enable_path_normalize {}
    non_default_loadbalancer {}
  }

  advertise_on_public_default_vip {}
}