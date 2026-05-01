// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package openapi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSchema_IsRef(t *testing.T) {
	tests := []struct {
		name     string
		schema   Schema
		expected bool
	}{
		{"with ref", Schema{Ref: "#/components/schemas/Foo"}, true},
		{"without ref", Schema{Type: "string"}, false},
		{"empty ref", Schema{Ref: ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.schema.IsRef(); got != tt.expected {
				t.Errorf("Schema.IsRef() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSchema_IsArray(t *testing.T) {
	tests := []struct {
		name     string
		schema   Schema
		expected bool
	}{
		{"array type", Schema{Type: "array"}, true},
		{"object type", Schema{Type: "object"}, false},
		{"string type", Schema{Type: "string"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.schema.IsArray(); got != tt.expected {
				t.Errorf("Schema.IsArray() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSchema_IsObject(t *testing.T) {
	tests := []struct {
		name     string
		schema   Schema
		expected bool
	}{
		{"object type", Schema{Type: "object"}, true},
		{"array type", Schema{Type: "array"}, false},
		{"string type", Schema{Type: "string"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.schema.IsObject(); got != tt.expected {
				t.Errorf("Schema.IsObject() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSchema_IsRequired(t *testing.T) {
	schema := Schema{
		Required: []string{"name", "namespace"},
	}

	tests := []struct {
		property string
		expected bool
	}{
		{"name", true},
		{"namespace", true},
		{"description", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.property, func(t *testing.T) {
			if got := schema.IsRequired(tt.property); got != tt.expected {
				t.Errorf("Schema.IsRequired(%q) = %v, want %v", tt.property, got, tt.expected)
			}
		})
	}
}

func TestSchema_HasProperties(t *testing.T) {
	tests := []struct {
		name     string
		schema   Schema
		expected bool
	}{
		{
			"with properties",
			Schema{Properties: map[string]Schema{"foo": {Type: "string"}}},
			true,
		},
		{
			"empty properties",
			Schema{Properties: map[string]Schema{}},
			false,
		},
		{
			"nil properties",
			Schema{},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.schema.HasProperties(); got != tt.expected {
				t.Errorf("Schema.HasProperties() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetRefName(t *testing.T) {
	tests := []struct {
		ref      string
		expected string
	}{
		{"#/components/schemas/Namespace", "Namespace"},
		{"#/components/schemas/foo.bar.Baz", "foo.bar.Baz"}, // GetRefName extracts after last /
		{"Namespace", "Namespace"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.ref, func(t *testing.T) {
			if got := GetRefName(tt.ref); got != tt.expected {
				t.Errorf("GetRefName(%q) = %q, want %q", tt.ref, got, tt.expected)
			}
		})
	}
}

func TestSpec_ResolveRef(t *testing.T) {
	spec := &Spec{
		Components: Components{
			Schemas: map[string]Schema{
				"Namespace": {Type: "object", Description: "A namespace"},
			},
		},
	}

	tests := []struct {
		ref     string
		wantErr bool
	}{
		{"#/components/schemas/Namespace", false},
		{"#/components/schemas/NotFound", true},
		{"external/ref", true},
	}

	for _, tt := range tests {
		t.Run(tt.ref, func(t *testing.T) {
			schema, err := spec.ResolveRef(tt.ref)
			if (err != nil) != tt.wantErr {
				t.Errorf("Spec.ResolveRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && schema == nil {
				t.Error("Spec.ResolveRef() returned nil schema")
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-spec.json")

	testJSON := `{
		"openapi": "3.0.0",
		"info": {"title": "Test API", "version": "1.0.0"},
		"paths": {},
		"components": {
			"schemas": {
				"TestSchema": {
					"type": "object",
					"description": "A test schema"
				}
			}
		}
	}`

	if err := os.WriteFile(testFile, []byte(testJSON), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	spec, err := ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	if spec.OpenAPI != "3.0.0" {
		t.Errorf("ParseFile().OpenAPI = %q, want %q", spec.OpenAPI, "3.0.0")
	}

	if _, ok := spec.Components.Schemas["TestSchema"]; !ok {
		t.Error("ParseFile() missing TestSchema in components")
	}
}

func TestParseIndex_Timestamp(t *testing.T) {
	tmpDir := t.TempDir()
	indexFile := filepath.Join(tmpDir, "index.json")
	indexJSON := `{
		"version": "2.1.62",
		"timestamp": "2026-04-23T09:16:13Z",
		"specifications": [],
		"x-f5xc-critical-resources": ["http_loadbalancer"],
		"x-f5xc-error-resolution": {"401": "Token expired"},
		"x-f5xc-guided-workflows": {"deploy_lb": {"steps": []}},
		"x-f5xc-acronyms": {"LB": "Load Balancer", "WAF": "Web Application Firewall"}
	}`
	if err := os.WriteFile(indexFile, []byte(indexJSON), 0644); err != nil {
		t.Fatalf("Failed to create index file: %v", err)
	}

	index, err := ParseIndex(indexFile)
	if err != nil {
		t.Fatalf("ParseIndex() error = %v", err)
	}

	if index.Version != "2.1.62" {
		t.Errorf("Version = %q, want %q", index.Version, "2.1.62")
	}
	if index.Timestamp != "2026-04-23T09:16:13Z" {
		t.Errorf("Timestamp = %q, want %q", index.Timestamp, "2026-04-23T09:16:13Z")
	}
	if len(index.CriticalResources) != 1 {
		t.Errorf("CriticalResources len = %d, want 1", len(index.CriticalResources))
	}
	if index.CriticalResources[0] != "http_loadbalancer" {
		t.Errorf("CriticalResources[0] = %q, want %q", index.CriticalResources[0], "http_loadbalancer")
	}
	if len(index.Acronyms) != 2 {
		t.Errorf("Acronyms len = %d, want 2", len(index.Acronyms))
	}
}

func TestSchema_NewExtensions(t *testing.T) {
	specJSON := `{
		"openapi": "3.0.0",
		"info": {"title": "Test", "version": "1.0.0"},
		"paths": {},
		"components": {
			"schemas": {
				"TestSchema": {
					"type": "object",
					"x-ves-deprecated": "Use v2 endpoint instead",
					"x-f5xc-conflicts-with": ["field_a", "field_b"],
					"x-field-mutability": "immutable",
					"x-original-maxLength": 256,
					"x-validation-rules": {"max_length": "256"},
					"x-required": true,
					"x-ves-required": "true",
					"x-ves-displayorder": "5",
					"x-ves-proto-enum": "ves.io.schema.SomeEnum",
					"x-f5xc-constraints": {"min_items": 1},
					"x-f5xc-recommended-oneof-variant": "option_a",
					"x-reconciled-at": "2026-04-23T09:16:13Z",
					"x-reconciled-from-discovery": true
				}
			}
		}
	}`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-spec.json")
	if err := os.WriteFile(testFile, []byte(specJSON), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	spec, err := ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	schema := spec.Components.Schemas["TestSchema"]

	if schema.XVesDeprecated != "Use v2 endpoint instead" {
		t.Errorf("XVesDeprecated = %q, want %q", schema.XVesDeprecated, "Use v2 endpoint instead")
	}
	if len(schema.XF5XCConflictsWith) != 2 || schema.XF5XCConflictsWith[0] != "field_a" {
		t.Errorf("XF5XCConflictsWith = %v, want [field_a, field_b]", schema.XF5XCConflictsWith)
	}
	if schema.XFieldMutability != "immutable" {
		t.Errorf("XFieldMutability = %q, want %q", schema.XFieldMutability, "immutable")
	}
	if schema.XOriginalMaxLength != 256 {
		t.Errorf("XOriginalMaxLength = %d, want 256", schema.XOriginalMaxLength)
	}
	if schema.XRequired != true {
		t.Errorf("XRequired = %v, want true", schema.XRequired)
	}
	if schema.XVesRequired != "true" {
		t.Errorf("XVesRequired = %q, want %q", schema.XVesRequired, "true")
	}
	if schema.XVesDisplayOrder != "5" {
		t.Errorf("XVesDisplayOrder = %q, want %q", schema.XVesDisplayOrder, "5")
	}
	if schema.XVesProtoEnum != "ves.io.schema.SomeEnum" {
		t.Errorf("XVesProtoEnum = %q, want %q", schema.XVesProtoEnum, "ves.io.schema.SomeEnum")
	}
	if schema.XF5XCRecommendedOneofVariant != "option_a" {
		t.Errorf("XF5XCRecommendedOneofVariant = %v, want %q", schema.XF5XCRecommendedOneofVariant, "option_a")
	}
	if schema.XReconciledFromDiscovery != true {
		t.Errorf("XReconciledFromDiscovery = %v, want true", schema.XReconciledFromDiscovery)
	}
}

func TestGetAPIDocURL(t *testing.T) {
	pathMap := map[string]string{
		"http_loadbalancer": "views-http-loadbalancer",
		"namespace":         "namespace",
	}

	tests := []struct {
		resource string
		expected string
	}{
		{"http_loadbalancer", "https://docs.cloud.f5.com/docs-v2/api/views-http-loadbalancer"},
		{"namespace", "https://docs.cloud.f5.com/docs-v2/api/namespace"},
		{"unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.resource, func(t *testing.T) {
			result := GetAPIDocURL(tt.resource, pathMap)
			if result != tt.expected {
				t.Errorf("GetAPIDocURL(%q) = %q, want %q", tt.resource, result, tt.expected)
			}
		})
	}
}
