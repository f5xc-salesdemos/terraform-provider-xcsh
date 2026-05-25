# WithListenPort — Verified Configuration Example
# This configuration is extracted from acceptance tests
# and verified against the live F5 XC API.

resource "f5xc_origin_pool" "test" {
  name       = "example-pool"
  namespace  = "system"
  port       = 443

  origin_servers {
    labels {}
    public_name {
      dns_name = "example.com"
    }
  }

  no_tls {}
  same_as_endpoint_port {}
}

resource "f5xc_tcp_loadbalancer" "test" {
  name       = "example"
  namespace  = "system"

  labels = {
    environment = "test"
    managed_by  = "terraform-acceptance-test"
  }

  domains = ["example.example.com"]
  listen_port = 443

  tcp {}
  sni {}

  origin_pools_weights {
    pool {
      name      = f5xc_origin_pool.test.name
      namespace = "system"
    }
    weight = 1
  }

  advertise_on_public_default_vip {}
}