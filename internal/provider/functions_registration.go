// This file is MANUALLY MAINTAINED and is NOT auto-generated from OpenAPI specifications.
// It registers provider-defined functions for utility operations that are not part of
// the XCSH API specification.
//
// DO NOT DELETE OR MODIFY during code generation. This file is preserved by the
// generate-all-schemas.go tool.

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"

	"github.com/f5xc-salesdemos/terraform-provider-xcsh/internal/functions"
)

// Ensure XCShProvider satisfies the provider.ProviderWithFunctions interface.
var _ provider.ProviderWithFunctions = &XCShProvider{}

// Functions returns utility functions provided by this provider.
// These include blindfold encryption functions for XCSH Secret Management.
//
// Available functions:
//   - blindfold: Encrypts base64-encoded plaintext using XCSH Secret Management
//   - blindfold_file: Reads a file and encrypts its contents using XCSH Secret Management
//
// These functions require Terraform 1.8.0 or later and are called using the provider namespace:
//
//	provider::xcsh::blindfold(plaintext, policy_name, namespace)
//	provider::xcsh::blindfold_file(path, policy_name, namespace)
func (p *XCShProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		functions.NewBlindfoldFunction,
		functions.NewBlindfoldFileFunction,
	}
}
