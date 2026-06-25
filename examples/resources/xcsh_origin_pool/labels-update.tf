# LabelsUpdate — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "xcsh_origin_pool" "test" {
  name      = "example"
  namespace = "system"

  port = 443

  labels = {
    environment = "example-value"
  }

  origin_servers {
    labels {} # API returns this even if not set
    public_name {
      dns_name = "example.com"
    }
  }

  no_tls {}
  same_as_endpoint_port {}
}