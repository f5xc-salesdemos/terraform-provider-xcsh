// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package schema

import (
	"testing"

	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/openapi"
)

func TestMapSchemaType(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"string", "string"},
		{"integer", "int64"},
		{"number", "int64"},
		{"boolean", "bool"},
		{"unknown", "string"},
		{"", "string"},
	}
	for _, tt := range tests {
		got := MapSchemaType(tt.input)
		if got != tt.want {
			t.Errorf("MapSchemaType(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestMapSchemaTypeToGo(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"string", "string"},
		{"integer", "int64"},
		{"number", "int64"},
		{"boolean", "bool"},
		{"unknown", "interface{}"},
		{"", "interface{}"},
	}
	for _, tt := range tests {
		got := MapSchemaTypeToGo(tt.input)
		if got != tt.want {
			t.Errorf("MapSchemaTypeToGo(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestIsMetadataField(t *testing.T) {
	metadataFields := []string{"name", "namespace", "labels", "annotations", "description", "id", "timeouts"}
	for _, f := range metadataFields {
		if !IsMetadataField(f) {
			t.Errorf("IsMetadataField(%q) = false, want true", f)
		}
	}
	specFields := []string{"spec_field", "custom_attr", "port", "protocol"}
	for _, f := range specFields {
		if IsMetadataField(f) {
			t.Errorf("IsMetadataField(%q) = true, want false", f)
		}
	}
}

func TestFilterSpecFields(t *testing.T) {
	attrs := []openapi.TerraformAttribute{
		{TfsdkTag: "name", IsSpecField: false},
		{TfsdkTag: "namespace", IsSpecField: false},
		{TfsdkTag: "port", IsSpecField: true},
		{TfsdkTag: "protocol", IsSpecField: true},
		{TfsdkTag: "description", IsSpecField: true}, // metadata name even though IsSpecField
		{TfsdkTag: "custom", IsSpecField: true},
	}
	result := FilterSpecFields(attrs)
	if len(result) != 3 {
		t.Fatalf("FilterSpecFields returned %d attrs, want 3", len(result))
	}
	names := map[string]bool{}
	for _, r := range result {
		names[r.TfsdkTag] = true
	}
	if !names["port"] || !names["protocol"] || !names["custom"] {
		t.Errorf("FilterSpecFields missing expected fields, got %v", names)
	}
}

func TestFilterOptional(t *testing.T) {
	attrs := []openapi.TerraformAttribute{
		{Name: "required_field", Required: true, Optional: false},
		{Name: "optional_field", Required: false, Optional: true},
		{Name: "computed_field", Required: false, Optional: false, Computed: true},
		{Name: "optional_computed", Required: false, Optional: true, Computed: true},
	}
	result := FilterOptional(attrs)
	if len(result) != 1 {
		t.Fatalf("FilterOptional returned %d attrs, want 1", len(result))
	}
	if result[0].Name != "optional_field" {
		t.Errorf("FilterOptional[0].Name = %q, want %q", result[0].Name, "optional_field")
	}
}

func TestParseMinConfigRequiredFields(t *testing.T) {
	// Valid input
	raw := map[string]interface{}{
		"required_fields": []interface{}{"field_a", "spec.field_b", "field_c"},
	}
	result := ParseMinConfigRequiredFields(raw)
	if len(result) != 3 {
		t.Fatalf("ParseMinConfigRequiredFields returned %d fields, want 3", len(result))
	}
	if result[1] != "field_b" {
		t.Errorf("ParseMinConfigRequiredFields[1] = %q, want %q (spec. stripped)", result[1], "field_b")
	}

	// Nil input
	if ParseMinConfigRequiredFields(nil) != nil {
		t.Error("ParseMinConfigRequiredFields(nil) should return nil")
	}

	// Wrong type
	if ParseMinConfigRequiredFields("not a map") != nil {
		t.Error("ParseMinConfigRequiredFields(string) should return nil")
	}
}

func TestPromoteMinConfigRequired(t *testing.T) {
	attrs := []openapi.TerraformAttribute{
		{TfsdkTag: "field_a", Optional: true, Required: false},
		{TfsdkTag: "field_b", Optional: true, Required: false},
		{TfsdkTag: "field_c", Optional: false, Required: true}, // already required
	}
	minFields := map[string]bool{"field_a": true, "field_c": true}
	PromoteMinConfigRequired(attrs, minFields)

	if !attrs[0].Required || attrs[0].Optional {
		t.Error("field_a should be promoted to Required")
	}
	if attrs[1].Required || !attrs[1].Optional {
		t.Error("field_b should remain Optional")
	}
	if !attrs[2].Required {
		t.Error("field_c should remain Required")
	}
}

func TestResolveRef(t *testing.T) {
	spec := &openapi.Spec{
		Components: openapi.Components{
			Schemas: map[string]openapi.Schema{
				"MyType": {Type: "object", Description: "test"},
			},
		},
	}

	// Resolve from spec
	result := ResolveRef("#/components/schemas/MyType", spec)
	if result.Type != "object" {
		t.Errorf("ResolveRef from spec: Type = %q, want %q", result.Type, "object")
	}

	// Resolve from cache
	SchemaCache["CachedType"] = openapi.Schema{Type: "integer"}
	result = ResolveRef("#/components/schemas/CachedType", spec)
	if result.Type != "integer" {
		t.Errorf("ResolveRef from cache: Type = %q, want %q", result.Type, "integer")
	}
	delete(SchemaCache, "CachedType") // cleanup

	// Unresolvable ref
	result = ResolveRef("#/components/schemas/Unknown", spec)
	if result.Type != "string" {
		t.Errorf("ResolveRef unknown: Type = %q, want %q", result.Type, "string")
	}
}

func TestConvertToTerraformAttribute_String(t *testing.T) {
	spec := &openapi.Spec{
		Components: openapi.Components{
			Schemas: map[string]openapi.Schema{},
		},
	}
	s := openapi.Schema{
		Type:        "string",
		Description: "A test field",
	}
	attr := ConvertToTerraformAttribute("test_field", s, true, "", spec)
	if attr.Type != "string" {
		t.Errorf("Type = %q, want %q", attr.Type, "string")
	}
	if !attr.Required {
		t.Error("Should be Required")
	}
	if attr.GoType != "string" {
		t.Errorf("GoType = %q, want %q", attr.GoType, "string")
	}
}

func TestConvertToTerraformAttribute_Object(t *testing.T) {
	spec := &openapi.Spec{
		Components: openapi.Components{
			Schemas: map[string]openapi.Schema{},
		},
	}
	s := openapi.Schema{
		Type: "object",
		Properties: map[string]openapi.Schema{
			"nested": {Type: "string"},
		},
	}
	attr := ConvertToTerraformAttribute("my_block", s, false, "", spec)
	if !attr.IsBlock {
		t.Error("Object with properties should be a block")
	}
	if attr.NestedBlockType != "single" {
		t.Errorf("NestedBlockType = %q, want %q", attr.NestedBlockType, "single")
	}
}

func TestConvertToTerraformAttribute_Array(t *testing.T) {
	spec := &openapi.Spec{
		Components: openapi.Components{
			Schemas: map[string]openapi.Schema{},
		},
	}
	s := openapi.Schema{
		Type: "array",
		Items: &openapi.Schema{
			Type: "string",
		},
	}
	attr := ConvertToTerraformAttribute("tags", s, false, "", spec)
	if attr.Type != "list" {
		t.Errorf("Type = %q, want %q", attr.Type, "list")
	}
	if attr.ElementType != "string" {
		t.Errorf("ElementType = %q, want %q", attr.ElementType, "string")
	}
	if attr.IsBlock {
		t.Error("Array of strings should not be a block")
	}
}

func TestConvertToTerraformAttribute_DescriptionReserved(t *testing.T) {
	spec := &openapi.Spec{
		Components: openapi.Components{
			Schemas: map[string]openapi.Schema{},
		},
	}
	s := openapi.Schema{Type: "string", Description: "A description field"}
	attr := ConvertToTerraformAttribute("description", s, false, "", spec)
	if attr.GoName != "DescriptionSpec" {
		t.Errorf("GoName = %q, want %q", attr.GoName, "DescriptionSpec")
	}
	if attr.TfsdkTag != "description_spec" {
		t.Errorf("TfsdkTag = %q, want %q", attr.TfsdkTag, "description_spec")
	}
}

func TestConvertToTerraformAttribute_RequiredForCreate(t *testing.T) {
	schema := openapi.Schema{
		Type:        "string",
		Description: "A domain name",
		XF5XCRequiredFor: openapi.RequiredFor{
			Create: true,
		},
	}
	spec := &openapi.Spec{Components: openapi.Components{Schemas: map[string]openapi.Schema{}}}
	attr := ConvertToTerraformAttributeWithDepth("domain", schema, false, "", spec, 0, "domain")
	if !attr.Required {
		t.Error("Expected Required=true when x-f5xc-required-for.create=true at depth 0")
	}
	if attr.Optional {
		t.Error("Expected Optional=false when Required=true")
	}
}

func TestConvertToTerraformAttribute_RequiredForCreate_NestedIgnored(t *testing.T) {
	schema := openapi.Schema{
		Type:        "string",
		Description: "A nested field",
		XF5XCRequiredFor: openapi.RequiredFor{Create: true},
	}
	spec := &openapi.Spec{Components: openapi.Components{Schemas: map[string]openapi.Schema{}}}
	attr := ConvertToTerraformAttributeWithDepth("nested_field", schema, false, "", spec, 1, "parent.nested_field")
	if attr.Required {
		t.Error("x-f5xc-required-for.create should NOT set Required at depth > 0")
	}
}

func TestConvertToTerraformAttribute_ServerDefault(t *testing.T) {
	schema := openapi.Schema{
		Type:               "string",
		Description:        "A server-defaulted field",
		XF5XCServerDefault: true,
	}
	spec := &openapi.Spec{Components: openapi.Components{Schemas: map[string]openapi.Schema{}}}
	attr := ConvertToTerraformAttributeWithDepth("mode", schema, false, "", spec, 0, "mode")
	if !attr.Computed {
		t.Error("Expected Computed=true when x-f5xc-server-default=true")
	}
	if !attr.Optional {
		t.Error("Expected Optional=true for server-default fields")
	}
	if attr.PlanModifier != "UseStateForUnknown" {
		t.Errorf("Expected PlanModifier=UseStateForUnknown, got %q", attr.PlanModifier)
	}
}

func TestConvertToTerraformAttribute_ServerDefault_RequiredNotOverridden(t *testing.T) {
	schema := openapi.Schema{
		Type:               "string",
		Description:        "Required but server has default",
		XF5XCServerDefault: true,
	}
	spec := &openapi.Spec{Components: openapi.Components{Schemas: map[string]openapi.Schema{}}}
	attr := ConvertToTerraformAttributeWithDepth("name", schema, true, "", spec, 0, "name")
	if !attr.Computed {
		t.Error("Expected Computed=true for server-default field")
	}
	if !attr.Required {
		t.Error("Required should not be overridden when field is in OpenAPI required array")
	}
}

func TestExtractNestedAttributes(t *testing.T) {
	spec := &openapi.Spec{
		Components: openapi.Components{
			Schemas: map[string]openapi.Schema{},
		},
	}
	s := openapi.Schema{
		Properties: map[string]openapi.Schema{
			"port":     {Type: "integer"},
			"protocol": {Type: "string"},
		},
		Required: []string{"port"},
	}
	attrs := ExtractNestedAttributes(s, spec, 1, "parent")
	if len(attrs) != 2 {
		t.Fatalf("ExtractNestedAttributes returned %d attrs, want 2", len(attrs))
	}
	// Required should be first
	if attrs[0].Name != "port" {
		t.Errorf("First attr should be required 'port', got %q", attrs[0].Name)
	}
}

func TestMaxNestedDepthGuard(t *testing.T) {
	spec := &openapi.Spec{
		Components: openapi.Components{
			Schemas: map[string]openapi.Schema{},
		},
	}
	s := openapi.Schema{
		Properties: map[string]openapi.Schema{
			"deep": {Type: "string"},
		},
	}
	attrs := ExtractNestedAttributes(s, spec, MaxNestedDepth+1, "")
	if attrs != nil {
		t.Error("ExtractNestedAttributes should return nil when depth exceeds MaxNestedDepth")
	}
}
