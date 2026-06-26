# Global Log Receiver Resource Example
# Manages new Global Log Receiver object. in F5 Distributed Cloud.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Global Log Receiver configuration
resource "xcsh_global_log_receiver" "example" {
  name      = "example-global-log-receiver"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # [OneOf: audit_logs, dns_logs, request_logs, security_even...
  audit_logs {
    # Configure audit_logs settings
  }
  # [OneOf: aws_cloud_watch_receiver, azure_event_hubs_receiv...
  aws_cloud_watch_receiver {
    # Configure aws_cloud_watch_receiver settings
  }
  # Type establishes a direct reference from one object(the r...
  aws_cred {
    # Configure aws_cred settings
  }
}

# The following optional fields have server-applied defaults and can be omitted:
# - ns_current
