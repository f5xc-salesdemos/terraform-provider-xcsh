# DNS LB Health Check Data Source Example
# Retrieves information about an existing DNS LB Health Check

# Look up an existing DNS LB Health Check by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_dns_lb_health_check" "example" {
  name      = "example-dns-lb-health-check"
  namespace = "system"
}

output "dns_lb_health_check_id" {
  value = data.f5xc_dns_lb_health_check.example.id
}
