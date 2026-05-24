# API Definition Resource Example
# Manages API Definition. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic API Definition configuration
resource "f5xc_api_definition" "example" {
  name      = "example-api-definition"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # API Definition configuration
  # OpenAPI spec
  swagger_specs = ["string:///base64-openapi-spec"]

  # Non-validation mode
  non_validation_mode {}
}

# The following optional fields have server-applied defaults and can be omitted:
# - swagger_specs
# - api_inventory_exclusion_list
# - api_inventory_inclusion_list
# - non_api_endpoints
# - strict_schema_origin
