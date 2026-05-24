# Subnet Data Source Example
# Retrieves information about an existing Subnet

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Subnet by name
data "f5xc_subnet" "example" {
  name      = "example-subnet"
  namespace = "system"
}

output "subnet_id" {
  value = data.f5xc_subnet.example.id
}
