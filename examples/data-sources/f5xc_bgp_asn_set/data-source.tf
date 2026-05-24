# BGP Asn Set Data Source Example
# Retrieves information about an existing BGP Asn Set

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing BGP Asn Set by name
data "f5xc_bgp_asn_set" "example" {
  name      = "example-bgp-asn-set"
  namespace = "system"
}

output "bgp_asn_set_id" {
  value = data.f5xc_bgp_asn_set.example.id
}
