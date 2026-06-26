# Ticket Tracking System Data Source Example
# Retrieves information about an existing Ticket Tracking System

# Look up an existing Ticket Tracking System by name
terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/xcsh"
      version = ">= 0.1.0"
    }
  }
}


data "xcsh_ticket_tracking_system" "example" {
  name      = "example-ticket-tracking-system"
  namespace = "system"
}

output "ticket_tracking_system_id" {
  value = data.xcsh_ticket_tracking_system.example.id
}
