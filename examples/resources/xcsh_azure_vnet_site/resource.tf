# Azure VNET Site Resource Example
# Manages a Azure VNET Site resource in F5 Distributed Cloud for deploying F5 sites within Azure Virtual Network environments.

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5-sales-demo/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Basic Azure VNET Site configuration
resource "xcsh_azure_vnet_site" "example" {
  name      = "example-azure-vnet-site"
  namespace = "staging"

  labels = {
    environment = "production"
    managed_by  = "terraform"
  }

  annotations = {
    "owner" = "platform-team"
  }

  # Azure VNET Site configuration
  azure_region = "westus2"

  # Azure credentials reference
  azure_cred {
    name      = "azure-credentials"
    namespace = "staging"
  }

  # Resource group
  resource_group = "xcsh-rg"

  # VNET configuration
  vnet {
    new_vnet {
      name         = "xcsh-vnet"
      primary_ipv4 = "10.0.0.0/16"
    }
  }

  # Machine type
  machine_type = "Standard_D3_v2"

  # Ingress/Egress gateway
  ingress_egress_gw {
    azure_certified_hw = "azure-byol-multi-nic-voltmesh"
    az_nodes {
      azure_az = "1"
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

# The following optional fields have server-applied defaults and can be omitted:
# - disk_size
# - block_all_services
# - logs_streaming_disabled
# - no_worker_nodes
# - tags
