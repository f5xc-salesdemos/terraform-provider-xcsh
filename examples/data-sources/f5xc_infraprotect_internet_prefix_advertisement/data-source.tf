# Infraprotect Internet Prefix Advertisement Data Source Example
# Retrieves information about an existing Infraprotect Internet Prefix Advertisement

# Look up an existing Infraprotect Internet Prefix Advertisement by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_infraprotect_internet_prefix_advertisement" "example" {
  name      = "example-infraprotect-internet-prefix-advertisement"
  namespace = "system"
}

output "infraprotect_internet_prefix_advertisement_id" {
  value = data.f5xc_infraprotect_internet_prefix_advertisement.example.id
}
