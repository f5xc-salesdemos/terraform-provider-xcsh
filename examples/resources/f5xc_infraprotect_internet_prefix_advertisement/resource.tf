# Infraprotect Internet Prefix Advertisement Resource Example
# Manages DDoS transit Internet Prefix in F5 Distributed Cloud.

# Basic Infraprotect Internet Prefix Advertisement configuration
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


resource "f5xc_infraprotect_internet_prefix_advertisement" "example" {
  name      = "example-infraprotect-internet-prefix-advertisement"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # [OneOf: activation_announce, activation_withdraw] Enable ...
  activation_announce {
    # Configure activation_announce settings
  }
  # Enable this option
  activation_withdraw {
    # Configure activation_withdraw settings
  }
  # [OneOf: expiration_never, expiration_timestamp] Enable th...
  expiration_never {
    # Configure expiration_never settings
  }
}
