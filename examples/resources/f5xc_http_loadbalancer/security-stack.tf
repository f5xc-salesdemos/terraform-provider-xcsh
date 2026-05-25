# SecurityStack — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "f5xc_healthcheck" "test" {
  name      = "example"
  namespace = "system"

  healthy_threshold   = 3
  unhealthy_threshold = 1
  timeout             = 3
  interval            = 15

  http_health_check {
    path        = "/health"
    host_header = "example.com"
  }
}

resource "f5xc_origin_pool" "test" {
  name      = "example"
  namespace = "system"
  port      = 443

  origin_servers {
    labels {}
    public_name {
      dns_name = "example.com"
    }
  }

  healthcheck {
    name      = f5xc_healthcheck.test.name
    namespace = "system"
  }

  no_tls {}
  same_as_endpoint_port {}
}

resource "f5xc_app_firewall" "test" {
  name      = "example"
  namespace = "system"

  default_detection_settings {}
  allow_all_response_codes {}
  blocking {}
  use_default_blocking_page {}
  default_bot_setting {}
  default_anonymization {}
}

resource "f5xc_http_loadbalancer" "test" {
  name      = "example"
  namespace = "system"

  domains = ["test.example.com"]

  http {
    port = 80
  }

  default_route_pools {
    pool {
      name      = f5xc_origin_pool.test.name
      namespace = "system"
    }
    weight   = 1
    priority = 1
  }

  app_firewall {
    name      = f5xc_app_firewall.test.name
    namespace = "system"
  }

  enable_malicious_user_detection {}
  enable_threat_mesh {}

  advertise_on_public_default_vip {}
}