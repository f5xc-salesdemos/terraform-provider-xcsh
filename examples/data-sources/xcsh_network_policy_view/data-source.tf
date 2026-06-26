# Network Policy View Data Source Example
# Retrieves information about an existing Network Policy View

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Network Policy View by name
data "xcsh_network_policy_view" "example" {
  name      = "example-network-policy-view"
  namespace = "staging"
}

output "network_policy_view_id" {
  value = data.xcsh_network_policy_view.example.id
}
