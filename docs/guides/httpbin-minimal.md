---
page_title: "Guide: Minimal HTTP Load Balancer with httpbin"
subcategory: "Guides"
description: |-
  Copy-paste Terraform configuration for an HTTP Load Balancer with WAF,
  origin pool, and health check pointing to httpbin.org.
---

# Minimal HTTP Load Balancer with httpbin

A complete, validated Terraform configuration that deploys four resources:

1. **Health Check** — HTTP GET on `/get` against httpbin.org
2. **Origin Pool** — Single public origin pointing to `httpbin.org:443` with TLS
3. **App Firewall** — WAF in blocking mode with default detection settings
4. **HTTP Load Balancer** — Ties it all together with auto-cert HTTPS

## Complete Configuration

Copy this entire block into a `main.tf` file. Set `F5XC_API_URL` and `F5XC_API_TOKEN` environment variables, then run `terraform init && terraform apply`.

~> **Important:** Replace `your-namespace` with your actual F5 XC namespace.

```terraform
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

provider "f5xc" {}

# 1. Health Check — monitors httpbin.org /get endpoint
resource "f5xc_healthcheck" "httpbin" {
  name      = "httpbin-health"
  namespace = "your-namespace"

  http_health_check {
    use_origin_server_name {}
    path = "/get"
  }

  healthy_threshold   = 3
  unhealthy_threshold = 3
  interval            = 15
  timeout             = 5
}

# 2. Origin Pool — httpbin.org over TLS
resource "f5xc_origin_pool" "httpbin" {
  name      = "httpbin-pool"
  namespace = "your-namespace"

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
    name      = f5xc_healthcheck.httpbin.name
    namespace = f5xc_healthcheck.httpbin.namespace
  }

  endpoint_selection     = "LOCAL_PREFERRED"
  loadbalancer_algorithm = "ROUND_ROBIN"
}

# 3. App Firewall — WAF in blocking mode
resource "f5xc_app_firewall" "httpbin" {
  name      = "httpbin-waf"
  namespace = "your-namespace"

  blocking {}
  use_default_blocking_page {}
  default_detection_settings {}
  allow_all_response_codes {}
}

# 4. HTTP Load Balancer — ties it all together
resource "f5xc_http_loadbalancer" "httpbin" {
  name      = "httpbin-lb"
  namespace = "your-namespace"
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
      name      = f5xc_origin_pool.httpbin.name
      namespace = f5xc_origin_pool.httpbin.namespace
    }
    weight   = 1
    priority = 1
  }

  app_firewall {
    name      = f5xc_app_firewall.httpbin.name
    namespace = f5xc_app_firewall.httpbin.namespace
  }
}
```

## How the Resources Connect

```text
Health Check ──► Origin Pool ──► HTTP Load Balancer
                                        │
App Firewall ───────────────────────────┘
```

- `f5xc_healthcheck` monitors origin health
- `f5xc_origin_pool` references the health check and defines the backend
- `f5xc_app_firewall` defines WAF policy
- `f5xc_http_loadbalancer` references the origin pool and app firewall

## What This Omits

This configuration relies on server-applied defaults for fields you don't need to set:

- **Load balancing:** `round_robin` (server default)
- **Rate limiting:** `disable_rate_limit` (server default)
- **Bot defense:** `disable_bot_defense` (server default)
- **Challenge:** `no_challenge` (server default)
- **API features:** `disable_api_definition`, `disable_api_discovery`, `disable_api_testing` (server defaults)

To enable any of these features, add them explicitly. See the [full HTTP Load Balancer guide](http-loadbalancer.md) for production configurations.
