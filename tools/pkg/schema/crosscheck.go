// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package schema

import (
	"fmt"
	"sort"
)

// Mismatch records a field where Terraform schema and OpenAPI disagree on type.
type Mismatch struct {
	Resource    string
	Field       string
	TFKind      FieldKind
	OpenAPIKind FieldKind
}

func (m Mismatch) String() string {
	return fmt.Sprintf("%s.%s: terraform=%s openapi=%s", m.Resource, m.Field, m.TFKind, m.OpenAPIKind)
}

// CrossCheck compares OpenAPI-derived field types against the Terraform schema.
// openAPIFields maps resourceName → fieldName → FieldKind.
// Returns mismatches where both sources have an opinion and disagree.
func CrossCheck(tfSchema *TerraformSchema, openAPIFields map[string]map[string]FieldKind) []Mismatch {
	var mismatches []Mismatch

	resources := make([]string, 0, len(openAPIFields))
	for r := range openAPIFields {
		resources = append(resources, r)
	}
	sort.Strings(resources)

	for _, resource := range resources {
		fields := openAPIFields[resource]
		tfResource, ok := tfSchema.Resources[resource]
		if !ok {
			continue
		}

		fieldNames := make([]string, 0, len(fields))
		for f := range fields {
			fieldNames = append(fieldNames, f)
		}
		sort.Strings(fieldNames)

		for _, field := range fieldNames {
			openAPIKind := fields[field]
			tfKind := tfResource.FieldKind(field)
			if tfKind == FieldKindUnknown {
				continue
			}
			if tfKind != openAPIKind {
				mismatches = append(mismatches, Mismatch{
					Resource:    resource,
					Field:       field,
					TFKind:      tfKind,
					OpenAPIKind: openAPIKind,
				})
			}
		}
	}

	return mismatches
}
