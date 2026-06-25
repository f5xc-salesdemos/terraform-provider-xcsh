// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package schema

import (
	"testing"
)

func TestParseTerraformSchema_Empty(t *testing.T) {
	input := `{"format_version":"1.0","provider_schemas":{}}`
	ts, err := ParseTerraformSchema([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ts.Resources) != 0 {
		t.Errorf("expected 0 resources, got %d", len(ts.Resources))
	}
}

func TestParseTerraformSchema_DistinguishesBlocksFromAttributes(t *testing.T) {
	input := `{
		"format_version": "1.0",
		"provider_schemas": {
			"registry.terraform.io/f5xc-salesdemos/f5xc": {
				"resource_schemas": {
					"xcsh_http_loadbalancer": {
						"version": 0,
						"block": {
							"attributes": {
								"name": {"type": "string", "required": true},
								"add_hsts": {"type": "bool", "optional": true}
							},
							"block_types": {
								"round_robin": {
									"nesting_mode": "single",
									"block": {"attributes": {}, "block_types": {}}
								},
								"https_auto_cert": {
									"nesting_mode": "single",
									"block": {
										"attributes": {
											"http_redirect": {"type": "bool", "optional": true}
										},
										"block_types": {
											"default_header": {
												"nesting_mode": "single",
												"block": {"attributes": {}, "block_types": {}}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}`

	ts, err := ParseTerraformSchema([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res, ok := ts.Resources["xcsh_http_loadbalancer"]
	if !ok {
		t.Fatal("expected xcsh_http_loadbalancer in resources")
	}

	tests := []struct {
		path     string
		expected FieldKind
	}{
		{"name", FieldKindAttribute},
		{"add_hsts", FieldKindAttribute},
		{"round_robin", FieldKindBlock},
		{"https_auto_cert", FieldKindBlock},
		{"https_auto_cert.http_redirect", FieldKindAttribute},
		{"https_auto_cert.default_header", FieldKindBlock},
	}

	for _, tt := range tests {
		got := res.FieldKind(tt.path)
		if got != tt.expected {
			t.Errorf("FieldKind(%q) = %v, want %v", tt.path, got, tt.expected)
		}
	}
}

func TestFieldKind_UnknownPath(t *testing.T) {
	ts := &TerraformSchema{
		Resources: map[string]*ResourceSchema{
			"xcsh_test": {Fields: map[string]FieldKind{}},
		},
	}
	got := ts.Resources["xcsh_test"].FieldKind("nonexistent")
	if got != FieldKindUnknown {
		t.Errorf("expected FieldKindUnknown, got %v", got)
	}
}

func TestListBlockFields(t *testing.T) {
	rs := &ResourceSchema{
		Fields: map[string]FieldKind{
			"name":                           FieldKindAttribute,
			"round_robin":                    FieldKindBlock,
			"https_auto_cert":                FieldKindBlock,
			"https_auto_cert.default_header": FieldKindBlock,
		},
	}
	blocks := rs.ListBlockFields()
	if len(blocks) != 2 {
		t.Fatalf("expected 2 top-level blocks, got %d: %v", len(blocks), blocks)
	}
	if blocks[0] != "https_auto_cert" || blocks[1] != "round_robin" {
		t.Errorf("expected [https_auto_cert round_robin], got %v", blocks)
	}
}
