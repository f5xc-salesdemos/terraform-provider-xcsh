// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package schema

import (
	"testing"
)

func TestCrossCheck_NoMismatches(t *testing.T) {
	tfSchema := &TerraformSchema{
		Resources: map[string]*ResourceSchema{
			"xcsh_http_loadbalancer": {
				Fields: map[string]FieldKind{
					"round_robin": FieldKindBlock,
					"name":        FieldKindAttribute,
				},
			},
		},
	}
	openAPIFields := map[string]map[string]FieldKind{
		"xcsh_http_loadbalancer": {
			"round_robin": FieldKindBlock,
			"name":        FieldKindAttribute,
		},
	}

	mismatches := CrossCheck(tfSchema, openAPIFields)
	if len(mismatches) != 0 {
		t.Errorf("expected 0 mismatches, got %d: %v", len(mismatches), mismatches)
	}
}

func TestCrossCheck_DetectsMismatch(t *testing.T) {
	tfSchema := &TerraformSchema{
		Resources: map[string]*ResourceSchema{
			"xcsh_http_loadbalancer": {
				Fields: map[string]FieldKind{
					"round_robin": FieldKindBlock,
				},
			},
		},
	}
	openAPIFields := map[string]map[string]FieldKind{
		"xcsh_http_loadbalancer": {
			"round_robin": FieldKindAttribute,
		},
	}

	mismatches := CrossCheck(tfSchema, openAPIFields)
	if len(mismatches) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(mismatches))
	}
	if mismatches[0].Resource != "xcsh_http_loadbalancer" {
		t.Errorf("expected resource xcsh_http_loadbalancer, got %s", mismatches[0].Resource)
	}
	if mismatches[0].Field != "round_robin" {
		t.Errorf("expected field round_robin, got %s", mismatches[0].Field)
	}
}

func TestCrossCheck_SkipsMissingResources(t *testing.T) {
	tfSchema := &TerraformSchema{
		Resources: map[string]*ResourceSchema{
			"xcsh_http_loadbalancer": {
				Fields: map[string]FieldKind{"name": FieldKindAttribute},
			},
		},
	}
	openAPIFields := map[string]map[string]FieldKind{
		"xcsh_nonexistent": {
			"name": FieldKindAttribute,
		},
	}

	mismatches := CrossCheck(tfSchema, openAPIFields)
	if len(mismatches) != 0 {
		t.Errorf("expected 0 mismatches for missing resource, got %d", len(mismatches))
	}
}
