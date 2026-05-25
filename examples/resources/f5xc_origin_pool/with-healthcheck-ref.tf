# WithHealthcheckRef — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "f5xc_healthcheck" "test" {
  name      = "example"
  namespace = "system"

  healthy_threshold   = 3
  unhealthy_threshold = 1
  timeout             = 3
  interval            = 15

  tcp_health_check {}
}

resource "f5xc_origin_pool" "test" {
  name      = "example"
  namespace = "system"

  port = 443

  origin_servers {
    labels {}
    public_name {
      dns_name = "example.com"
    }
  }

  healthcheck {
    name      = f5xc_healthcheck.test.name
    namespace = f5xc_healthcheck.test.namespace
  }

  no_tls {}
  same_as_endpoint_port {}
}