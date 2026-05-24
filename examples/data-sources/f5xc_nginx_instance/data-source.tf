# Nginx Instance Data Source Example
# Retrieves information about an existing Nginx Instance

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Nginx Instance by name
data "f5xc_nginx_instance" "example" {
  name      = "example-nginx-instance"
  namespace = "staging"
}

output "nginx_instance_id" {
  value = data.f5xc_nginx_instance.example.id
}
