# HTTP Loadbalancer Resource Example
# Manages a HTTP Load Balancer resource in F5 Distributed Cloud for load balancing HTTP/HTTPS traffic with advanced routing and security.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Basic HTTP Loadbalancer configuration
resource "xcsh_http_loadbalancer" "example" {
  name      = "example-http-loadbalancer"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  // One of the arguments from this list "advertise_custom advertise_on_public advertise_on_public_default_vip do_not_advertise" must be set

  advertise_on_public_default_vip {}

  // One of the arguments from this list "api_specification disable_api_definition" must be set

  disable_api_definition {}

  // One of the arguments from this list "disable_api_discovery enable_api_discovery" must be set

  disable_api_discovery {}

  // One of the arguments from this list "api_testing disable_api_testing" must be set

  disable_api_testing {}

  // One of the arguments from this list "captcha_challenge enable_challenge js_challenge no_challenge policy_based_challenge" must be set

  no_challenge {}

  domains = ["app.example.com", "www.example.com"]

  // One of the arguments from this list "cookie_stickiness least_active random ring_hash round_robin source_ip_stickiness" must be set

  round_robin {}

  // One of the arguments from this list "http https https_auto_cert" must be set

  https_auto_cert {
    http_redirect = true
    add_hsts      = true

    // One of the arguments from this list "default_header no_headers server_name" must be set

    default_header {}

    tls_config {
      // One of the arguments from this list "custom_security default_security low_security medium_security" must be set

      default_security {}
    }

    // One of the arguments from this list "no_mtls use_mtls" must be set

    no_mtls {}
  }

  // One of the arguments from this list "disable_malicious_user_detection enable_malicious_user_detection" must be set

  enable_malicious_user_detection {}

  // One of the arguments from this list "disable_malware_protection malware_protection_settings" must be set

  disable_malware_protection {}

  // One of the arguments from this list "api_rate_limit disable_rate_limit rate_limit" must be set

  disable_rate_limit {}

  // One of the arguments from this list "default_sensitive_data_policy sensitive_data_policy" must be set

  default_sensitive_data_policy {}

  // One of the arguments from this list "active_service_policies no_service_policies service_policies_from_namespace" must be set

  service_policies_from_namespace {}

  // One of the arguments from this list "disable_threat_mesh enable_threat_mesh" must be set

  enable_threat_mesh {}

  // One of the arguments from this list "disable_trust_client_ip_headers enable_trust_client_ip_headers" must be set

  disable_trust_client_ip_headers {}

  // One of the arguments from this list "user_id_client_ip user_identification" must be set

  user_id_client_ip {}

  // One of the arguments from this list "app_firewall disable_waf" must be set

  app_firewall {
    name      = "example-app-firewall"
    namespace = "staging"
  }

  // One of the arguments from this list "bot_defense bot_defense_advanced disable_bot_defense" must be set

  disable_bot_defense {}

  // Default route pools configuration
  default_route_pools {
    pool {
      name      = "example-origin-pool"
      namespace = "staging"
    }
    weight   = 1
    priority = 1
  }
}

# The following optional fields have server-applied defaults and can be omitted:
# - add_location
# - endpoint_selection
# - loadbalancer_algorithm
# - healthcheck
# - no_tls
# - same_as_endpoint_port
# - default_sensitive_data_policy
# - disable_api_definition
# - disable_api_discovery
# - disable_api_testing
# - disable_malware_protection
# - disable_rate_limit
# - disable_threat_mesh
# - disable_trust_client_ip_headers
# - l7_ddos_protection
# - round_robin
# - service_policies_from_namespace
# - user_id_client_ip
