# Geo Location Set Resource Example
# Manages Geolocation Set in F5 Distributed Cloud.

# Basic Geo Location Set configuration
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


resource "f5xc_geo_location_set" "example" {
  name      = "example-geo-location-set"
  namespace = "shared"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Geo Location Set configuration
  country_codes = ["US", "CA", "GB"]
}
