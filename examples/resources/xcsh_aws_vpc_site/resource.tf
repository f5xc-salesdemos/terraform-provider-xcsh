# AWS VPC Site Resource Example
# Manages a AWS VPC Site resource in F5 Distributed Cloud for deploying F5 sites within AWS VPC environments.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic AWS VPC Site configuration
resource "xcsh_aws_vpc_site" "example" {
  name      = "example-aws-vpc-site"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # AWS VPC Site configuration
  aws_region = "us-west-2"

  # AWS credentials reference
  aws_cred {
    name      = "aws-credentials"
    namespace = "staging"
  }

  # VPC configuration
  vpc {
    new_vpc {
      name_tag     = "xcsh-vpc"
      primary_ipv4 = "10.0.0.0/16"
    }
  }

  # Instance type
  instance_type = "t3.xlarge"

  # Ingress/Egress gateway
  ingress_egress_gw {
    aws_certified_hw = "aws-byol-multi-nic-voltmesh"
    az_nodes {
      aws_az_name = "us-west-2a"
      inside_subnet {
        subnet_param {
          ipv4 = "10.0.1.0/24"
        }
      }
      outside_subnet {
        subnet_param {
          ipv4 = "10.0.2.0/24"
        }
      }
    }
  }

  # No worker nodes by default
  no_worker_nodes {}
}
