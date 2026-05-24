# DNS Zone Resource Example
# Manages DNS Zone in a given namespace. If one already exist it will give a error. in F5 Distributed Cloud.

# Basic DNS Zone configuration
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


resource "f5xc_dns_zone" "example" {
  name      = "example-dns-zone"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # DNS Zone configuration
  # Primary DNS zone
  primary {
    soa_record_parameters {
      refresh = 86400
      retry   = 7200
      expire  = 3600000
      ttl     = 86400
      neg_ttl = 1800
    }
    default_rr_set_group {}
    default_soa_parameters {}
    dnssec_mode {
      disable {}
    }
  }
}
