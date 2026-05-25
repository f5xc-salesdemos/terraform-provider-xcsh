# MultipleOrigins — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "f5xc_origin_pool" "test" {
  name      = "example"
  namespace = "system"

  port = 443

  origin_servers {
    labels {}
    public_name {
      dns_name = "backend1.example.com"
    }
  }

  origin_servers {
    labels {}
    public_name {
      dns_name = "backend2.example.com"
    }
  }

  no_tls {}
  same_as_endpoint_port {}
}