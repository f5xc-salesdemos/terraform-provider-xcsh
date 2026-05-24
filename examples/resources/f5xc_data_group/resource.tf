# Data Group Resource Example
# Manages data group in a given namespace. If one already exists it will give an error. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Data Group configuration
resource "f5xc_data_group" "example" {
  name      = "example-data-group"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # [OneOf: address_records, integer_records, string_records]...
  address_records {
    # Configure address_records settings
  }
  # Address records. Configuration parameter for records
  records {
    # Configure records settings
  }
  # Configuration parameter for integer records.
  integer_records {
    # Configure integer_records settings
  }
}
