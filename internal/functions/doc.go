// Package functions provides Terraform provider-defined functions for F5XC.
//
// This package implements utility functions that extend the provider's capabilities
// beyond the auto-generated resources from OpenAPI specifications. These functions
// are available in Terraform 1.8+ and can be called using the provider namespace:
//
//	provider::xcsh::blindfold(plaintext, policy_name, namespace)
//	provider::xcsh::blindfold_file(path, policy_name, namespace)
//
// This package is MANUALLY MAINTAINED and is NOT auto-generated from OpenAPI
// specifications. Changes to this package should be committed directly.
//
// # Available Functions
//
// The following functions are provided:
//
//   - blindfold: Encrypts base64-encoded plaintext using XCSH Secret Management
//   - blindfold_file: Reads a file and encrypts its contents using XCSH Secret Management
//
// # Example Usage
//
//	resource "xcsh_http_loadbalancer" "example" {
//	  name = "my-lb"
//
//	  tls_config {
//	    private_key {
//	      blindfold_secret_info {
//	        location = provider::xcsh::blindfold_file(
//	          "${path.module}/certs/private.key",
//	          "my-secret-policy",
//	          "shared"
//	        )
//	      }
//	    }
//	  }
//	}
//
// # Requirements
//
//   - Terraform 1.8.0 or later
//   - Valid XCSH provider configuration with API credentials
//   - Existing SecretPolicy in the specified namespace
package functions
