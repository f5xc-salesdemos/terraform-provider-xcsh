# Addon Subscription Data Source Example
# Retrieves information about an existing Addon Subscription

# Look up an existing Addon Subscription by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_addon_subscription" "example" {
  name      = "example-addon-subscription"
  namespace = "system"
}

output "addon_subscription_id" {
  value = data.xcsh_addon_subscription.example.id
}
