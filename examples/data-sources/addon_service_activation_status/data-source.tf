# Example: Check addon service activation status
# Use this to determine if an addon service can be activated

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_addon_service_activation_status" "bot_defense" {
  addon_service = "bot_defense"
}

# Output the activation status
output "state" {
  description = "Current subscription state (AS_NONE, AS_PENDING, AS_SUBSCRIBED, AS_ERROR)"
  value       = data.xcsh_addon_service_activation_status.bot_defense.state
}

output "can_activate" {
  description = "Whether the addon can be activated"
  value       = data.xcsh_addon_service_activation_status.bot_defense.can_activate
}

output "status_message" {
  description = "Human-readable status message"
  value       = data.xcsh_addon_service_activation_status.bot_defense.message
}

# Example: Conditional subscription based on activation status
# Only create the subscription if the addon is available for activation
resource "xcsh_addon_subscription" "bot_defense" {
  count = data.xcsh_addon_service_activation_status.bot_defense.can_activate && data.xcsh_addon_service_activation_status.bot_defense.state == "AS_NONE" ? 1 : 0

  name      = "my-bot-defense-subscription"
  namespace = "system"

  addon_service {
    name      = "bot_defense"
    namespace = "shared"
  }
}
