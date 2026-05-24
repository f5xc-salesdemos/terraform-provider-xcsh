# CDN Purge Command Data Source Example
# Retrieves information about an existing CDN Purge Command

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing CDN Purge Command by name
data "f5xc_cdn_purge_command" "example" {
  name      = "example-cdn-purge-command"
  namespace = "staging"
}

output "cdn_purge_command_id" {
  value = data.f5xc_cdn_purge_command.example.id
}
