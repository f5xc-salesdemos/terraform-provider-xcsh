# LabelsUpdate — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "xcsh_http_loadbalancer" "test" {
  name      = "example"
  namespace = "system"

  labels = {
    environment = "example-value"
    managed_by  = "terraform"
  }

  domains = ["test.example.com"]

  http {
    port = 80
  }

  advertise_on_public_default_vip {}
}