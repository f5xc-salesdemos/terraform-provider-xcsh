// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package schema

import (
	"testing"

	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/openapi"
)

func TestExtractOneOfGroups_Empty(t *testing.T) {
	spec := &openapi.Spec{}
	result := ExtractOneOfGroups(spec, "NonExistent")
	if len(result) != 0 {
		t.Errorf("Expected empty map, got %d groups", len(result))
	}
}

func TestExtractOneOfGroups_FromCache(t *testing.T) {
	spec := &openapi.Spec{}

	// Put raw schema data into cache
	RawSpecCache["TestType"] = map[string]interface{}{
		"x-ves-oneof-field-choice": []interface{}{"field_a", "field_b"},
		"type":                     "object",
	}
	defer delete(RawSpecCache, "TestType")

	result := ExtractOneOfGroups(spec, "TestType")
	if len(result) != 1 {
		t.Fatalf("Expected 1 group, got %d", len(result))
	}
	fields, ok := result["choice"]
	if !ok {
		t.Fatal("Expected group 'choice'")
	}
	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields in group, got %d", len(fields))
	}
	if fields[0] != "field_a" || fields[1] != "field_b" {
		t.Errorf("Unexpected fields: %v", fields)
	}
}

func TestExtractOneOfGroups_StringFormat(t *testing.T) {
	spec := &openapi.Spec{}

	RawSpecCache["StringType"] = map[string]interface{}{
		"x-ves-oneof-field-mode": `["fast","slow"]`,
	}
	defer delete(RawSpecCache, "StringType")

	result := ExtractOneOfGroups(spec, "StringType")
	if len(result) != 1 {
		t.Fatalf("Expected 1 group, got %d", len(result))
	}
	fields := result["mode"]
	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(fields))
	}
}

func TestExtractResourceSchema_NoCreateSpecType(t *testing.T) {
	spec := &openapi.Spec{
		Components: openapi.Components{
			Schemas: map[string]openapi.Schema{
				"SomeOtherType": {Type: "object"},
			},
		},
	}
	extractAPIPath := func(spec *openapi.Spec, resourceName string) (string, string, bool) {
		return "/api/test", "/api/test/{name}", true
	}
	_, err := ExtractResourceSchema(spec, "test_resource", extractAPIPath)
	if err == nil {
		t.Error("Expected error when no CreateSpecType found")
	}
}

func TestExtractResourceSchema_Basic(t *testing.T) {
	spec := &openapi.Spec{
		Components: openapi.Components{
			Schemas: map[string]openapi.Schema{
				"ves.io.schema.my_resource.CreateSpecType": {
					Type:        "object",
					Description: "A test resource",
					Properties: map[string]openapi.Schema{
						"port": {Type: "integer", Description: "Port number"},
					},
				},
			},
		},
	}
	extractAPIPath := func(spec *openapi.Spec, resourceName string) (string, string, bool) {
		return "/api/config/namespaces/%s/my_resources", "/api/config/namespaces/%s/my_resources/%s", true
	}
	result, err := ExtractResourceSchema(spec, "my_resource", extractAPIPath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result.Name != "my_resource" {
		t.Errorf("Name = %q, want %q", result.Name, "my_resource")
	}
	if result.HasNamespaceInPath != true {
		t.Error("HasNamespaceInPath should be true")
	}
	// Should have standard attrs (name, namespace, annotations, description, disable, labels, id) + port
	foundPort := false
	for _, attr := range result.Attributes {
		if attr.Name == "port" {
			foundPort = true
		}
	}
	if !foundPort {
		t.Error("Expected 'port' attribute in result")
	}
}
