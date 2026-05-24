# Flow Anomaly Data Source Example
# Retrieves information about an existing Flow Anomaly

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Flow Anomaly by name
data "f5xc_flow_anomaly" "example" {
  name      = "example-flow-anomaly"
  namespace = "staging"
}

output "flow_anomaly_id" {
  value = data.f5xc_flow_anomaly.example.id
}
