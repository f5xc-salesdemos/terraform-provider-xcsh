// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

// Package schema provides functions for converting OpenAPI schemas to Terraform
// resource attributes. It handles type mapping, nested attribute extraction,
// and schema resolution for the F5 XC Terraform provider code generation tools.
package schema

import (
	"fmt"
	"sort"
	"strings"

	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/constraints"
	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/description"
	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/naming"
	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/openapi"
)

// MaxNestedDepth is the maximum recursion depth for nested schemas to prevent infinite loops.
// Set high enough to capture all nested defaults while still preventing infinite recursion.
const MaxNestedDepth = 20

// SchemaCache stores parsed schema definitions keyed by schema name.
// Populated by parseOpenAPISpec in the monolith; used by ResolveRef for cross-file lookups.
var SchemaCache = make(map[string]openapi.Schema)

// RawSpecCache stores raw JSON schema maps for extracting x-ves-oneof-field annotations.
// Populated by parseOpenAPISpec in the monolith; used by ExtractOneOfGroups.
var RawSpecCache = make(map[string]map[string]interface{})

// MapSchemaType converts an OpenAPI type string to a Terraform attribute type string.
func MapSchemaType(t string) string {
	switch t {
	case "string":
		return "string"
	case "integer", "number":
		return "int64"
	case "boolean":
		return "bool"
	default:
		return "string"
	}
}

// MapSchemaTypeToGo converts OpenAPI types to Go types for client struct generation.
func MapSchemaTypeToGo(t string) string {
	switch t {
	case "string":
		return "string"
	case "integer", "number":
		return "int64"
	case "boolean":
		return "bool"
	default:
		return "interface{}"
	}
}

// IsMetadataField returns true if the field is a metadata field (not a spec field).
func IsMetadataField(name string) bool {
	metadataFields := map[string]bool{
		"name":        true,
		"namespace":   true,
		"labels":      true,
		"annotations": true,
		"description": true,
		"id":          true,
		"timeouts":    true,
	}
	return metadataFields[name]
}

// FilterSpecFields returns only the spec fields (non-metadata) from attributes.
func FilterSpecFields(attrs []openapi.TerraformAttribute) []openapi.TerraformAttribute {
	var specFields []openapi.TerraformAttribute
	for _, attr := range attrs {
		if !IsMetadataField(attr.TfsdkTag) && attr.IsSpecField {
			specFields = append(specFields, attr)
		}
	}
	return specFields
}

// FilterOptional returns only optional (non-required, non-computed) attributes.
func FilterOptional(attrs []openapi.TerraformAttribute) []openapi.TerraformAttribute {
	var result []openapi.TerraformAttribute
	for _, attr := range attrs {
		if attr.Optional && !attr.Required && !attr.Computed {
			result = append(result, attr)
		}
	}
	return result
}

// ParseMinConfigRequiredFields extracts the required_fields list from x-f5xc-minimum-configuration.
// The structure is: {"required_fields": ["field_name", "spec.field_name", ...], "example_yaml": "..."}
// Strips "spec." prefix which commonly appears in minimum config field names.
func ParseMinConfigRequiredFields(raw interface{}) []string {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	fieldsRaw, ok := m["required_fields"].([]interface{})
	if !ok {
		return nil
	}
	var result []string
	for _, f := range fieldsRaw {
		if s, ok := f.(string); ok && s != "" {
			result = append(result, strings.TrimPrefix(s, "spec."))
		}
	}
	return result
}

// PromoteMinConfigRequired marks attributes as Required if their tfsdk tag appears in the minimum config.
// This promotes Optional fields that are needed for a minimum viable configuration.
func PromoteMinConfigRequired(attrs []openapi.TerraformAttribute, minFields map[string]bool) {
	for i := range attrs {
		if minFields[attrs[i].TfsdkTag] && attrs[i].Optional && !attrs[i].Required {
			attrs[i].Required = true
			attrs[i].Optional = false
		}
	}
}

// ResolveRef resolves a $ref string to a schema definition.
// It first checks the spec's own components, then falls back to the global SchemaCache.
func ResolveRef(ref string, spec *openapi.Spec) openapi.Schema {
	parts := strings.Split(ref, "/")
	schemaName := parts[len(parts)-1]

	if schema, ok := spec.Components.Schemas[schemaName]; ok {
		return schema
	}
	if schema, ok := SchemaCache[schemaName]; ok {
		return schema
	}
	return openapi.Schema{Type: "string"}
}

// ConvertToTerraformAttribute converts an OpenAPI schema property to a TerraformAttribute.
// This is the top-level entry point that delegates to ConvertToTerraformAttributeWithDepth.
func ConvertToTerraformAttribute(name string, schema openapi.Schema, required bool, oneOfGroup string, spec *openapi.Spec) openapi.TerraformAttribute {
	return ConvertToTerraformAttributeWithDepth(name, schema, required, oneOfGroup, spec, 0, name)
}

// ConvertToTerraformAttributeWithDepth converts an OpenAPI schema property to a TerraformAttribute,
// tracking recursion depth to prevent infinite loops on deeply nested schemas.
func ConvertToTerraformAttributeWithDepth(name string, schema openapi.Schema, required bool, oneOfGroup string, spec *openapi.Spec, depth int, fieldPath string) openapi.TerraformAttribute {
	// Preserve extensions from original schema before resolving $ref
	// Extensions like x-f5xc-server-default are on the property, not the referenced schema
	serverDefault := schema.XF5XCServerDefault
	defaultValue := schema.Default // Preserve actual default value
	descShort := schema.XF5XCDescriptionShort
	descMed := schema.XF5XCDescriptionMed
	requiredFor := schema.XF5XCRequiredFor
	recommendedValue := schema.XF5XCRecommendedValue
	validationRules := schema.XVesValidationRules
	complexity := schema.XF5XCComplexity
	useCases := schema.XF5XCUseCases

	if schema.Ref != "" {
		schema = ResolveRef(schema.Ref, spec)
		// Restore preserved extensions (property-level extensions take precedence)
		if serverDefault {
			schema.XF5XCServerDefault = serverDefault
		}
		if defaultValue != nil && schema.Default == nil {
			schema.Default = defaultValue
		}
		if descShort != "" && schema.XF5XCDescriptionShort == "" {
			schema.XF5XCDescriptionShort = descShort
		}
		if descMed != "" && schema.XF5XCDescriptionMed == "" {
			schema.XF5XCDescriptionMed = descMed
		}
		// Restore enhanced metadata extensions
		if requiredFor.MinimumConfig || requiredFor.Create || requiredFor.Update || requiredFor.Read {
			schema.XF5XCRequiredFor = requiredFor
		}
		if recommendedValue != nil {
			schema.XF5XCRecommendedValue = recommendedValue
		}
		if len(validationRules) > 0 && len(schema.XVesValidationRules) == 0 {
			schema.XVesValidationRules = validationRules
		}
		if complexity != "" && schema.XF5XCComplexity == "" {
			schema.XF5XCComplexity = complexity
		}
		if len(useCases) > 0 && len(schema.XF5XCUseCases) == 0 {
			schema.XF5XCUseCases = useCases
		}
	}

	// Convert name to valid Terraform attribute name (lowercase with underscores)
	tfsdkName := naming.ToSnakeCaseTerraform(name)

	// Handle reserved field names to avoid conflicts with standard metadata fields
	// "description" is used for resource metadata in TF schema, so rename spec fields
	goName := naming.ToResourceTypeName(name)
	if strings.ToLower(name) == "description" {
		goName = "DescriptionSpec"
		tfsdkName = "description_spec"
	}

	attr := openapi.TerraformAttribute{
		Name:        name,
		GoName:      goName,
		TfsdkTag:    tfsdkName,
		Required:    required,
		Optional:    !required,
		OneOfGroup:  oneOfGroup,
		MaxDepth:      depth,
		IsSpecField:   true, // Attributes from OpenAPI spec are spec fields
		JsonName:      name, // Original OpenAPI property name for JSON marshaling
		ServerDefault: schema.XF5XCServerDefault,
		Default:       schema.Default, // Store actual default value for metadata
		// Enhanced metadata from enriched API specs
		MinimumConfigRequired: schema.XF5XCRequiredFor.MinimumConfig,
		RecommendedValue:      schema.XF5XCRecommendedValue,
		ValidationRules:       schema.XVesValidationRules,
		Complexity:            schema.XF5XCComplexity,
		UseCases:              schema.XF5XCUseCases,
		DeprecationMessage:    schema.XVesDeprecated,
		ConflictsWith:         schema.XF5XCConflictsWith,
		MaxLength:             schema.XOriginalMaxLength,
		Immutable:             schema.XFieldMutability == "immutable" || schema.XFieldMutability == "create-only",
	}

	// Apply x-f5xc-constraints when confidence/determinism thresholds are met
	if c := constraints.Parse(schema.XF5XCConstraints); c != nil {
		if c.MinLength > 0 {
			attr.MinLength = c.MinLength
		}
		if c.MaxLength > 0 && attr.MaxLength == 0 {
			attr.MaxLength = c.MaxLength
		}
		if c.Pattern != "" {
			attr.Pattern = c.Pattern
		}
		if c.MinItems > 0 {
			attr.MinItems = c.MinItems
		}
		if c.MaxItems > 0 {
			attr.MaxItems = c.MaxItems
		}
	}

	// Prefer x-validation-rules over x-ves-validation-rules when both present
	if len(schema.XValidationRules) > 0 {
		attr.ValidationRules = schema.XValidationRules
	}
	// x-required / x-ves-required are additional Required signals, but only
	// apply at depth 0 (top-level resource spec attributes). For nested
	// attributes inside optional blocks, Required: true causes Terraform
	// Plugin Framework to fail validation even when the block is absent.
	if depth == 0 && (schema.XRequired || schema.XVesRequired == "true") {
		attr.Required = true
		attr.Optional = false
	}
	// Set PlanModifier for immutable fields (reuses existing template infrastructure)
	if attr.Immutable && attr.PlanModifier == "" {
		attr.PlanModifier = "RequiresReplace"
	}

	// Build description with enrichment extension priority:
	// 1. x-f5xc-description-medium (preferred - detailed but concise)
	// 2. x-f5xc-description-short (fallback - ultra-short)
	// 3. x-displayname + description (original behavior)
	// 4. description only
	fieldDesc := schema.XF5XCDescriptionMed
	if fieldDesc == "" {
		fieldDesc = schema.XF5XCDescriptionShort
	}
	if fieldDesc == "" {
		if schema.XDisplayName != "" && schema.Description != "" {
			fieldDesc = schema.XDisplayName + ". " + schema.Description
		} else if schema.XDisplayName != "" {
			fieldDesc = schema.XDisplayName
		} else {
			fieldDesc = schema.Description
		}
	}
	attr.Description = description.Clean(fieldDesc, fieldPath)
	if attr.Description == "" {
		attr.Description = fmt.Sprintf("Configuration for %s.", name)
	}

	// Format enum values per HashiCorp standards: "Possible values are `value1`, `value2`"
	if len(schema.Enum) > 0 {
		attr.Description = description.FormatEnum(attr.Description, schema.Enum)
	}

	// Populate EnumValues for string-type enum fields (for OneOf validator generation)
	if schema.Type == "string" && len(schema.Enum) > 0 {
		for _, v := range schema.Enum {
			if str, ok := v.(string); ok && str != "" && len(str) <= 50 {
				attr.EnumValues = append(attr.EnumValues, str)
			}
		}
	}

	// Format default values per HashiCorp standards: "Defaults to `value`."
	// First check for explicit default field, then try to extract from description text
	if schema.Default != nil {
		attr.Description = description.FormatDefault(attr.Description, schema.Default)
	} else {
		// Try to extract default from description text (F5 XC often documents defaults this way)
		extractedDefault, cleanedDesc := description.ExtractDefault(attr.Description)
		if extractedDefault != nil {
			attr.Description = description.FormatDefault(cleanedDesc, extractedDefault)
		}
	}

	// Add server default note to description (x-f5xc-server-default extension)
	// This indicates fields where F5XC server applies sensible defaults when omitted
	if schema.XF5XCServerDefault {
		attr.Description += " Server applies default when omitted."
	}

	// Determine type and extract nested attributes
	switch schema.Type {
	case "string":
		attr.Type = "string"
		attr.GoType = "string"
	case "integer", "number":
		attr.Type = "int64"
		attr.GoType = "int64"
	case "boolean":
		attr.Type = "bool"
		attr.GoType = "bool"
	case "array":
		attr.Type = "list"
		attr.GoType = "[]interface{}" // Default, refined below
		if schema.Items != nil {
			itemSchema := *schema.Items
			if itemSchema.Ref != "" {
				itemSchema = ResolveRef(itemSchema.Ref, spec)
			}
			if itemSchema.Type == "object" || len(itemSchema.Properties) > 0 {
				attr.IsBlock = true
				attr.NestedBlockType = "list"
				attr.GoType = "[]map[string]interface{}"
				// Extract nested attributes if within depth limit
				if depth < MaxNestedDepth {
					attr.NestedAttributes = ExtractNestedAttributes(itemSchema, spec, depth+1, fieldPath)
				}
			} else {
				attr.ElementType = MapSchemaType(itemSchema.Type)
				attr.GoType = "[]" + MapSchemaTypeToGo(itemSchema.Type)
				// Capture enum and default values from item schema for list element documentation
				if len(itemSchema.Enum) > 0 {
					attr.Description = description.FormatEnum(attr.Description, itemSchema.Enum)
				}
				if itemSchema.Default != nil {
					attr.Description = description.FormatDefault(attr.Description, itemSchema.Default)
				}
			}
		}
	case "object":
		if len(schema.Properties) > 0 {
			attr.Type = "object"
			attr.IsBlock = true
			attr.NestedBlockType = "single"
			attr.GoType = "map[string]interface{}"
			// Extract nested attributes if within depth limit
			if depth < MaxNestedDepth {
				attr.NestedAttributes = ExtractNestedAttributes(schema, spec, depth+1, fieldPath)
			}
		} else if schema.AdditionalProperties != nil {
			attr.Type = "map"
			attr.ElementType = "string"
			attr.GoType = "map[string]string"
		} else {
			attr.Type = "object"
			attr.IsBlock = true
			attr.NestedBlockType = "single"
			attr.GoType = "map[string]interface{}"
		}
	default:
		if len(schema.Properties) > 0 {
			attr.Type = "object"
			attr.IsBlock = true
			attr.NestedBlockType = "single"
			attr.GoType = "map[string]interface{}"
			// Extract nested attributes if within depth limit
			if depth < MaxNestedDepth {
				attr.NestedAttributes = ExtractNestedAttributes(schema, spec, depth+1, fieldPath)
			}
		} else {
			attr.Type = "string"
			attr.GoType = "string"
		}
	}

	// For optional scalar attributes (string, bool, int64), mark as computed with
	// UseStateForUnknown to prevent drift from API defaults. The API may return
	// default values for optional fields that weren't set in the config, causing
	// spurious plan diffs. UseStateForUnknown tells Terraform to preserve the
	// prior state value when the config doesn't specify a value.
	// Note: We only apply this to top-level attributes (depth == 0) because:
	// 1. Only top-level attributes are rendered with plan modifiers in the template
	// 2. Nested attributes inside blocks don't need this protection as the block
	//    structure itself handles the state management
	if depth == 0 && attr.Optional && !attr.IsBlock && !attr.Required {
		if attr.Type == "string" || attr.Type == "bool" || attr.Type == "int64" {
			attr.Computed = true
			attr.PlanModifier = "UseStateForUnknown"
		}
	}

	return attr
}

// ExtractNestedAttributes extracts attributes from an object schema's properties.
func ExtractNestedAttributes(schema openapi.Schema, spec *openapi.Spec, depth int, parentPath string) []openapi.TerraformAttribute {
	if depth > MaxNestedDepth {
		return nil
	}

	// Build the required set from the parent schema's "required" array.
	// Note: x-required / x-ves-required enrichment overrides are intentionally
	// restricted to depth 0 (see ConvertToTerraformAttributeWithDepth) so they
	// do not cascade here and cause "Missing required argument" errors when the
	// parent block is omitted from configuration. The OpenAPI spec's own
	// "required" array at nested levels is preserved here because it reflects
	// genuine schema constraints (e.g., a sub-object that truly requires certain
	// fields when the block is present).
	requiredSet := make(map[string]bool)
	for _, r := range schema.Required {
		requiredSet[r] = true
	}

	var attrs []openapi.TerraformAttribute
	for propName, propSchema := range schema.Properties {
		nestedPath := propName
		if parentPath != "" {
			nestedPath = parentPath + "." + propName
		}
		// Pass the required status from the parent schema's "required" array.
		// x-required/x-ves-required overrides are guarded at depth==0 so they
		// will not fire here.
		attr := ConvertToTerraformAttributeWithDepth(propName, propSchema, requiredSet[propName], "", spec, depth, nestedPath)

		// Mark 'namespace', 'tenant', 'uid', and 'kind' fields as Computed in nested Object Reference blocks.
		// The API always returns these values even when not specified in config,
		// which causes state drift if not marked as Computed.
		// Object Reference types have pattern: kind, name, namespace, tenant, uid
		// Using UseStateForUnknown plan modifier prevents perpetual drift.
		propNameLower := strings.ToLower(propName)
		if (propNameLower == "namespace" || propNameLower == "tenant" || propNameLower == "uid" || propNameLower == "kind") && !attr.Required {
			attr.Computed = true
			attr.Optional = true
			attr.PlanModifier = "UseStateForUnknown"
		}

		attrs = append(attrs, attr)
	}

	// Sort attributes: required first, then alphabetically
	sort.Slice(attrs, func(i, j int) bool {
		if attrs[i].Required != attrs[j].Required {
			return attrs[i].Required
		}
		return attrs[i].Name < attrs[j].Name
	})

	return attrs
}
