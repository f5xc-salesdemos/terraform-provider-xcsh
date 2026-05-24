# Addon Subscription Resource Example
# Manages new Addon Subscription with Addon Subscription State in F5 Distributed Cloud.

# Basic Addon Subscription configuration
terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}


resource "f5xc_addon_subscription" "example" {
  name      = "example-addon-subscription"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Object reference. This type establishes a direct reference...
  addon_service {
    # Configure addon_service settings
  }
  # Notification Preference. NotificationPreference preferenc...
  notification_preference {
    # Configure notification_preference settings
  }
  # Addon Subscription Associated Emails. Addon Subscription ...
  emails {
    # Configure emails settings
  }
}
