# Cloud Credentials Data Source Example
# Retrieves information about an existing Cloud Credentials

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Cloud Credentials by name
data "f5xc_cloud_credentials" "example" {
  name      = "example-cloud-credentials"
  namespace = "staging"
}

output "cloud_credentials_id" {
  value = data.f5xc_cloud_credentials.example.id
}

# Example: Reference cloud credentials in site configuration
# resource "f5xc_aws_vpc_site" "example" {
#   name      = "example-aws-site"
#   namespace = "staging"
#
#   aws_cred {
#     name      = data.f5xc_cloud_credentials.example.name
#     namespace = data.f5xc_cloud_credentials.example.namespace
#   }
# }
