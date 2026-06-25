# TCP Loadbalancer Resource Example
# Manages a TCP Load Balancer resource in F5 Distributed Cloud for load balancing TCP traffic across origin pools.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Basic TCP Loadbalancer configuration
resource "xcsh_tcp_loadbalancer" "example" {
  name      = "example-tcp-loadbalancer"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # TCP Load Balancer specific configuration
  listen_port = 8443

  // One of the arguments from this list "advertise_custom advertise_on_public advertise_on_public_default_vip do_not_advertise" must be set

  advertise_on_public_default_vip {}

  # Origin pools
  origin_pools_weights {
    pool {
      name      = "example-tcp-pool"
      namespace = "staging"
    }
    weight = 1
  }

  # DNS for TCP load balancer
  dns_volterra_managed = true

  # No retract cluster by default
  retract_cluster {}
}

# The following optional fields have server-applied defaults and can be omitted:
# - dns_volterra_managed
# - idle_timeout
# - hash_policy_choice_round_robin
# - no_sni
# - retract_cluster
# - service_policies_from_namespace
# - tcp
