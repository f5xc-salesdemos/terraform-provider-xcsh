terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

provider "xcsh" {}

resource "xcsh_healthcheck" "httpbin" {
  name      = "httpbin-health"
  namespace = "demo"

  http_health_check {
    use_origin_server_name {}
    path = "/get"
  }

  healthy_threshold   = 3
  unhealthy_threshold = 3
  interval            = 15
  timeout             = 5
}

resource "xcsh_origin_pool" "httpbin" {
  name      = "httpbin-pool"
  namespace = "demo"

  origin_servers {
    public_name {
      dns_name = "httpbin.org"
    }
  }

  port = 443

  use_tls {
    sni = "httpbin.org"

    tls_config {
      default_security {}
    }

    no_mtls {}
    volterra_trusted_ca {}
  }

  healthcheck {
    name      = xcsh_healthcheck.httpbin.name
    namespace = xcsh_healthcheck.httpbin.namespace
  }

  endpoint_selection     = "LOCAL_PREFERRED"
  loadbalancer_algorithm = "ROUND_ROBIN"
}

resource "xcsh_app_firewall" "httpbin" {
  name      = "httpbin-waf"
  namespace = "demo"

  blocking {}
  use_default_blocking_page {}
  default_detection_settings {}
  allow_all_response_codes {}
}

resource "xcsh_http_loadbalancer" "httpbin" {
  name      = "httpbin-lb"
  namespace = "demo"
  domains   = ["httpbin.example.com"]

  https_auto_cert {
    http_redirect = true
    add_hsts      = false
    default_header {}

    tls_config {
      default_security {}
    }

    no_mtls {}
  }

  advertise_on_public_default_vip {}

  default_route_pools {
    pool {
      name      = xcsh_origin_pool.httpbin.name
      namespace = xcsh_origin_pool.httpbin.namespace
    }
    weight   = 1
    priority = 1
  }

  app_firewall {
    name      = xcsh_app_firewall.httpbin.name
    namespace = xcsh_app_firewall.httpbin.namespace
  }
}
