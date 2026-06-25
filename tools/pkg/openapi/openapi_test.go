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

// =============================================================================
// Test helpers
// =============================================================================

func assertEqual(t *testing.T, field string, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, want %v", field, got, want)
	}
}

func assertTrue(t *testing.T, field string, got bool) {
	t.Helper()
	if !got {
		t.Errorf("%s = false, want true", field)
	}
}

func assertFalse(t *testing.T, field string, got bool) {
	t.Helper()
	if got {
		t.Errorf("%s = true, want false", field)
	}
}

func assertNotNil(t *testing.T, field string, got interface{}) {
	t.Helper()
	if got == nil {
		t.Errorf("%s is nil, want non-nil", field)
	}
}

func assertLen(t *testing.T, field string, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("%s length = %d, want %d", field, got, want)
	}
}

func assertIntEqual(t *testing.T, field string, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %d, want %d", field, got, want)
	}
}

// =============================================================================
// SP-1 Extension Parsing Tests
// =============================================================================

func TestSchema_AllExtensions(t *testing.T) {
	testFile := filepath.Join("testdata", "full-extensions.json")
	spec, err := ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile(%q) error = %v", testFile, err)
	}

	schema, ok := spec.Components.Schemas["FullExtensionsSchema"]
	if !ok {
		t.Fatal("FullExtensionsSchema not found in parsed spec")
	}

	// --- Standard OpenAPI fields ---
	assertEqual(t, "Type", schema.Type, "object")
	assertEqual(t, "Description", schema.Description, "Schema with all extension fields populated")
	assertEqual(t, "Title", schema.Title, "Full Extensions Schema")
	assertLen(t, "Required", len(schema.Required), 1)
	assertEqual(t, "Required[0]", schema.Required[0], "name")

	// --- Original x-ves-* extensions ---
	assertEqual(t, "XDisplayName", schema.XDisplayName, "Full Extensions")
	assertEqual(t, "XVesExample", schema.XVesExample, "ves-example-value")
	assertLen(t, "XVesValidationRules", len(schema.XVesValidationRules), 2)
	assertEqual(t, "XVesProtoMessage", schema.XVesProtoMessage, "ves.io.schema.full_extensions.Object")

	// --- Existing x-f5xc-* extensions ---
	assertEqual(t, "XF5XCCategory", schema.XF5XCCategory, "security")
	assertEqual(t, "XF5XCRequiresTier", schema.XF5XCRequiresTier, "Standard")
	assertEqual(t, "XF5XCComplexity", schema.XF5XCComplexity, "advanced")
	assertEqual(t, "XF5XCExample", schema.XF5XCExample, "f5xc-example-value")
	assertEqual(t, "XF5XCDescriptionShort", schema.XF5XCDescriptionShort, "Short description")
	assertEqual(t, "XF5XCDescriptionMed", schema.XF5XCDescriptionMed, "Medium length description")
	assertLen(t, "XF5XCUseCases", len(schema.XF5XCUseCases), 2)
	assertLen(t, "XF5XCRelatedDomains", len(schema.XF5XCRelatedDomains), 2)
	assertTrue(t, "XF5XCIsPreview", schema.XF5XCIsPreview)
	assertNotNil(t, "XF5XCNamespaceProfile", schema.XF5XCNamespaceProfile)
	assertEqual(t, "XF5XCNamespaceProfile.Constraint.Allowed[0]", schema.XF5XCNamespaceProfile.Constraint.Allowed[0], "system")
	assertTrue(t, "XF5XCNamespaceProfile.Constraint.Enforced", schema.XF5XCNamespaceProfile.Constraint.Enforced)
	assertEqual(t, "XF5XCNamespaceProfile.Recommendation.Primary", schema.XF5XCNamespaceProfile.Recommendation.Primary, "system")
	assertEqual(t, "XF5XCNamespaceProfile.Classification.Category", schema.XF5XCNamespaceProfile.Classification.Category, "infrastructure")
	assertEqual(t, "XF5XCNamespaceProfile.Classification.MultiTenantPattern", schema.XF5XCNamespaceProfile.Classification.MultiTenantPattern, "none")
	assertEqual(t, "XF5XCIcon", schema.XF5XCIcon, "shield-icon")
	assertEqual(t, "XF5XCCLIDomain", schema.XF5XCCLIDomain, "security")

	// --- Additional upstream extensions ---
	assertEqual(t, "XVesDeprecated", schema.XVesDeprecated, "Use v2 endpoint instead")
	assertEqual(t, "XVesDisplayOrder", schema.XVesDisplayOrder, "5")
	assertEqual(t, "XVesProtoEnum", schema.XVesProtoEnum, "ves.io.schema.SomeEnum")
	assertEqual(t, "XVesRequired", schema.XVesRequired, "true")
	assertTrue(t, "XRequired", schema.XRequired)

	// --- Enrichment — actionable ---
	assertLen(t, "XF5XCConflictsWith", len(schema.XF5XCConflictsWith), 2)
	assertNotNil(t, "XF5XCConstraints", schema.XF5XCConstraints)
	assertEqual(t, "XFieldMutability", schema.XFieldMutability, "immutable")
	assertIntEqual(t, "XOriginalMaxLength", schema.XOriginalMaxLength, 256)

	// --- Provenance ---
	assertEqual(t, "XReconciledAt", schema.XReconciledAt, "2026-04-23T09:16:13Z")
	assertTrue(t, "XReconciledFromDiscovery", schema.XReconciledFromDiscovery)

	// --- SP-1: field-level extensions ---
	assertTrue(t, "XF5XCServerDefault", schema.XF5XCServerDefault)
	assertTrue(t, "XF5XCRequiredFor.MinimumConfig", schema.XF5XCRequiredFor.MinimumConfig)
	assertTrue(t, "XF5XCRequiredFor.Create", schema.XF5XCRequiredFor.Create)
	assertFalse(t, "XF5XCRequiredFor.Update", schema.XF5XCRequiredFor.Update)
	assertTrue(t, "XF5XCRequiredFor.Read", schema.XF5XCRequiredFor.Read)
	assertNotNil(t, "XF5XCRecommendedValue", schema.XF5XCRecommendedValue)
	assertNotNil(t, "XF5XCMinimumConfiguration", schema.XF5XCMinimumConfiguration)

	// --- SP-1: schema-level extensions ---
	assertNotNil(t, "XF5XCValidation", schema.XF5XCValidation)
	assertNotNil(t, "XF5XCDefaults", schema.XF5XCDefaults)
	assertNotNil(t, "XF5XCConditions", schema.XF5XCConditions)
	assertEqual(t, "XF5XCDeprecated", schema.XF5XCDeprecated, "Deprecated: use new_field instead")
	assertNotNil(t, "XF5XCCompletion", schema.XF5XCCompletion)
	assertEqual(t, "XF5XCDisplayName", schema.XF5XCDisplayName, "Full Extensions Display Name")
	assertEqual(t, "XF5XCDescription", schema.XF5XCDescription, "Full enriched description for display")
	assertLen(t, "XF5XCExamples", len(schema.XF5XCExamples), 3)
	assertNotNil(t, "XF5XCRequiredForOps", schema.XF5XCRequiredForOps)
	assertIntEqual(t, "XF5XCDisplayOrder", schema.XF5XCDisplayOrder, 10)
	assertEqual(t, "XF5XCUniqueness", schema.XF5XCUniqueness, "global")
	assertEqual(t, "XF5XCTerraformResource", schema.XF5XCTerraformResource, "xcsh_full_extensions")

	// --- SP-1: operation-level extensions ---
	assertEqual(t, "XF5XCDangerLevel", schema.XF5XCDangerLevel, "high")
	assertTrue(t, "XF5XCConfirmationRequired", schema.XF5XCConfirmationRequired)
	assertLen(t, "XF5XCSideEffects", len(schema.XF5XCSideEffects), 2)
	assertNotNil(t, "XF5XCOperationMetadata", schema.XF5XCOperationMetadata)
	assertLen(t, "XF5XCRequiredFields", len(schema.XF5XCRequiredFields), 2)

	// --- Nested property extensions ---
	nameProp, ok := schema.Properties["name"]
	if !ok {
		t.Fatal("name property not found")
	}
	assertFalse(t, "name.XF5XCServerDefault", nameProp.XF5XCServerDefault)
	assertTrue(t, "name.XF5XCRequiredFor.MinimumConfig", nameProp.XF5XCRequiredFor.MinimumConfig)
	assertTrue(t, "name.XF5XCRequiredFor.Create", nameProp.XF5XCRequiredFor.Create)
	assertFalse(t, "name.XF5XCRequiredFor.Update", nameProp.XF5XCRequiredFor.Update)
	assertFalse(t, "name.XF5XCRequiredFor.Read", nameProp.XF5XCRequiredFor.Read)
	assertNotNil(t, "name.XF5XCRecommendedValue", nameProp.XF5XCRecommendedValue)
	assertNotNil(t, "name.XF5XCMinimumConfiguration", nameProp.XF5XCMinimumConfiguration)
}

func TestTerraformAttribute_AllFields(t *testing.T) {
	attr := TerraformAttribute{
		// Existing fields
		Name:        "test_field",
		GoName:      "TestField",
		TfsdkTag:    "test_field",
		Type:        "string",
		Description: "A test field",
		Required:    true,
		Optional:    false,
		Computed:    false,
		Sensitive:   false,
		IsBlock:     false,
		IsSpecField: true,
		JsonName:    "test_field",
		GoType:      "string",

		// SP-1 additions
		ServerDefault:        true,
		Default:              "default-value",
		MinimumConfigRequired: true,
		RecommendedValue:     "recommended-value",
		ValidationRules:      map[string]string{"max_length": "256"},
		Complexity:           "advanced",
		UseCases:             []string{"use-case-1"},
		DeprecationMessage:   "Deprecated: use new_field",
		ConflictsWith:        []string{"other_field"},
		MaxLength:            256,
		Immutable:            true,
		EnumValues:           []string{"a", "b", "c"},
		MinLength:            1,
		Pattern:              "^[a-z]+$",
		MinItems:             1,
		MaxItems:             10,
	}

	// Existing fields
	assertEqual(t, "Name", attr.Name, "test_field")
	assertEqual(t, "GoName", attr.GoName, "TestField")
	assertEqual(t, "Type", attr.Type, "string")
	assertTrue(t, "Required", attr.Required)
	assertTrue(t, "IsSpecField", attr.IsSpecField)

	// SP-1 additions
	assertTrue(t, "ServerDefault", attr.ServerDefault)
	assertEqual(t, "Default", attr.Default, "default-value")
	assertTrue(t, "MinimumConfigRequired", attr.MinimumConfigRequired)
	assertEqual(t, "RecommendedValue", attr.RecommendedValue, "recommended-value")
	assertLen(t, "ValidationRules", len(attr.ValidationRules), 1)
	assertEqual(t, "Complexity", attr.Complexity, "advanced")
	assertLen(t, "UseCases", len(attr.UseCases), 1)
	assertEqual(t, "DeprecationMessage", attr.DeprecationMessage, "Deprecated: use new_field")
	assertLen(t, "ConflictsWith", len(attr.ConflictsWith), 1)
	assertIntEqual(t, "MaxLength", attr.MaxLength, 256)
	assertTrue(t, "Immutable", attr.Immutable)
	assertLen(t, "EnumValues", len(attr.EnumValues), 3)
	assertIntEqual(t, "MinLength", attr.MinLength, 1)
	assertEqual(t, "Pattern", attr.Pattern, "^[a-z]+$")
	assertIntEqual(t, "MinItems", attr.MinItems, 1)
	assertIntEqual(t, "MaxItems", attr.MaxItems, 10)
}

func TestResourceTemplate_AllFields(t *testing.T) {
	tmpl := ResourceTemplate{
		// Existing fields
		Name:                   "test_resource",
		TitleCase:              "TestResource",
		APIPath:                "/api/config/namespaces/{namespace}/test_resources",
		HasNamespaceInPath:     true,
		Description:            "A test resource",
		UsesBoolPlanModifier:   true,
		UsesInt64PlanModifier:  false,
		UsesStringPlanModifier: true,

		// SP-1 additions
		HasBlocks:              true,
		HasMaxLengthValidators: true,
		HasEnumValidators:      true,
		HasPatternValidators:   false,
		HasListSizeValidators:  true,
		HasConflicts:           true,
		ConflictCheckCode:      "// conflict check code here",
	}

	// Existing fields
	assertEqual(t, "Name", tmpl.Name, "test_resource")
	assertEqual(t, "TitleCase", tmpl.TitleCase, "TestResource")
	assertTrue(t, "HasNamespaceInPath", tmpl.HasNamespaceInPath)
	assertTrue(t, "UsesBoolPlanModifier", tmpl.UsesBoolPlanModifier)
	assertFalse(t, "UsesInt64PlanModifier", tmpl.UsesInt64PlanModifier)

	// SP-1 additions
	assertTrue(t, "HasBlocks", tmpl.HasBlocks)
	assertTrue(t, "HasMaxLengthValidators", tmpl.HasMaxLengthValidators)
	assertTrue(t, "HasEnumValidators", tmpl.HasEnumValidators)
	assertFalse(t, "HasPatternValidators", tmpl.HasPatternValidators)
	assertTrue(t, "HasListSizeValidators", tmpl.HasListSizeValidators)
	assertTrue(t, "HasConflicts", tmpl.HasConflicts)
	assertEqual(t, "ConflictCheckCode", tmpl.ConflictCheckCode, "// conflict check code here")
}

func TestPrimaryResourceMetadata_AllFields(t *testing.T) {
	meta := PrimaryResourceMetadata{
		// Existing fields
		Name:             "http_loadbalancer",
		Description:      "HTTP Load Balancer",
		DescriptionShort: "HTTP LB",
		Tier:             "Standard",
		Icon:             "lb-icon",
		Category:         "networking",
		SupportsLogs:     true,
		SupportsMetrics:  true,
		Dependencies: ResourceDependencies{
			Required: []string{"namespace"},
			Optional: []string{"waf"},
		},
		RelationshipHints: []string{"origin_pool"},

		// SP-1 additions
		SchemaComponents: []string{"ves.io.schema.views.http_loadbalancer.Object"},
		APIPaths:         []string{"/api/config/namespaces/{namespace}/http_loadbalancers"},
	}

	// Existing fields
	assertEqual(t, "Name", meta.Name, "http_loadbalancer")
	assertEqual(t, "Tier", meta.Tier, "Standard")
	assertTrue(t, "SupportsLogs", meta.SupportsLogs)
	assertLen(t, "Dependencies.Required", len(meta.Dependencies.Required), 1)
	assertLen(t, "RelationshipHints", len(meta.RelationshipHints), 1)

	// SP-1 additions
	assertLen(t, "SchemaComponents", len(meta.SchemaComponents), 1)
	assertEqual(t, "SchemaComponents[0]", meta.SchemaComponents[0], "ves.io.schema.views.http_loadbalancer.Object")
	assertLen(t, "APIPaths", len(meta.APIPaths), 1)
	assertEqual(t, "APIPaths[0]", meta.APIPaths[0], "/api/config/namespaces/{namespace}/http_loadbalancers")
}

func TestDomainInfo_SpecLevelExtensions(t *testing.T) {
	specJSON := `{
		"openapi": "3.0.0",
		"info": {
			"title": "Security Domain",
			"description": "Security domain spec",
			"version": "1.0.0",
			"x-f5xc-description-short": "Security",
			"x-f5xc-description-medium": "Security domain for F5 XC",
			"x-f5xc-icon": "shield",
			"x-f5xc-logo-svg": "<svg>logo</svg>",
			"x-f5xc-description-long": "A comprehensive security domain",
			"x-f5xc-summary": "Security summary",
			"x-f5xc-namespace-profile": {
				"constraint": {"allowed": ["system"], "enforced": true},
				"recommendation": {"primary": "system"},
				"classification": {"category": "infrastructure", "multi_tenant_pattern": "none"}
			},
			"x-f5xc-cli-metadata": {"command_prefix": "security"},
			"x-f5xc-glossary": {"WAF": "Web Application Firewall"},
			"x-f5xc-guided-workflows": [{"name": "deploy_waf", "steps": []}],
			"x-f5xc-acronyms": {"LB": "Load Balancer", "WAF": "Web Application Firewall"}
		},
		"paths": {},
		"components": {"schemas": {}},
		"x-f5xc-category": "security",
		"x-f5xc-doc-section": "security-services",
		"x-f5xc-primary-resources": [
			{
				"name": "http_loadbalancer",
				"description": "HTTP LB",
				"tier": "Standard",
				"schema_components": ["ves.io.schema.views.http_loadbalancer.Object"],
				"api_paths": ["/api/config/namespaces/{namespace}/http_loadbalancers"]
			}
		],
		"x-f5xc-critical-resources": ["http_loadbalancer"],
		"x-f5xc-logo-svg": "<svg>domain-logo</svg>",
		"x-f5xc-icon": "domain-shield"
	}`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "security.json")
	if err := os.WriteFile(testFile, []byte(specJSON), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	spec, err := ParseDomainSpec(testFile)
	if err != nil {
		t.Fatalf("ParseDomainSpec() error = %v", err)
	}

	// DomainInfo existing fields
	assertEqual(t, "Info.Title", spec.Info.Title, "Security Domain")
	assertEqual(t, "Info.XF5XCDescriptionShort", spec.Info.XF5XCDescriptionShort, "Security")
	assertEqual(t, "Info.XF5XCIcon", spec.Info.XF5XCIcon, "shield")
	assertEqual(t, "Info.XF5XCLogoSVG", spec.Info.XF5XCLogoSVG, "<svg>logo</svg>")
	assertNotNil(t, "Info.XF5XCNamespaceProfile", spec.Info.XF5XCNamespaceProfile)
	assertEqual(t, "Info.XF5XCNamespaceProfile.Recommendation.Primary", spec.Info.XF5XCNamespaceProfile.Recommendation.Primary, "system")

	// DomainInfo SP-1 additions
	assertNotNil(t, "Info.XF5XCCLIMetadata", spec.Info.XF5XCCLIMetadata)
	assertNotNil(t, "Info.XF5XCGlossary", spec.Info.XF5XCGlossary)
	assertLen(t, "Info.XF5XCGuidedWorkflows", len(spec.Info.XF5XCGuidedWorkflows), 1)
	assertLen(t, "Info.XF5XCAcronyms", len(spec.Info.XF5XCAcronyms), 2)

	// DomainSpec SP-1 additions
	assertEqual(t, "XF5XCDocSection", spec.XF5XCDocSection, "security-services")
	assertLen(t, "XF5XCPrimaryResources", len(spec.XF5XCPrimaryResources), 1)
	assertEqual(t, "XF5XCPrimaryResources[0].Name", spec.XF5XCPrimaryResources[0].Name, "http_loadbalancer")
	assertLen(t, "XF5XCPrimaryResources[0].SchemaComponents", len(spec.XF5XCPrimaryResources[0].SchemaComponents), 1)
	assertLen(t, "XF5XCPrimaryResources[0].APIPaths", len(spec.XF5XCPrimaryResources[0].APIPaths), 1)
	assertLen(t, "XF5XCCriticalResources", len(spec.XF5XCCriticalResources), 1)
	assertEqual(t, "XF5XCCriticalResources[0]", spec.XF5XCCriticalResources[0], "http_loadbalancer")
	assertEqual(t, "XF5XCLogoSVG", spec.XF5XCLogoSVG, "<svg>domain-logo</svg>")
	assertEqual(t, "XF5XCIcon", spec.XF5XCIcon, "domain-shield")
}

func TestSchema_GetBestDescription_WithEnrichedDescription(t *testing.T) {
	tests := []struct {
		name     string
		schema   Schema
		expected string
	}{
		{
			"x-f5xc-description takes priority",
			Schema{
				Description:           "Original",
				XF5XCDescriptionShort: "Short",
				XF5XCDescriptionMed:   "Medium",
				XF5XCDescription:      "Full enriched description",
			},
			"Full enriched description",
		},
		{
			"falls back to medium when description absent",
			Schema{
				Description:           "Original",
				XF5XCDescriptionShort: "Short",
				XF5XCDescriptionMed:   "Medium",
			},
			"Medium",
		},
		{
			"falls back to short when medium absent",
			Schema{
				Description:           "Original",
				XF5XCDescriptionShort: "Short",
			},
			"Short",
		},
		{
			"falls back to original when all enriched absent",
			Schema{
				Description: "Original",
			},
			"Original",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.schema.GetBestDescription()
			if got != tt.expected {
				t.Errorf("GetBestDescription() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGenerationResult_AllFields(t *testing.T) {
	result := GenerationResult{
		ResourceName: "test_resource",
		Success:      true,
		Error:        "",
		FilePath:     "/path/to/resource.go",
		AttrCount:    15,
		BlockCount:   3,
	}

	assertEqual(t, "ResourceName", result.ResourceName, "test_resource")
	assertTrue(t, "Success", result.Success)
	assertEqual(t, "FilePath", result.FilePath, "/path/to/resource.go")
	assertIntEqual(t, "AttrCount", result.AttrCount, 15)
	assertIntEqual(t, "BlockCount", result.BlockCount, 3)
}
