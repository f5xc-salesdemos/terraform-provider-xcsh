# Token Data Source Example
# Retrieves information about an existing Token

# Look up an existing Token by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


data "f5xc_token" "example" {
  name      = "example-token"
  namespace = "system"
}

output "token_id" {
  value = data.f5xc_token.example.id
}
