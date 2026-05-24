# TCP Loadbalancer Resource Example
# Manages a TCP Load Balancer resource in F5 Distributed Cloud for load balancing TCP traffic across origin pools.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic TCP Loadbalancer configuration
resource "f5xc_tcp_loadbalancer" "example" {
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

  # Advertise on public internet
  advertise_on_internet {
    default_vip {}
  }

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
