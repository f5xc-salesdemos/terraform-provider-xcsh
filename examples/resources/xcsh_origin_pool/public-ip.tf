# PublicIp — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "xcsh_origin_pool" "test" {
  name      = "example"
  namespace = "system"

  port = 8080

  origin_servers {
    labels {}
    public_ip {
      ip = "93.184.216.34"
    }
  }

  no_tls {}
  same_as_endpoint_port {}
}