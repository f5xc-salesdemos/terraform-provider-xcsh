// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package schema

import (
	"fmt"
	"strings"

	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/namespace"
	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/openapi"
)

// GenerateExampleUsage creates a sample HCL configuration for the resource.
func GenerateExampleUsage(resourceName string, attributes []openapi.TerraformAttribute) string {
	// Get the correct namespace for this resource (fixes bug: was hardcoded to "system")
	_, ns := namespace.ForResource(resourceName)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("resource \"f5xc_%s\" \"example\" {\n", resourceName))
	sb.WriteString("  name      = \"example\"\n")
	sb.WriteString(fmt.Sprintf("  namespace = %q\n", ns))
	sb.WriteString("\n")
	sb.WriteString("  labels = {\n")
	sb.WriteString("    \"env\" = \"production\"\n")
	sb.WriteString("  }\n")

	// Add example for optional scalar attributes (skip standard ones)
	standardAttrs := map[string]bool{"name": true, "namespace": true, "labels": true, "annotations": true, "id": true}
	for _, attr := range attributes {
		if standardAttrs[attr.Name] || attr.Computed || attr.IsBlock {
			continue
		}
		if attr.Optional {
			switch attr.Type {
			case "string":
				sb.WriteString(fmt.Sprintf("\n  %s = \"example-value\"\n", attr.Name))
			case "int64":
				sb.WriteString(fmt.Sprintf("\n  %s = 0\n", attr.Name))
			case "bool":
				sb.WriteString(fmt.Sprintf("\n  %s = false\n", attr.Name))
			}
			break // Just add one example optional attribute
		}
	}

	// Add example for first nested block with attributes
	for _, attr := range attributes {
		if attr.IsBlock && len(attr.NestedAttributes) > 0 {
			sb.WriteString(fmt.Sprintf("\n  %s {\n", attr.Name))
			for _, nested := range attr.NestedAttributes {
				if nested.IsBlock {
					continue // Skip deeply nested blocks in example
				}
				switch nested.Type {
				case "string":
					sb.WriteString(fmt.Sprintf("    %s = \"example\"\n", nested.Name))
				case "int64":
					sb.WriteString(fmt.Sprintf("    %s = 0\n", nested.Name))
				case "bool":
					sb.WriteString(fmt.Sprintf("    %s = false\n", nested.Name))
				}
				if len(sb.String()) > 500 { // Keep example concise
					break
				}
			}
			sb.WriteString("  }\n")
			break // Just add one example block
		}
	}

	sb.WriteString("}")
	return sb.String()
}
