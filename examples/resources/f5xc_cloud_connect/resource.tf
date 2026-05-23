# Cloud Connect Resource Example
# Manages a Cloud Connect resource in F5 Distributed Cloud for establishing connectivity to cloud provider networks.

# Basic Cloud Connect configuration
resource "f5xc_cloud_connect" "example" {
  name      = "example-cloud-connect"
  namespace = "system"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Resource-specific configuration
  # Configuration parameter for aws provider.
  aws_provider {
    # Configure aws_provider settings
  }
  # AWS TGW Site Type. Cloud Connect AWS TGW Site Type.
  aws_tgw_site {
    # Configure aws_tgw_site settings
  }
  # Type establishes a direct reference from one object(the r...
  cred {
    # Configure cred settings
  }
}
