# Example: Query addon service details
# This data source retrieves information about an F5 Distributed Cloud addon service

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_addon_service" "bot_defense" {
  name = "bot_defense"
}

# Output the addon service details
output "display_name" {
  description = "Human-readable name of the addon service"
  value       = data.xcsh_addon_service.bot_defense.display_name
}

output "tier" {
  description = "Required subscription tier (NO_TIER, BASIC, STANDARD, ADVANCED, PREMIUM)"
  value       = data.xcsh_addon_service.bot_defense.tier
}

output "activation_type" {
  description = "How the service is activated (self, partial, managed)"
  value       = data.xcsh_addon_service.bot_defense.activation_type
}

# Example: Check if self-activation is available
output "is_self_activatable" {
  description = "Whether the addon can be activated without manual intervention"
  value       = data.xcsh_addon_service.bot_defense.activation_type == "self"
}
