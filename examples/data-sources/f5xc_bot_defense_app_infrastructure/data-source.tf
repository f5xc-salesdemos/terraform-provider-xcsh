# Bot Defense App Infrastructure Data Source Example
# Retrieves information about an existing Bot Defense App Infrastructure

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Bot Defense App Infrastructure by name
data "f5xc_bot_defense_app_infrastructure" "example" {
  name      = "example-bot-defense-app-infrastructure"
  namespace = "shared"
}

output "bot_defense_app_infrastructure_id" {
  value = data.f5xc_bot_defense_app_infrastructure.example.id
}
