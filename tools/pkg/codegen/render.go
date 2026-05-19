// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

// Package codegen provides code rendering functions for generating Go source
// code from Terraform attribute definitions. These functions produce Go code
// strings that are embedded into generated resource files via text/template.
package codegen

import (
	"fmt"
	"strings"

	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/naming"
	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/openapi"
	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/schema"
)

// NestedModelInfo holds information needed to generate a nested model type
type NestedModelInfo struct {
	TypeName    string
	Description string
	Prefix      string // Full prefix path for generating nested type references
	Attributes  []openapi.TerraformAttribute
	IsEmpty     bool
}

// EscapeGoString escapes a string for use in a Go string literal
func EscapeGoString(s string) string {
	// Replace backslashes first, then other special characters
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

// RegexLiteral returns a Go literal representation of a regex pattern.
// Uses a raw string literal (backticks) unless the pattern contains backticks,
// in which case it falls back to a quoted string with proper escaping.
func RegexLiteral(pattern string) string {
	if !strings.Contains(pattern, "`") {
		return "`" + pattern + "`"
	}
	// Fall back to quoted string - need to escape backslashes and quotes
	escaped := strings.ReplaceAll(pattern, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `"`, `\"`)
	return `"` + escaped + `"`
}

// GetGoClientType returns the Go type for use in client structs
func GetGoClientType(attr openapi.TerraformAttribute) string {
	if attr.IsBlock {
		// Nested blocks become pointers to nested structs or slices
		if attr.NestedBlockType == "list" {
			return "[]map[string]interface{}"
		}
		return "map[string]interface{}"
	}

	switch attr.Type {
	case "string":
		return "string"
	case "int64":
		return "int64"
	case "bool":
		return "bool"
	case "list":
		if attr.ElementType == "string" {
			return "[]string"
		} else if attr.ElementType == "int64" {
			return "[]int64"
		}
		return "[]interface{}"
	case "map":
		return "map[string]string"
	default:
		return "interface{}"
	}
}

// RenderSpecStructFields generates Go struct fields for spec attributes
func RenderSpecStructFields(attrs []openapi.TerraformAttribute, indent string) string {
	specFields := schema.FilterSpecFields(attrs)
	if len(specFields) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, attr := range specFields {
		goType := GetGoClientType(attr)
		jsonTag := attr.JsonName
		if jsonTag == "" {
			jsonTag = attr.TfsdkTag
		}
		// For nested blocks, don't use omitempty - the API needs empty objects to be sent
		// For primitive fields, use omitempty to avoid sending zero values
		if attr.IsBlock {
			sb.WriteString(fmt.Sprintf("%s%s %s `json:\"%s\"`\n", indent, attr.GoName, goType, jsonTag))
		} else {
			sb.WriteString(fmt.Sprintf("%s%s %s `json:\"%s,omitempty\"`\n", indent, attr.GoName, goType, jsonTag))
		}
	}
	return sb.String()
}

// RenderSpecMarshalCodeForCreate generates Go code for Create (uses "createReq" variable)
func RenderSpecMarshalCodeForCreate(attrs []openapi.TerraformAttribute, indent string, resourceTitleCase string) string {
	return RenderSpecMarshalCodeWithVar(attrs, indent, "createReq", resourceTitleCase)
}

// RenderSpecMarshalCode generates Go code to marshal spec fields from Terraform state to API struct (uses "apiResource" variable)
func RenderSpecMarshalCode(attrs []openapi.TerraformAttribute, indent string, resourceTitleCase string) string {
	return RenderSpecMarshalCodeWithVar(attrs, indent, "apiResource", resourceTitleCase)
}

// RenderSpecMarshalCodeWithVar generates Go code to marshal spec fields with configurable variable name
func RenderSpecMarshalCodeWithVar(attrs []openapi.TerraformAttribute, indent string, varName string, resourceTitleCase string) string {
	specFields := schema.FilterSpecFields(attrs)
	if len(specFields) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, attr := range specFields {
		fieldName := attr.GoName
		tfsdkName := attr.TfsdkTag
		jsonName := attr.JsonName
		if jsonName == "" {
			jsonName = tfsdkName
		}

		if attr.IsBlock {
			// Handle list nested blocks
			if attr.NestedBlockType == "list" {
				// Check if there are any attributes that need the item variable
				// (primitives or nested blocks)
				needsItemVar := false
				for _, nestedAttr := range attr.NestedAttributes {
					switch nestedAttr.Type {
					case "string", "int64", "bool":
						needsItemVar = true
					}
					// Nested blocks also need access to item
					if nestedAttr.IsBlock {
						needsItemVar = true
					}
				}

				// Extract types.List to Go slice for iteration
				sb.WriteString(fmt.Sprintf("%sif !data.%s.IsNull() && !data.%s.IsUnknown() {\n", indent, fieldName, fieldName))
				sb.WriteString(fmt.Sprintf("%s\tvar %sItems []%s%sModel\n", indent, tfsdkName, resourceTitleCase, attr.GoName))
				sb.WriteString(fmt.Sprintf("%s\tdiags := data.%s.ElementsAs(ctx, &%sItems, false)\n", indent, fieldName, tfsdkName))
				sb.WriteString(fmt.Sprintf("%s\tresp.Diagnostics.Append(diags...)\n", indent))
				sb.WriteString(fmt.Sprintf("%s\tif !resp.Diagnostics.HasError() && len(%sItems) > 0 {\n", indent, tfsdkName))
				sb.WriteString(fmt.Sprintf("%s\t\tvar %sList []map[string]interface{}\n", indent, tfsdkName))
				if needsItemVar {
					sb.WriteString(fmt.Sprintf("%s\t\tfor _, item := range %sItems {\n", indent, tfsdkName))
				} else {
					sb.WriteString(fmt.Sprintf("%s\t\tfor range %sItems {\n", indent, tfsdkName))
				}
				sb.WriteString(fmt.Sprintf("%s\t\t\titemMap := make(map[string]interface{})\n", indent))
				// Marshal nested block fields
				for _, nestedAttr := range attr.NestedAttributes {
					nestedFieldName := nestedAttr.GoName
					nestedTfsdkName := nestedAttr.TfsdkTag
					nestedJsonName := nestedAttr.JsonName
					if nestedJsonName == "" {
						nestedJsonName = nestedTfsdkName
					}

					// Handle single nested blocks within list items
					if nestedAttr.IsBlock && nestedAttr.NestedBlockType == "single" {
						sb.WriteString(fmt.Sprintf("%s\t\tif item.%s != nil {\n", indent, nestedFieldName))
						if len(nestedAttr.NestedAttributes) == 0 {
							// Empty block - just send empty map
							sb.WriteString(fmt.Sprintf("%s\t\t\titemMap[\"%s\"] = map[string]interface{}{}\n", indent, nestedJsonName))
						} else {
							// Block with nested attributes - marshal them
							sb.WriteString(fmt.Sprintf("%s\t\t\t%sNestedMap := make(map[string]interface{})\n", indent, nestedTfsdkName))
							for _, deepAttr := range nestedAttr.NestedAttributes {
								deepFieldName := deepAttr.GoName
								deepTfsdkName := deepAttr.TfsdkTag
								deepJsonName := deepAttr.JsonName
								if deepJsonName == "" {
									deepJsonName = deepAttr.TfsdkTag
								}
								// Handle nested blocks within SingleNestedBlocks that are within ListNestedBlock items
								// e.g., allow{}, deny{}, ip_prefixes{} within action{} or match{} within rules[]
								if deepAttr.IsBlock && deepAttr.NestedBlockType == "single" {
									sb.WriteString(fmt.Sprintf("%s\t\t\tif item.%s.%s != nil {\n", indent, nestedFieldName, deepFieldName))
									if len(deepAttr.NestedAttributes) == 0 {
										// Empty nested block
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t%sNestedMap[\"%s\"] = map[string]interface{}{}\n", indent, nestedTfsdkName, deepJsonName))
									} else {
										// Nested block with attributes - create a deep map
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t%sDeepMap := make(map[string]interface{})\n", indent, deepTfsdkName))
										for _, leafAttr := range deepAttr.NestedAttributes {
											leafFieldName := leafAttr.GoName
											leafJsonName := leafAttr.JsonName
											if leafJsonName == "" {
												leafJsonName = leafAttr.TfsdkTag
											}
											// Handle empty blocks at leaf level
											if leafAttr.IsBlock && leafAttr.NestedBlockType == "single" && len(leafAttr.NestedAttributes) == 0 {
												sb.WriteString(fmt.Sprintf("%s\t\t\t\tif item.%s.%s.%s != nil {\n", indent, nestedFieldName, deepFieldName, leafFieldName))
												sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t%sDeepMap[\"%s\"] = map[string]interface{}{}\n", indent, deepTfsdkName, leafJsonName))
												sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
											} else {
												switch leafAttr.Type {
												case "string":
													sb.WriteString(fmt.Sprintf("%s\t\t\t\tif !item.%s.%s.%s.IsNull() && !item.%s.%s.%s.IsUnknown() {\n", indent, nestedFieldName, deepFieldName, leafFieldName, nestedFieldName, deepFieldName, leafFieldName))
													sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t%sDeepMap[\"%s\"] = item.%s.%s.%s.ValueString()\n", indent, deepTfsdkName, leafJsonName, nestedFieldName, deepFieldName, leafFieldName))
													sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
												case "int64":
													sb.WriteString(fmt.Sprintf("%s\t\t\t\tif !item.%s.%s.%s.IsNull() && !item.%s.%s.%s.IsUnknown() {\n", indent, nestedFieldName, deepFieldName, leafFieldName, nestedFieldName, deepFieldName, leafFieldName))
													sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t%sDeepMap[\"%s\"] = item.%s.%s.%s.ValueInt64()\n", indent, deepTfsdkName, leafJsonName, nestedFieldName, deepFieldName, leafFieldName))
													sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
												case "bool":
													sb.WriteString(fmt.Sprintf("%s\t\t\t\tif !item.%s.%s.%s.IsNull() && !item.%s.%s.%s.IsUnknown() {\n", indent, nestedFieldName, deepFieldName, leafFieldName, nestedFieldName, deepFieldName, leafFieldName))
													sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t%sDeepMap[\"%s\"] = item.%s.%s.%s.ValueBool()\n", indent, deepTfsdkName, leafJsonName, nestedFieldName, deepFieldName, leafFieldName))
													sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
												case "list":
													if leafAttr.ElementType == "string" {
														sb.WriteString(fmt.Sprintf("%s\t\t\t\tif !item.%s.%s.%s.IsNull() && !item.%s.%s.%s.IsUnknown() {\n", indent, nestedFieldName, deepFieldName, leafFieldName, nestedFieldName, deepFieldName, leafFieldName))
														sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tvar %sItems []string\n", indent, leafFieldName))
														sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tdiags := item.%s.%s.%s.ElementsAs(ctx, &%sItems, false)\n", indent, nestedFieldName, deepFieldName, leafFieldName, leafFieldName))
														sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif !diags.HasError() {\n", indent))
														sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t%sDeepMap[\"%s\"] = %sItems\n", indent, deepTfsdkName, leafJsonName, leafFieldName))
														sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
														sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
													}
												}
											}
										}
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t%sNestedMap[\"%s\"] = %sDeepMap\n", indent, nestedTfsdkName, deepJsonName, deepTfsdkName))
									}
									sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
									continue
								}
								// Handle ListNestedBlocks within SingleNestedBlocks that are within ListNestedBlock items
								// e.g., prefixes[] within ip_prefixes{} within match{} within rules[]
								if deepAttr.IsBlock && deepAttr.NestedBlockType == "list" {
									sb.WriteString(fmt.Sprintf("%s\t\t\tif len(item.%s.%s) > 0 {\n", indent, nestedFieldName, deepFieldName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\tvar %sDeepList []map[string]interface{}\n", indent, deepTfsdkName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\tfor _, deepListItem := range item.%s.%s {\n", indent, nestedFieldName, deepFieldName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tdeepListItemMap := make(map[string]interface{})\n", indent))
									for _, leafAttr := range deepAttr.NestedAttributes {
										leafFieldName := leafAttr.GoName
										leafJsonName := leafAttr.JsonName
										if leafJsonName == "" {
											leafJsonName = leafAttr.TfsdkTag
										}
										// Handle empty blocks at leaf level
										if leafAttr.IsBlock && leafAttr.NestedBlockType == "single" && len(leafAttr.NestedAttributes) == 0 {
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif deepListItem.%s != nil {\n", indent, leafFieldName))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tdeepListItemMap[\"%s\"] = map[string]interface{}{}\n", indent, leafJsonName))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
										} else {
											switch leafAttr.Type {
											case "string":
												sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif !deepListItem.%s.IsNull() && !deepListItem.%s.IsUnknown() {\n", indent, leafFieldName, leafFieldName))
												sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tdeepListItemMap[\"%s\"] = deepListItem.%s.ValueString()\n", indent, leafJsonName, leafFieldName))
												sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
											case "int64":
												sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif !deepListItem.%s.IsNull() && !deepListItem.%s.IsUnknown() {\n", indent, leafFieldName, leafFieldName))
												sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tdeepListItemMap[\"%s\"] = deepListItem.%s.ValueInt64()\n", indent, leafJsonName, leafFieldName))
												sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
											case "bool":
												sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif !deepListItem.%s.IsNull() && !deepListItem.%s.IsUnknown() {\n", indent, leafFieldName, leafFieldName))
												sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tdeepListItemMap[\"%s\"] = deepListItem.%s.ValueBool()\n", indent, leafJsonName, leafFieldName))
												sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
											}
										}
									}
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t%sDeepList = append(%sDeepList, deepListItemMap)\n", indent, deepTfsdkName, deepTfsdkName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t%sNestedMap[\"%s\"] = %sDeepList\n", indent, nestedTfsdkName, deepJsonName, deepTfsdkName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
									continue
								}
								switch deepAttr.Type {
								case "string":
									sb.WriteString(fmt.Sprintf("%s\t\t\tif !item.%s.%s.IsNull() && !item.%s.%s.IsUnknown() {\n", indent, nestedFieldName, deepFieldName, nestedFieldName, deepFieldName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t%sNestedMap[\"%s\"] = item.%s.%s.ValueString()\n", indent, nestedTfsdkName, deepJsonName, nestedFieldName, deepFieldName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
								case "int64":
									sb.WriteString(fmt.Sprintf("%s\t\t\tif !item.%s.%s.IsNull() && !item.%s.%s.IsUnknown() {\n", indent, nestedFieldName, deepFieldName, nestedFieldName, deepFieldName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t%sNestedMap[\"%s\"] = item.%s.%s.ValueInt64()\n", indent, nestedTfsdkName, deepJsonName, nestedFieldName, deepFieldName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
								case "bool":
									sb.WriteString(fmt.Sprintf("%s\t\t\tif !item.%s.%s.IsNull() && !item.%s.%s.IsUnknown() {\n", indent, nestedFieldName, deepFieldName, nestedFieldName, deepFieldName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t%sNestedMap[\"%s\"] = item.%s.%s.ValueBool()\n", indent, nestedTfsdkName, deepJsonName, nestedFieldName, deepFieldName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
								case "list":
									// Handle list types inside single nested blocks within list items
									if deepAttr.ElementType == "string" {
										sb.WriteString(fmt.Sprintf("%s\t\t\tif !item.%s.%s.IsNull() && !item.%s.%s.IsUnknown() {\n", indent, nestedFieldName, deepFieldName, nestedFieldName, deepFieldName))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\tvar %sItems []string\n", indent, deepFieldName))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\tdiags := item.%s.%s.ElementsAs(ctx, &%sItems, false)\n", indent, nestedFieldName, deepFieldName, deepFieldName))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\tif !diags.HasError() {\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t%sNestedMap[\"%s\"] = %sItems\n", indent, nestedTfsdkName, deepJsonName, deepFieldName))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
									} else if deepAttr.ElementType == "int64" {
										sb.WriteString(fmt.Sprintf("%s\t\t\tif !item.%s.%s.IsNull() && !item.%s.%s.IsUnknown() {\n", indent, nestedFieldName, deepFieldName, nestedFieldName, deepFieldName))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\tvar %sItems []int64\n", indent, deepFieldName))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\tdiags := item.%s.%s.ElementsAs(ctx, &%sItems, false)\n", indent, nestedFieldName, deepFieldName, deepFieldName))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\tif !diags.HasError() {\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t%sNestedMap[\"%s\"] = %sItems\n", indent, nestedTfsdkName, deepJsonName, deepFieldName))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
									}
								}
							}
							sb.WriteString(fmt.Sprintf("%s\t\t\titemMap[\"%s\"] = %sNestedMap\n", indent, nestedJsonName, nestedTfsdkName))
						}
						sb.WriteString(fmt.Sprintf("%s\t\t}\n", indent))
						continue
					}

					// Handle list nested blocks within list items (e.g., app_type_ref within app_type_settings)
					if nestedAttr.IsBlock && nestedAttr.NestedBlockType == "list" {
						sb.WriteString(fmt.Sprintf("%s\t\tif len(item.%s) > 0 {\n", indent, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t\tvar %sNestedList []map[string]interface{}\n", indent, nestedTfsdkName))
						sb.WriteString(fmt.Sprintf("%s\t\t\tfor _, nestedItem := range item.%s {\n", indent, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\tnestedItemMap := make(map[string]interface{})\n", indent))
						// Marshal fields from the nested list block
						for _, deepAttr := range nestedAttr.NestedAttributes {
							deepFieldName := deepAttr.GoName
							deepJsonName := deepAttr.JsonName
							if deepJsonName == "" {
								deepJsonName = deepAttr.TfsdkTag
							}
							switch deepAttr.Type {
							case "string":
								sb.WriteString(fmt.Sprintf("%s\t\t\t\tif !nestedItem.%s.IsNull() && !nestedItem.%s.IsUnknown() {\n", indent, deepFieldName, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tnestedItemMap[\"%s\"] = nestedItem.%s.ValueString()\n", indent, deepJsonName, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
							case "int64":
								sb.WriteString(fmt.Sprintf("%s\t\t\t\tif !nestedItem.%s.IsNull() && !nestedItem.%s.IsUnknown() {\n", indent, deepFieldName, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tnestedItemMap[\"%s\"] = nestedItem.%s.ValueInt64()\n", indent, deepJsonName, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
							case "bool":
								sb.WriteString(fmt.Sprintf("%s\t\t\t\tif !nestedItem.%s.IsNull() && !nestedItem.%s.IsUnknown() {\n", indent, deepFieldName, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tnestedItemMap[\"%s\"] = nestedItem.%s.ValueBool()\n", indent, deepJsonName, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
							}
						}
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t%sNestedList = append(%sNestedList, nestedItemMap)\n", indent, nestedTfsdkName, nestedTfsdkName))
						sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\titemMap[\"%s\"] = %sNestedList\n", indent, nestedJsonName, nestedTfsdkName))
						sb.WriteString(fmt.Sprintf("%s\t\t}\n", indent))
						continue
					}

					switch nestedAttr.Type {
					case "string":
						sb.WriteString(fmt.Sprintf("%s\t\tif !item.%s.IsNull() && !item.%s.IsUnknown() {\n", indent, nestedFieldName, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t\titemMap[\"%s\"] = item.%s.ValueString()\n", indent, nestedJsonName, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t}\n", indent))
					case "int64":
						sb.WriteString(fmt.Sprintf("%s\t\tif !item.%s.IsNull() && !item.%s.IsUnknown() {\n", indent, nestedFieldName, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t\titemMap[\"%s\"] = item.%s.ValueInt64()\n", indent, nestedJsonName, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t}\n", indent))
					case "bool":
						sb.WriteString(fmt.Sprintf("%s\t\tif !item.%s.IsNull() && !item.%s.IsUnknown() {\n", indent, nestedFieldName, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t\titemMap[\"%s\"] = item.%s.ValueBool()\n", indent, nestedJsonName, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t}\n", indent))
					}
				}
				sb.WriteString(fmt.Sprintf("%s\t\t\t%sList = append(%sList, itemMap)\n", indent, tfsdkName, tfsdkName))
				sb.WriteString(fmt.Sprintf("%s\t\t}\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t\t%s.Spec[\"%s\"] = %sList\n", indent, varName, jsonName, tfsdkName))
				sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
				sb.WriteString(fmt.Sprintf("%s}\n", indent))
				continue
			}
			// Handle single nested blocks by converting model struct to map
			sb.WriteString(fmt.Sprintf("%sif data.%s != nil {\n", indent, fieldName))
			sb.WriteString(fmt.Sprintf("%s\t%sMap := make(map[string]interface{})\n", indent, tfsdkName))
			// Marshal nested block fields
			for _, nestedAttr := range attr.NestedAttributes {
				nestedFieldName := nestedAttr.GoName
				nestedTfsdkName := nestedAttr.TfsdkTag
				nestedJsonName := nestedAttr.JsonName
				if nestedJsonName == "" {
					nestedJsonName = nestedTfsdkName
				}

				// Handle single nested blocks within single nested blocks (not list nested blocks)
				if nestedAttr.IsBlock && nestedAttr.NestedBlockType == "single" {
					if len(nestedAttr.NestedAttributes) == 0 {
						// Empty nested block - send empty map if present
						sb.WriteString(fmt.Sprintf("%s\tif data.%s.%s != nil {\n", indent, fieldName, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t%sMap[\"%s\"] = map[string]interface{}{}\n", indent, tfsdkName, nestedJsonName))
						sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
					} else {
						// Nested block with attributes - marshal them
						sb.WriteString(fmt.Sprintf("%s\tif data.%s.%s != nil {\n", indent, fieldName, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t%sNestedMap := make(map[string]interface{})\n", indent, nestedTfsdkName))
						for _, deepAttr := range nestedAttr.NestedAttributes {
							deepFieldName := deepAttr.GoName
							deepJsonName := deepAttr.JsonName
							if deepJsonName == "" {
								deepJsonName = deepAttr.TfsdkTag
							}
							switch deepAttr.Type {
							case "string":
								sb.WriteString(fmt.Sprintf("%s\t\tif !data.%s.%s.%s.IsNull() && !data.%s.%s.%s.IsUnknown() {\n", indent, fieldName, nestedFieldName, deepFieldName, fieldName, nestedFieldName, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t%sNestedMap[\"%s\"] = data.%s.%s.%s.ValueString()\n", indent, nestedTfsdkName, deepJsonName, fieldName, nestedFieldName, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t}\n", indent))
							case "int64":
								sb.WriteString(fmt.Sprintf("%s\t\tif !data.%s.%s.%s.IsNull() && !data.%s.%s.%s.IsUnknown() {\n", indent, fieldName, nestedFieldName, deepFieldName, fieldName, nestedFieldName, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t%sNestedMap[\"%s\"] = data.%s.%s.%s.ValueInt64()\n", indent, nestedTfsdkName, deepJsonName, fieldName, nestedFieldName, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t}\n", indent))
							case "bool":
								sb.WriteString(fmt.Sprintf("%s\t\tif !data.%s.%s.%s.IsNull() && !data.%s.%s.%s.IsUnknown() {\n", indent, fieldName, nestedFieldName, deepFieldName, fieldName, nestedFieldName, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t%sNestedMap[\"%s\"] = data.%s.%s.%s.ValueBool()\n", indent, nestedTfsdkName, deepJsonName, fieldName, nestedFieldName, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t}\n", indent))
							}
						}
						sb.WriteString(fmt.Sprintf("%s\t\t%sMap[\"%s\"] = %sNestedMap\n", indent, tfsdkName, nestedJsonName, nestedTfsdkName))
						sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
					}
					continue
				}

				// Handle list nested blocks within single nested blocks (e.g., rules within mitigation_type)
				if nestedAttr.IsBlock && nestedAttr.NestedBlockType == "list" {
					sb.WriteString(fmt.Sprintf("%s\tif len(data.%s.%s) > 0 {\n", indent, fieldName, nestedFieldName))
					sb.WriteString(fmt.Sprintf("%s\t\tvar %sList []map[string]interface{}\n", indent, nestedTfsdkName))
					sb.WriteString(fmt.Sprintf("%s\t\tfor _, listItem := range data.%s.%s {\n", indent, fieldName, nestedFieldName))
					sb.WriteString(fmt.Sprintf("%s\t\t\tlistItemMap := make(map[string]interface{})\n", indent))

					// Marshal deep nested attributes within the list items
					for _, deepAttr := range nestedAttr.NestedAttributes {
						deepFieldName := deepAttr.GoName
						deepTfsdkName := deepAttr.TfsdkTag
						deepJsonName := deepAttr.JsonName
						if deepJsonName == "" {
							deepJsonName = deepTfsdkName
						}

						// Handle single nested blocks within list items
						if deepAttr.IsBlock && deepAttr.NestedBlockType == "single" {
							sb.WriteString(fmt.Sprintf("%s\t\t\tif listItem.%s != nil {\n", indent, deepFieldName))
							if len(deepAttr.NestedAttributes) == 0 {
								// Empty block - just send empty map
								sb.WriteString(fmt.Sprintf("%s\t\t\t\tlistItemMap[\"%s\"] = map[string]interface{}{}\n", indent, deepJsonName))
							} else {
								// Block with nested attributes (like choice blocks)
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t%sDeepMap := make(map[string]interface{})\n", indent, deepTfsdkName))
								for _, leafAttr := range deepAttr.NestedAttributes {
									leafFieldName := leafAttr.GoName
									leafTfsdkName := leafAttr.TfsdkTag
									leafJsonName := leafAttr.JsonName
									if leafJsonName == "" {
										leafJsonName = leafTfsdkName
									}
									// Handle empty choice blocks
									if leafAttr.IsBlock && leafAttr.NestedBlockType == "single" && len(leafAttr.NestedAttributes) == 0 {
										sb.WriteString(fmt.Sprintf("%s\t\t\t\tif listItem.%s.%s != nil {\n", indent, deepFieldName, leafFieldName))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t%sDeepMap[\"%s\"] = map[string]interface{}{}\n", indent, deepTfsdkName, leafJsonName))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
									} else {
										// Handle primitive types at leaf level
										switch leafAttr.Type {
										case "string":
											sb.WriteString(fmt.Sprintf("%s\t\t\t\tif !listItem.%s.%s.IsNull() && !listItem.%s.%s.IsUnknown() {\n", indent, deepFieldName, leafFieldName, deepFieldName, leafFieldName))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t%sDeepMap[\"%s\"] = listItem.%s.%s.ValueString()\n", indent, deepTfsdkName, leafJsonName, deepFieldName, leafFieldName))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
										case "int64":
											sb.WriteString(fmt.Sprintf("%s\t\t\t\tif !listItem.%s.%s.IsNull() && !listItem.%s.%s.IsUnknown() {\n", indent, deepFieldName, leafFieldName, deepFieldName, leafFieldName))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t%sDeepMap[\"%s\"] = listItem.%s.%s.ValueInt64()\n", indent, deepTfsdkName, leafJsonName, deepFieldName, leafFieldName))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
										case "bool":
											sb.WriteString(fmt.Sprintf("%s\t\t\t\tif !listItem.%s.%s.IsNull() && !listItem.%s.%s.IsUnknown() {\n", indent, deepFieldName, leafFieldName, deepFieldName, leafFieldName))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t%sDeepMap[\"%s\"] = listItem.%s.%s.ValueBool()\n", indent, deepTfsdkName, leafJsonName, deepFieldName, leafFieldName))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
										}
									}
								}
								sb.WriteString(fmt.Sprintf("%s\t\t\t\tlistItemMap[\"%s\"] = %sDeepMap\n", indent, deepJsonName, deepTfsdkName))
							}
							sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
							continue
						}

						// Handle primitive types within list items
						switch deepAttr.Type {
						case "string":
							sb.WriteString(fmt.Sprintf("%s\t\t\tif !listItem.%s.IsNull() && !listItem.%s.IsUnknown() {\n", indent, deepFieldName, deepFieldName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\tlistItemMap[\"%s\"] = listItem.%s.ValueString()\n", indent, deepJsonName, deepFieldName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
						case "int64":
							sb.WriteString(fmt.Sprintf("%s\t\t\tif !listItem.%s.IsNull() && !listItem.%s.IsUnknown() {\n", indent, deepFieldName, deepFieldName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\tlistItemMap[\"%s\"] = listItem.%s.ValueInt64()\n", indent, deepJsonName, deepFieldName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
						case "bool":
							sb.WriteString(fmt.Sprintf("%s\t\t\tif !listItem.%s.IsNull() && !listItem.%s.IsUnknown() {\n", indent, deepFieldName, deepFieldName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\tlistItemMap[\"%s\"] = listItem.%s.ValueBool()\n", indent, deepJsonName, deepFieldName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
						}
					}

					sb.WriteString(fmt.Sprintf("%s\t\t\t%sList = append(%sList, listItemMap)\n", indent, nestedTfsdkName, nestedTfsdkName))
					sb.WriteString(fmt.Sprintf("%s\t\t}\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t%sMap[\"%s\"] = %sList\n", indent, tfsdkName, nestedJsonName, nestedTfsdkName))
					sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
					continue
				}

				switch nestedAttr.Type {
				case "string":
					sb.WriteString(fmt.Sprintf("%s\tif !data.%s.%s.IsNull() && !data.%s.%s.IsUnknown() {\n", indent, fieldName, nestedFieldName, fieldName, nestedFieldName))
					sb.WriteString(fmt.Sprintf("%s\t\t%sMap[\"%s\"] = data.%s.%s.ValueString()\n", indent, tfsdkName, nestedTfsdkName, fieldName, nestedFieldName))
					sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
				case "int64":
					sb.WriteString(fmt.Sprintf("%s\tif !data.%s.%s.IsNull() && !data.%s.%s.IsUnknown() {\n", indent, fieldName, nestedFieldName, fieldName, nestedFieldName))
					sb.WriteString(fmt.Sprintf("%s\t\t%sMap[\"%s\"] = data.%s.%s.ValueInt64()\n", indent, tfsdkName, nestedTfsdkName, fieldName, nestedFieldName))
					sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
				case "bool":
					sb.WriteString(fmt.Sprintf("%s\tif !data.%s.%s.IsNull() && !data.%s.%s.IsUnknown() {\n", indent, fieldName, nestedFieldName, fieldName, nestedFieldName))
					sb.WriteString(fmt.Sprintf("%s\t\t%sMap[\"%s\"] = data.%s.%s.ValueBool()\n", indent, tfsdkName, nestedTfsdkName, fieldName, nestedFieldName))
					sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
				case "list":
					// Handle list types (e.g., expressions: types.List) inside single nested blocks
					if nestedAttr.ElementType == "string" {
						sb.WriteString(fmt.Sprintf("%s\tif !data.%s.%s.IsNull() && !data.%s.%s.IsUnknown() {\n", indent, fieldName, nestedFieldName, fieldName, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\tvar %sItems []string\n", indent, nestedTfsdkName))
						sb.WriteString(fmt.Sprintf("%s\t\tdiags := data.%s.%s.ElementsAs(ctx, &%sItems, false)\n", indent, fieldName, nestedFieldName, nestedTfsdkName))
						sb.WriteString(fmt.Sprintf("%s\t\tif !diags.HasError() {\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t%sMap[\"%s\"] = %sItems\n", indent, tfsdkName, nestedTfsdkName, nestedTfsdkName))
						sb.WriteString(fmt.Sprintf("%s\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
					} else if nestedAttr.ElementType == "int64" {
						sb.WriteString(fmt.Sprintf("%s\tif !data.%s.%s.IsNull() && !data.%s.%s.IsUnknown() {\n", indent, fieldName, nestedFieldName, fieldName, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\tvar %sItems []int64\n", indent, nestedTfsdkName))
						sb.WriteString(fmt.Sprintf("%s\t\tdiags := data.%s.%s.ElementsAs(ctx, &%sItems, false)\n", indent, fieldName, nestedFieldName, nestedTfsdkName))
						sb.WriteString(fmt.Sprintf("%s\t\tif !diags.HasError() {\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t%sMap[\"%s\"] = %sItems\n", indent, tfsdkName, nestedTfsdkName, nestedTfsdkName))
						sb.WriteString(fmt.Sprintf("%s\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
					}
				}
			}
			sb.WriteString(fmt.Sprintf("%s\t%s.Spec[\"%s\"] = %sMap\n", indent, varName, jsonName, tfsdkName))
			sb.WriteString(fmt.Sprintf("%s}\n", indent))
			continue
		}

		switch attr.Type {
		case "string":
			sb.WriteString(fmt.Sprintf("%sif !data.%s.IsNull() && !data.%s.IsUnknown() {\n", indent, fieldName, fieldName))
			sb.WriteString(fmt.Sprintf("%s\t%s.Spec[\"%s\"] = data.%s.ValueString()\n", indent, varName, jsonName, fieldName))
			sb.WriteString(fmt.Sprintf("%s}\n", indent))
		case "int64":
			sb.WriteString(fmt.Sprintf("%sif !data.%s.IsNull() && !data.%s.IsUnknown() {\n", indent, fieldName, fieldName))
			sb.WriteString(fmt.Sprintf("%s\t%s.Spec[\"%s\"] = data.%s.ValueInt64()\n", indent, varName, jsonName, fieldName))
			sb.WriteString(fmt.Sprintf("%s}\n", indent))
		case "bool":
			sb.WriteString(fmt.Sprintf("%sif !data.%s.IsNull() && !data.%s.IsUnknown() {\n", indent, fieldName, fieldName))
			sb.WriteString(fmt.Sprintf("%s\t%s.Spec[\"%s\"] = data.%s.ValueBool()\n", indent, varName, jsonName, fieldName))
			sb.WriteString(fmt.Sprintf("%s}\n", indent))
		case "list":
			if attr.ElementType == "string" {
				sb.WriteString(fmt.Sprintf("%sif !data.%s.IsNull() && !data.%s.IsUnknown() {\n", indent, fieldName, fieldName))
				sb.WriteString(fmt.Sprintf("%s\tvar %sList []string\n", indent, tfsdkName))
				sb.WriteString(fmt.Sprintf("%s\tresp.Diagnostics.Append(data.%s.ElementsAs(ctx, &%sList, false)...)\n", indent, fieldName, tfsdkName))
				sb.WriteString(fmt.Sprintf("%s\tif !resp.Diagnostics.HasError() {\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t\t%s.Spec[\"%s\"] = %sList\n", indent, varName, jsonName, tfsdkName))
				sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
				sb.WriteString(fmt.Sprintf("%s}\n", indent))
			} else if attr.ElementType == "int64" {
				sb.WriteString(fmt.Sprintf("%sif !data.%s.IsNull() && !data.%s.IsUnknown() {\n", indent, fieldName, fieldName))
				sb.WriteString(fmt.Sprintf("%s\tvar %sList []int64\n", indent, tfsdkName))
				sb.WriteString(fmt.Sprintf("%s\tresp.Diagnostics.Append(data.%s.ElementsAs(ctx, &%sList, false)...)\n", indent, fieldName, tfsdkName))
				sb.WriteString(fmt.Sprintf("%s\tif !resp.Diagnostics.HasError() {\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t\t%s.Spec[\"%s\"] = %sList\n", indent, varName, jsonName, tfsdkName))
				sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
				sb.WriteString(fmt.Sprintf("%s}\n", indent))
			}
		}
	}
	return sb.String()
}

// RenderComputedFieldsCode generates Go code to set Computed+Optional fields from API response
// This ensures that fields with UseStateForUnknown plan modifier have known values after Create/Update
// The varName parameter specifies the API response variable name (e.g., "created" or "updated")
func RenderComputedFieldsCode(attrs []openapi.TerraformAttribute, indent string, varName string) string {
	specFields := schema.FilterSpecFields(attrs)
	if len(specFields) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(indent + "// Set computed fields from API response\n")

	for _, attr := range specFields {
		// Only generate code for Computed+Optional scalar fields (UseStateForUnknown pattern)
		if !attr.Computed || !attr.Optional || attr.IsBlock {
			continue
		}

		fieldName := attr.GoName
		jsonName := attr.JsonName
		if jsonName == "" {
			jsonName = attr.TfsdkTag
		}

		switch attr.Type {
		case "bool":
			// Set value if API returns it; otherwise handle based on plan value:
			// - If plan was unknown, set to null (resolves unknown state after apply)
			// - If plan had a value, preserve it (user specified this value)
			sb.WriteString(fmt.Sprintf("%sif v, ok := %s.Spec[\"%s\"].(bool); ok {\n", indent, varName, jsonName))
			sb.WriteString(fmt.Sprintf("%s\tdata.%s = types.BoolValue(v)\n", indent, fieldName))
			sb.WriteString(fmt.Sprintf("%s} else if data.%s.IsUnknown() {\n", indent, fieldName))
			sb.WriteString(fmt.Sprintf("%s\t// API didn't return value and plan was unknown - set to null\n", indent))
			sb.WriteString(fmt.Sprintf("%s\tdata.%s = types.BoolNull()\n", indent, fieldName))
			sb.WriteString(fmt.Sprintf("%s}\n", indent))
			sb.WriteString(fmt.Sprintf("%s// If plan had a value, preserve it\n", indent))
		case "int64":
			// Set value if API returns it; otherwise handle based on plan value
			sb.WriteString(fmt.Sprintf("%sif v, ok := %s.Spec[\"%s\"].(float64); ok {\n", indent, varName, jsonName))
			sb.WriteString(fmt.Sprintf("%s\tdata.%s = types.Int64Value(int64(v))\n", indent, fieldName))
			sb.WriteString(fmt.Sprintf("%s} else if data.%s.IsUnknown() {\n", indent, fieldName))
			sb.WriteString(fmt.Sprintf("%s\t// API didn't return value and plan was unknown - set to null\n", indent))
			sb.WriteString(fmt.Sprintf("%s\tdata.%s = types.Int64Null()\n", indent, fieldName))
			sb.WriteString(fmt.Sprintf("%s}\n", indent))
			sb.WriteString(fmt.Sprintf("%s// If plan had a value, preserve it\n", indent))
		case "string":
			// Set value if API returns it; otherwise handle based on plan value
			sb.WriteString(fmt.Sprintf("%sif v, ok := %s.Spec[\"%s\"].(string); ok && v != \"\" {\n", indent, varName, jsonName))
			sb.WriteString(fmt.Sprintf("%s\tdata.%s = types.StringValue(v)\n", indent, fieldName))
			sb.WriteString(fmt.Sprintf("%s} else if data.%s.IsUnknown() {\n", indent, fieldName))
			sb.WriteString(fmt.Sprintf("%s\t// API didn't return value and plan was unknown - set to null\n", indent))
			sb.WriteString(fmt.Sprintf("%s\tdata.%s = types.StringNull()\n", indent, fieldName))
			sb.WriteString(fmt.Sprintf("%s}\n", indent))
			sb.WriteString(fmt.Sprintf("%s// If plan had a value, preserve it\n", indent))
		}
	}

	return sb.String()
}

// RenderCreateComputedFieldsCode generates code for Create method (uses "created" variable)
func RenderCreateComputedFieldsCode(attrs []openapi.TerraformAttribute, indent string) string {
	return RenderComputedFieldsCode(attrs, indent, "created")
}

// RenderUpdateComputedFieldsCode generates code for Update method (uses "updated" variable)
// Deprecated: Use RenderFetchedComputedFieldsCode instead since Update now uses GET after PUT
func RenderUpdateComputedFieldsCode(attrs []openapi.TerraformAttribute, indent string) string {
	return RenderComputedFieldsCode(attrs, indent, "updated")
}

// RenderFetchedComputedFieldsCode generates code for Update method after GET (uses "fetched" variable)
func RenderFetchedComputedFieldsCode(attrs []openapi.TerraformAttribute, indent string) string {
	return RenderComputedFieldsCode(attrs, indent, "fetched")
}

// RenderSpecUnmarshalCode generates Go code to unmarshal spec fields from API response to Terraform state
func RenderSpecUnmarshalCode(attrs []openapi.TerraformAttribute, indent string, resourceTitleCase string) string {
	specFields := schema.FilterSpecFields(attrs)
	if len(specFields) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, attr := range specFields {
		fieldName := attr.GoName
		jsonName := attr.JsonName
		if jsonName == "" {
			jsonName = attr.TfsdkTag
		}

		if attr.IsBlock {
			// Handle list nested blocks
			if attr.NestedBlockType == "list" {
				tfsdkName := attr.TfsdkTag

				// Check if there are any attributes that need the itemMap variable
				// (primitives, lists, or nested blocks with attributes)
				needsItemMap := false
				for _, nestedAttr := range attr.NestedAttributes {
					switch nestedAttr.Type {
					case "string", "int64", "bool", "list":
						needsItemMap = true
					}
					// Nested blocks with attributes need access to itemMap
					// Empty marker blocks (no nested attributes) don't need itemMap since
					// they're skipped to prevent state drift
					if nestedAttr.IsBlock && len(nestedAttr.NestedAttributes) > 0 {
						needsItemMap = true
					}
				}

				sb.WriteString(fmt.Sprintf("%sif listData, ok := apiResource.Spec[\"%s\"].([]interface{}); ok && len(listData) > 0 {\n", indent, jsonName))
				sb.WriteString(fmt.Sprintf("%s\tvar %sList []%s%sModel\n", indent, tfsdkName, resourceTitleCase, attr.GoName))
				// Extract current state's types.List to Go slice for preserving empty blocks
				sb.WriteString(fmt.Sprintf("%s\tvar existing%sItems []%s%sModel\n", indent, attr.GoName, resourceTitleCase, attr.GoName))
				sb.WriteString(fmt.Sprintf("%s\tif !data.%s.IsNull() && !data.%s.IsUnknown() {\n", indent, attr.GoName, attr.GoName))
				sb.WriteString(fmt.Sprintf("%s\t\tdata.%s.ElementsAs(ctx, &existing%sItems, false)\n", indent, attr.GoName, attr.GoName))
				sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
				if needsItemMap {
					sb.WriteString(fmt.Sprintf("%s\tfor listIdx, item := range listData {\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t_ = listIdx // May be unused if no empty marker blocks in list item\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\tif itemMap, ok := item.(map[string]interface{}); ok {\n", indent))
				} else {
					sb.WriteString(fmt.Sprintf("%s\tfor listIdx := range listData {\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t_ = listIdx // May be unused if no empty marker blocks in list item\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t{\n", indent))
				}
				sb.WriteString(fmt.Sprintf("%s\t\t\t%sList = append(%sList, %s%sModel{\n", indent, tfsdkName, tfsdkName, resourceTitleCase, attr.GoName))
				// Unmarshal nested block fields
				for _, nestedAttr := range attr.NestedAttributes {
					nestedFieldName := nestedAttr.GoName
					nestedJsonName := nestedAttr.JsonName
					if nestedJsonName == "" {
						nestedJsonName = nestedAttr.TfsdkTag
					}

					// Handle single nested blocks within list items
					if nestedAttr.IsBlock && nestedAttr.NestedBlockType == "single" {
						if len(nestedAttr.NestedAttributes) == 0 {
							// Empty block (marker block) - Preserve from existing state if user configured it
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t%s: func() *%sEmptyModel {\n", indent, nestedFieldName, resourceTitleCase))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif !isImport && len(existing%sItems) > listIdx && existing%sItems[listIdx].%s != nil {\n", indent, attr.GoName, attr.GoName, nestedFieldName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\treturn &%sEmptyModel{}\n", indent, resourceTitleCase))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\treturn nil\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t}(),\n", indent))
						} else {
							// Block with nested attributes - unmarshal them
							nestedModelType := resourceTitleCase + attr.GoName + nestedAttr.GoName + "Model"

							// Check if nested block has any primitive or list attributes that we can unmarshal
							hasDeepPrimitives := false
							for _, deepAttr := range nestedAttr.NestedAttributes {
								switch deepAttr.Type {
								case "string", "int64", "bool":
									hasDeepPrimitives = true
								case "list":
									if deepAttr.ElementType == "string" || deepAttr.ElementType == "int64" {
										hasDeepPrimitives = true
									}
								}
							}

							sb.WriteString(fmt.Sprintf("%s\t\t\t\t%s: func() *%s {\n", indent, nestedFieldName, nestedModelType))
							if hasDeepPrimitives {
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif nestedMap, ok := itemMap[\"%s\"].(map[string]interface{}); ok {\n", indent, nestedJsonName))
							} else {
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif _, ok := itemMap[\"%s\"].(map[string]interface{}); ok {\n", indent, nestedJsonName))
							}
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\treturn &%s{\n", indent, nestedModelType))
							for _, deepAttr := range nestedAttr.NestedAttributes {
								deepFieldName := deepAttr.GoName
								deepJsonName := deepAttr.JsonName
								if deepJsonName == "" {
									deepJsonName = deepAttr.TfsdkTag
								}
								// Handle deeply nested empty blocks
								if deepAttr.IsBlock && len(deepAttr.NestedAttributes) == 0 {
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t%s: func() *%sEmptyModel {\n", indent, deepFieldName, resourceTitleCase))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\tif !isImport && len(existing%sItems) > listIdx && existing%sItems[listIdx].%s != nil && existing%sItems[listIdx].%s.%s != nil {\n", indent, attr.GoName, attr.GoName, nestedAttr.GoName, attr.GoName, nestedAttr.GoName, deepFieldName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\treturn &%sEmptyModel{}\n", indent, resourceTitleCase))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t}\n", indent))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\treturn nil\n", indent))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t}(),\n", indent))
									continue
								}
								switch deepAttr.Type {
								case "string":
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t%s: func() types.String {\n", indent, deepFieldName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\tif v, ok := nestedMap[\"%s\"].(string); ok && v != \"\" {\n", indent, deepJsonName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\treturn types.StringValue(v)\n", indent))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t}\n", indent))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\treturn types.StringNull()\n", indent))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t}(),\n", indent))
								case "int64":
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t%s: func() types.Int64 {\n", indent, deepFieldName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\tif v, ok := nestedMap[\"%s\"].(float64); ok && v != 0 {\n", indent, deepJsonName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\treturn types.Int64Value(int64(v))\n", indent))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t}\n", indent))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\treturn types.Int64Null()\n", indent))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t}(),\n", indent))
								case "bool":
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t%s: func() types.Bool {\n", indent, deepFieldName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\tif v, ok := nestedMap[\"%s\"].(bool); ok {\n", indent, deepJsonName))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\treturn types.BoolValue(v)\n", indent))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t}\n", indent))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\treturn types.BoolNull()\n", indent))
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t}(),\n", indent))
								case "list":
									// Handle list types inside single nested blocks within list items
									if deepAttr.ElementType == "string" {
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t%s: func() types.List {\n", indent, deepFieldName))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\tif v, ok := nestedMap[\"%s\"].([]interface{}); ok && len(v) > 0 {\n", indent, deepJsonName))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\tvar items []string\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\tfor _, item := range v {\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\tif s, ok := item.(string); ok {\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\titems = append(items, s)\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t}\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t}\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\tlistVal, _ := types.ListValueFrom(ctx, types.StringType, items)\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\treturn listVal\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t}\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\treturn types.ListNull(types.StringType)\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t}(),\n", indent))
									} else if deepAttr.ElementType == "int64" {
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t%s: func() types.List {\n", indent, deepFieldName))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\tif v, ok := nestedMap[\"%s\"].([]interface{}); ok && len(v) > 0 {\n", indent, deepJsonName))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\tvar items []int64\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\tfor _, item := range v {\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\tif n, ok := item.(float64); ok {\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\titems = append(items, int64(n))\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t}\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t}\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\tlistVal, _ := types.ListValueFrom(ctx, types.Int64Type, items)\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\treturn listVal\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t}\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\treturn types.ListNull(types.Int64Type)\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t}(),\n", indent))
									}
								}
							}
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t}\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\treturn nil\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t}(),\n", indent))
						}
						continue
					}

					// Handle list nested blocks within list items (e.g., app_type_ref within app_type_settings)
					if nestedAttr.IsBlock && nestedAttr.NestedBlockType == "list" {
						nestedListModelType := resourceTitleCase + attr.GoName + nestedAttr.GoName + "Model"

						sb.WriteString(fmt.Sprintf("%s\t\t\t\t%s: func() []%s {\n", indent, nestedFieldName, nestedListModelType))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif nestedListData, ok := itemMap[\"%s\"].([]interface{}); ok && len(nestedListData) > 0 {\n", indent, nestedJsonName))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tvar result []%s\n", indent, nestedListModelType))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tfor _, nestedItem := range nestedListData {\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\tif nestedItemMap, ok := nestedItem.(map[string]interface{}); ok {\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\tresult = append(result, %s{\n", indent, nestedListModelType))
						// Generate field reading for each nested attribute
						for _, deepAttr := range nestedAttr.NestedAttributes {
							deepFieldName := deepAttr.GoName
							deepJsonName := deepAttr.JsonName
							if deepJsonName == "" {
								deepJsonName = deepAttr.TfsdkTag
							}
							switch deepAttr.Type {
							case "string":
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t%s: func() types.String {\n", indent, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\tif v, ok := nestedItemMap[\"%s\"].(string); ok && v != \"\" {\n", indent, deepJsonName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\treturn types.StringValue(v)\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t}\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\treturn types.StringNull()\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t}(),\n", indent))
							case "int64":
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t%s: func() types.Int64 {\n", indent, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\tif v, ok := nestedItemMap[\"%s\"].(float64); ok && v != 0 {\n", indent, deepJsonName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\treturn types.Int64Value(int64(v))\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t}\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\treturn types.Int64Null()\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t}(),\n", indent))
							case "bool":
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t%s: func() types.Bool {\n", indent, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\tif v, ok := nestedItemMap[\"%s\"].(bool); ok {\n", indent, deepJsonName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\treturn types.BoolValue(v)\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t}\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\treturn types.BoolNull()\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t}(),\n", indent))
							}
						}
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t})\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\treturn result\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\treturn nil\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t}(),\n", indent))
						continue
					}

					switch nestedAttr.Type {
					case "string":
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t%s: func() types.String {\n", indent, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif v, ok := itemMap[\"%s\"].(string); ok && v != \"\" {\n", indent, nestedJsonName))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\treturn types.StringValue(v)\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\treturn types.StringNull()\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t}(),\n", indent))
					case "int64":
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t%s: func() types.Int64 {\n", indent, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif v, ok := itemMap[\"%s\"].(float64); ok && v != 0 {\n", indent, nestedJsonName))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\treturn types.Int64Value(int64(v))\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\treturn types.Int64Null()\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t}(),\n", indent))
					case "bool":
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t%s: func() types.Bool {\n", indent, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif v, ok := itemMap[\"%s\"].(bool); ok {\n", indent, nestedJsonName))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\treturn types.BoolValue(v)\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\treturn types.BoolNull()\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t}(),\n", indent))
					case "list":
						// Handle list types inside list nested blocks
						if nestedAttr.ElementType == "string" {
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t%s: func() types.List {\n", indent, nestedFieldName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif v, ok := itemMap[\"%s\"].([]interface{}); ok && len(v) > 0 {\n", indent, nestedJsonName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tvar items []string\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tfor _, item := range v {\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\tif s, ok := item.(string); ok {\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\titems = append(items, s)\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t}\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t}\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tlistVal, _ := types.ListValueFrom(ctx, types.StringType, items)\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\treturn listVal\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\treturn types.ListNull(types.StringType)\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t}(),\n", indent))
						} else if nestedAttr.ElementType == "int64" {
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t%s: func() types.List {\n", indent, nestedFieldName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif v, ok := itemMap[\"%s\"].([]interface{}); ok && len(v) > 0 {\n", indent, nestedJsonName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tvar items []int64\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tfor _, item := range v {\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\tif n, ok := item.(float64); ok {\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\titems = append(items, int64(n))\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t}\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t}\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tlistVal, _ := types.ListValueFrom(ctx, types.Int64Type, items)\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\treturn listVal\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\treturn types.ListNull(types.Int64Type)\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t}(),\n", indent))
						}
					}
				}
				sb.WriteString(fmt.Sprintf("%s\t\t\t})\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t\t}\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
				// Convert Go slice to types.List for proper unknown value handling
				sb.WriteString(fmt.Sprintf("%s\tlistVal, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: %s%sModelAttrTypes}, %sList)\n", indent, resourceTitleCase, attr.GoName, tfsdkName))
				sb.WriteString(fmt.Sprintf("%s\tresp.Diagnostics.Append(diags...)\n", indent))
				sb.WriteString(fmt.Sprintf("%s\tif !resp.Diagnostics.HasError() {\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t\tdata.%s = listVal\n", indent, fieldName))
				sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
				sb.WriteString(fmt.Sprintf("%s} else {\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t// No data from API - set to null list\n", indent))
				sb.WriteString(fmt.Sprintf("%s\tdata.%s = types.ListNull(types.ObjectType{AttrTypes: %s%sModelAttrTypes})\n", indent, fieldName, resourceTitleCase, attr.GoName))
				sb.WriteString(fmt.Sprintf("%s}\n", indent))
				continue
			}

			// Empty blocks (no nested attributes) use EmptyModel
			if len(attr.NestedAttributes) == 0 {
				sb.WriteString(fmt.Sprintf("%sif _, ok := apiResource.Spec[\"%s\"].(map[string]interface{}); ok && isImport && data.%s == nil {\n", indent, jsonName, fieldName))
				sb.WriteString(fmt.Sprintf("%s\t// Import case: populate from API since state is nil and psd is empty\n", indent))
				sb.WriteString(fmt.Sprintf("%s\tdata.%s = &%sEmptyModel{}\n", indent, fieldName, resourceTitleCase))
				sb.WriteString(fmt.Sprintf("%s}\n", indent))
				sb.WriteString(fmt.Sprintf("%s// Normal Read: preserve existing state value\n", indent))
				continue
			}

			// Check if we need to read data (for primitives, lists, or nested list blocks)
			hasReadableData := false
			for _, nestedAttr := range attr.NestedAttributes {
				switch nestedAttr.Type {
				case "string", "int64", "bool", "list":
					hasReadableData = true
				}
				if nestedAttr.IsBlock && nestedAttr.NestedBlockType == "list" {
					hasReadableData = true
				}
			}

			// If the block contains ONLY nested empty blocks (no primitives, lists, or nested list blocks), handle specially
			if !hasReadableData {
				sb.WriteString(fmt.Sprintf("%sif _, ok := apiResource.Spec[\"%s\"].(map[string]interface{}); ok && isImport && data.%s == nil {\n", indent, jsonName, fieldName))
				sb.WriteString(fmt.Sprintf("%s\t// Import case: populate from API since state is nil and psd is empty\n", indent))
				sb.WriteString(fmt.Sprintf("%s\tdata.%s = &%s%sModel{}\n", indent, fieldName, resourceTitleCase, attr.GoName))
				sb.WriteString(fmt.Sprintf("%s}\n", indent))
				sb.WriteString(fmt.Sprintf("%s// Normal Read: preserve existing state value\n", indent))
				continue
			}
			// Unmarshal single nested block from API response map
			sb.WriteString(fmt.Sprintf("%sif blockData, ok := apiResource.Spec[\"%s\"].(map[string]interface{}); ok && (isImport || data.%s != nil) {\n", indent, jsonName, fieldName))
			sb.WriteString(fmt.Sprintf("%s\tdata.%s = &%s%sModel{\n", indent, fieldName, resourceTitleCase, attr.GoName))
			// Unmarshal nested block fields
			for _, nestedAttr := range attr.NestedAttributes {
				nestedFieldName := nestedAttr.GoName
				nestedJsonName := nestedAttr.JsonName
				if nestedJsonName == "" {
					nestedJsonName = nestedAttr.TfsdkTag
				}

				// Handle empty single nested blocks (presence markers like default_action_deny {})
				if nestedAttr.IsBlock && nestedAttr.NestedBlockType == "single" && len(nestedAttr.NestedAttributes) == 0 {
					sb.WriteString(fmt.Sprintf("%s\t\t%s: func() *%sEmptyModel {\n", indent, nestedFieldName, resourceTitleCase))
					sb.WriteString(fmt.Sprintf("%s\t\t\tif !isImport && data.%s != nil {\n", indent, fieldName))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\t// Normal Read: preserve existing state value (even if nil)\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\t// This prevents API returning empty objects from overwriting user's 'not configured' intent\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\treturn data.%s.%s\n", indent, fieldName, nestedFieldName))
					sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\t// Import case: read from API\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\tif _, ok := blockData[\"%s\"].(map[string]interface{}); ok {\n", indent, nestedJsonName))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\treturn &%sEmptyModel{}\n", indent, resourceTitleCase))
					sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\treturn nil\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t}(),\n", indent))
					continue
				}

				switch nestedAttr.Type {
				case "string":
					sb.WriteString(fmt.Sprintf("%s\t\t%s: func() types.String {\n", indent, nestedFieldName))
					sb.WriteString(fmt.Sprintf("%s\t\t\tif v, ok := blockData[\"%s\"].(string); ok && v != \"\" {\n", indent, nestedJsonName))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\treturn types.StringValue(v)\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\treturn types.StringNull()\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t}(),\n", indent))
				case "int64":
					sb.WriteString(fmt.Sprintf("%s\t\t%s: func() types.Int64 {\n", indent, nestedFieldName))
					if nestedAttr.Optional {
						sb.WriteString(fmt.Sprintf("%s\t\t\tif !isImport && data.%s != nil {\n", indent, fieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t// Preserve existing state (null or user-set value)\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t// This prevents API defaults (like 0) from overwriting user intent\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\treturn data.%s.%s\n", indent, fieldName, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\tif !isImport {\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t// Block not in user config - return null, not API default\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\treturn types.Int64Null()\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
					}
					sb.WriteString(fmt.Sprintf("%s\t\t\t// Import case: read from API\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\tif v, ok := blockData[\"%s\"].(float64); ok {\n", indent, nestedJsonName))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\treturn types.Int64Value(int64(v))\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\treturn types.Int64Null()\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t}(),\n", indent))
				case "bool":
					sb.WriteString(fmt.Sprintf("%s\t\t%s: func() types.Bool {\n", indent, nestedFieldName))
					if nestedAttr.Optional {
						sb.WriteString(fmt.Sprintf("%s\t\t\tif !isImport && data.%s != nil {\n", indent, fieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t// Preserve existing state (null or user-set value)\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t// This prevents API defaults from overwriting user intent\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\treturn data.%s.%s\n", indent, fieldName, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\tif !isImport {\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t// Block not in user config - return null, not API default\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\treturn types.BoolNull()\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
					}
					sb.WriteString(fmt.Sprintf("%s\t\t\t// Import case: read from API\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\tif v, ok := blockData[\"%s\"].(bool); ok {\n", indent, nestedJsonName))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\treturn types.BoolValue(v)\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\treturn types.BoolNull()\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t}(),\n", indent))
				case "list":
					// Handle list types inside single nested blocks
					if nestedAttr.ElementType == "string" {
						sb.WriteString(fmt.Sprintf("%s\t\t%s: func() types.List {\n", indent, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t\tif v, ok := blockData[\"%s\"].([]interface{}); ok && len(v) > 0 {\n", indent, nestedJsonName))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\tvar items []string\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\tfor _, item := range v {\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif s, ok := item.(string); ok {\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\titems = append(items, s)\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\tlistVal, _ := types.ListValueFrom(ctx, types.StringType, items)\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\treturn listVal\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\treturn types.ListNull(types.StringType)\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t}(),\n", indent))
					} else if nestedAttr.ElementType == "int64" {
						sb.WriteString(fmt.Sprintf("%s\t\t%s: func() types.List {\n", indent, nestedFieldName))
						sb.WriteString(fmt.Sprintf("%s\t\t\tif v, ok := blockData[\"%s\"].([]interface{}); ok && len(v) > 0 {\n", indent, nestedJsonName))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\tvar items []int64\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\tfor _, item := range v {\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif n, ok := item.(float64); ok {\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\titems = append(items, int64(n))\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\tlistVal, _ := types.ListValueFrom(ctx, types.Int64Type, items)\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t\treturn listVal\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t\treturn types.ListNull(types.Int64Type)\n", indent))
						sb.WriteString(fmt.Sprintf("%s\t\t}(),\n", indent))
					}
				}
				// Handle nested list blocks within single nested blocks (e.g., rules in mitigation_type)
				if nestedAttr.IsBlock && nestedAttr.NestedBlockType == "list" {
					nestedListModelType := resourceTitleCase + attr.GoName + nestedAttr.GoName + "Model"

					sb.WriteString(fmt.Sprintf("%s\t\t%s: func() []%s {\n", indent, nestedFieldName, nestedListModelType))
					sb.WriteString(fmt.Sprintf("%s\t\t\tif listData, ok := blockData[\"%s\"].([]interface{}); ok && len(listData) > 0 {\n", indent, nestedJsonName))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\tvar result []%s\n", indent, nestedListModelType))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\tfor _, item := range listData {\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\t\tif itemMap, ok := item.(map[string]interface{}); ok {\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tresult = append(result, %s{\n", indent, nestedListModelType))

					// Process nested attributes within the list items
					for _, deepAttr := range nestedAttr.NestedAttributes {
						deepFieldName := deepAttr.GoName
						deepJsonName := deepAttr.JsonName
						if deepJsonName == "" {
							deepJsonName = deepAttr.TfsdkTag
						}

						// Handle single nested blocks within list items (like threat_level, mitigation_action)
						if deepAttr.IsBlock && deepAttr.NestedBlockType == "single" {
							deepModelType := resourceTitleCase + attr.GoName + nestedAttr.GoName + deepAttr.GoName + "Model"

							if len(deepAttr.NestedAttributes) == 0 {
								// Empty marker block - check if key exists
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t%s: func() *%sEmptyModel {\n", indent, deepFieldName, resourceTitleCase))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\tif _, ok := itemMap[\"%s\"].(map[string]interface{}); ok {\n", indent, deepJsonName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\treturn &%sEmptyModel{}\n", indent, resourceTitleCase))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t}\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\treturn nil\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t}(),\n", indent))
							} else {
								// Block with nested attributes
								needsDeepMap := false
								for _, checkAttr := range deepAttr.NestedAttributes {
									if checkAttr.IsBlock && checkAttr.NestedBlockType == "single" && len(checkAttr.NestedAttributes) == 0 {
										needsDeepMap = true
										break
									}
									switch checkAttr.Type {
									case "string", "int64", "bool":
										needsDeepMap = true
									}
									if needsDeepMap {
										break
									}
								}

								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t%s: func() *%s {\n", indent, deepFieldName, deepModelType))
								if needsDeepMap {
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\tif deepMap, ok := itemMap[\"%s\"].(map[string]interface{}); ok {\n", indent, deepJsonName))
								} else {
									sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\tif _, ok := itemMap[\"%s\"].(map[string]interface{}); ok {\n", indent, deepJsonName))
								}
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\treturn &%s{\n", indent, deepModelType))

								// Handle the deepest level attributes
								for _, leafAttr := range deepAttr.NestedAttributes {
									leafFieldName := leafAttr.GoName
									leafJsonName := leafAttr.JsonName
									if leafJsonName == "" {
										leafJsonName = leafAttr.TfsdkTag
									}

									if leafAttr.IsBlock && leafAttr.NestedBlockType == "single" && len(leafAttr.NestedAttributes) == 0 {
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t%s: func() *%sEmptyModel {\n", indent, leafFieldName, resourceTitleCase))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\tif _, ok := deepMap[\"%s\"].(map[string]interface{}); ok {\n", indent, leafJsonName))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\t\treturn &%sEmptyModel{}\n", indent, resourceTitleCase))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\t}\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\treturn nil\n", indent))
										sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t}(),\n", indent))
									} else {
										switch leafAttr.Type {
										case "string":
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t%s: func() types.String {\n", indent, leafFieldName))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\tif v, ok := deepMap[\"%s\"].(string); ok && v != \"\" {\n", indent, leafJsonName))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\t\treturn types.StringValue(v)\n", indent))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\t}\n", indent))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\treturn types.StringNull()\n", indent))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t}(),\n", indent))
										case "int64":
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t%s: func() types.Int64 {\n", indent, leafFieldName))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\tif v, ok := deepMap[\"%s\"].(float64); ok {\n", indent, leafJsonName))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\t\treturn types.Int64Value(int64(v))\n", indent))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\t}\n", indent))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\treturn types.Int64Null()\n", indent))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t}(),\n", indent))
										case "bool":
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t%s: func() types.Bool {\n", indent, leafFieldName))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\tif v, ok := deepMap[\"%s\"].(bool); ok {\n", indent, leafJsonName))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\t\treturn types.BoolValue(v)\n", indent))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\t}\n", indent))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t\treturn types.BoolNull()\n", indent))
											sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t\t}(),\n", indent))
										}
									}
								}

								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\t}\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t}\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\treturn nil\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t}(),\n", indent))
							}
						} else {
							// Handle primitive types within list items
							switch deepAttr.Type {
							case "string":
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t%s: func() types.String {\n", indent, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\tif v, ok := itemMap[\"%s\"].(string); ok && v != \"\" {\n", indent, deepJsonName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\treturn types.StringValue(v)\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t}\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\treturn types.StringNull()\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t}(),\n", indent))
							case "int64":
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t%s: func() types.Int64 {\n", indent, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\tif v, ok := itemMap[\"%s\"].(float64); ok {\n", indent, deepJsonName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\treturn types.Int64Value(int64(v))\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t}\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\treturn types.Int64Null()\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t}(),\n", indent))
							case "bool":
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t%s: func() types.Bool {\n", indent, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\tif v, ok := itemMap[\"%s\"].(bool); ok {\n", indent, deepJsonName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\treturn types.BoolValue(v)\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t}\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\treturn types.BoolNull()\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t}(),\n", indent))
							}
						}
					}

					sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t})\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\treturn result\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\treturn nil\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t}(),\n", indent))
				}
				// Handle single nested blocks with attributes within single nested blocks
				if nestedAttr.IsBlock && nestedAttr.NestedBlockType == "single" && len(nestedAttr.NestedAttributes) > 0 {
					nestedSingleModelType := resourceTitleCase + attr.GoName + nestedAttr.GoName + "Model"

					// Check if any attributes will actually use nestedBlockData
					hasUsableAttrs := false
					for _, da := range nestedAttr.NestedAttributes {
						switch da.Type {
						case "string", "int64", "bool":
							hasUsableAttrs = true
						case "list":
							if da.ElementType == "string" || da.ElementType == "int64" {
								hasUsableAttrs = true
							}
						}
					}

					sb.WriteString(fmt.Sprintf("%s\t\t%s: func() *%s {\n", indent, nestedFieldName, nestedSingleModelType))
					sb.WriteString(fmt.Sprintf("%s\t\t\tif !isImport && data.%s != nil && data.%s.%s != nil {\n", indent, fieldName, fieldName, nestedFieldName))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\t// Normal Read: preserve existing state value\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\treturn data.%s.%s\n", indent, fieldName, nestedFieldName))
					sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\t// Import case: read from API\n", indent))

					// Use _ if no attributes will use nestedBlockData
					localVarName := "nestedBlockData"
					if !hasUsableAttrs {
						localVarName = "_"
					}
					sb.WriteString(fmt.Sprintf("%s\t\t\tif %s, ok := blockData[\"%s\"].(map[string]interface{}); ok {\n", indent, localVarName, nestedJsonName))
					sb.WriteString(fmt.Sprintf("%s\t\t\t\treturn &%s{\n", indent, nestedSingleModelType))

					// Process attributes within the single nested block
					for _, deepAttr := range nestedAttr.NestedAttributes {
						deepFieldName := deepAttr.GoName
						deepJsonName := deepAttr.JsonName
						if deepJsonName == "" {
							deepJsonName = deepAttr.TfsdkTag
						}

						switch deepAttr.Type {
						case "string":
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t%s: func() types.String {\n", indent, deepFieldName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tif v, ok := nestedBlockData[\"%s\"].(string); ok && v != \"\" {\n", indent, deepJsonName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\treturn types.StringValue(v)\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t}\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\treturn types.StringNull()\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}(),\n", indent))
						case "int64":
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t%s: func() types.Int64 {\n", indent, deepFieldName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tif v, ok := nestedBlockData[\"%s\"].(float64); ok {\n", indent, deepJsonName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\treturn types.Int64Value(int64(v))\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t}\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\treturn types.Int64Null()\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}(),\n", indent))
						case "bool":
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t%s: func() types.Bool {\n", indent, deepFieldName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tif v, ok := nestedBlockData[\"%s\"].(bool); ok {\n", indent, deepJsonName))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\treturn types.BoolValue(v)\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t}\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\treturn types.BoolNull()\n", indent))
							sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}(),\n", indent))
						case "list":
							if deepAttr.ElementType == "string" {
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t%s: func() types.List {\n", indent, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tif v, ok := nestedBlockData[\"%s\"].([]interface{}); ok && len(v) > 0 {\n", indent, deepJsonName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\tvar items []string\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\tfor _, item := range v {\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\tif s, ok := item.(string); ok {\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\titems = append(items, s)\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t}\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t}\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\tlistVal, _ := types.ListValueFrom(ctx, types.StringType, items)\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\treturn listVal\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t}\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\treturn types.ListNull(types.StringType)\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}(),\n", indent))
							} else if deepAttr.ElementType == "int64" {
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t%s: func() types.List {\n", indent, deepFieldName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\tif v, ok := nestedBlockData[\"%s\"].([]interface{}); ok && len(v) > 0 {\n", indent, deepJsonName))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\tvar items []int64\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\tfor _, item := range v {\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\tif n, ok := item.(float64); ok {\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t\titems = append(items, int64(n))\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t\t}\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\t}\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\tlistVal, _ := types.ListValueFrom(ctx, types.Int64Type, items)\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\treturn listVal\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t}\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t\treturn types.ListNull(types.Int64Type)\n", indent))
								sb.WriteString(fmt.Sprintf("%s\t\t\t\t\t}(),\n", indent))
							}
						}
					}

					sb.WriteString(fmt.Sprintf("%s\t\t\t\t}\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t\treturn nil\n", indent))
					sb.WriteString(fmt.Sprintf("%s\t\t}(),\n", indent))
				}
			}
			sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
			sb.WriteString(fmt.Sprintf("%s}\n", indent))
			continue
		}

		switch attr.Type {
		case "string":
			sb.WriteString(fmt.Sprintf("%sif v, ok := apiResource.Spec[\"%s\"].(string); ok && v != \"\" {\n", indent, jsonName))
			sb.WriteString(fmt.Sprintf("%s\tdata.%s = types.StringValue(v)\n", indent, fieldName))
			sb.WriteString(fmt.Sprintf("%s} else {\n", indent))
			sb.WriteString(fmt.Sprintf("%s\tdata.%s = types.StringNull()\n", indent, fieldName))
			sb.WriteString(fmt.Sprintf("%s}\n", indent))
		case "int64":
			sb.WriteString(fmt.Sprintf("%sif v, ok := apiResource.Spec[\"%s\"].(float64); ok {\n", indent, jsonName))
			sb.WriteString(fmt.Sprintf("%s\tdata.%s = types.Int64Value(int64(v))\n", indent, fieldName))
			sb.WriteString(fmt.Sprintf("%s} else {\n", indent))
			sb.WriteString(fmt.Sprintf("%s\tdata.%s = types.Int64Null()\n", indent, fieldName))
			sb.WriteString(fmt.Sprintf("%s}\n", indent))
		case "bool":
			if attr.Optional {
				sb.WriteString(fmt.Sprintf("%s// Top-level Optional bool: preserve prior state to avoid API default drift\n", indent))
				sb.WriteString(fmt.Sprintf("%sif !isImport && !data.%s.IsNull() && !data.%s.IsUnknown() {\n", indent, fieldName, fieldName))
				sb.WriteString(fmt.Sprintf("%s\t// Normal Read: preserve existing state value (do nothing)\n", indent))
				sb.WriteString(fmt.Sprintf("%s} else {\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t// Import case, null state, or unknown (after Create): read from API\n", indent))
				sb.WriteString(fmt.Sprintf("%s\tif v, ok := apiResource.Spec[\"%s\"].(bool); ok {\n", indent, jsonName))
				sb.WriteString(fmt.Sprintf("%s\t\tdata.%s = types.BoolValue(v)\n", indent, fieldName))
				sb.WriteString(fmt.Sprintf("%s\t} else {\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t\tdata.%s = types.BoolNull()\n", indent, fieldName))
				sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
				sb.WriteString(fmt.Sprintf("%s}\n", indent))
			} else {
				sb.WriteString(fmt.Sprintf("%sif v, ok := apiResource.Spec[\"%s\"].(bool); ok {\n", indent, jsonName))
				sb.WriteString(fmt.Sprintf("%s\tdata.%s = types.BoolValue(v)\n", indent, fieldName))
				sb.WriteString(fmt.Sprintf("%s} else {\n", indent))
				sb.WriteString(fmt.Sprintf("%s\tdata.%s = types.BoolNull()\n", indent, fieldName))
				sb.WriteString(fmt.Sprintf("%s}\n", indent))
			}
		case "list":
			if attr.ElementType == "string" {
				sb.WriteString(fmt.Sprintf("%sif v, ok := apiResource.Spec[\"%s\"].([]interface{}); ok && len(v) > 0 {\n", indent, jsonName))
				sb.WriteString(fmt.Sprintf("%s\tvar %sList []string\n", indent, attr.TfsdkTag))
				sb.WriteString(fmt.Sprintf("%s\tfor _, item := range v {\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t\tif s, ok := item.(string); ok {\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t\t\t%sList = append(%sList, s)\n", indent, attr.TfsdkTag, attr.TfsdkTag))
				sb.WriteString(fmt.Sprintf("%s\t\t}\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
				sb.WriteString(fmt.Sprintf("%s\tlistVal, diags := types.ListValueFrom(ctx, types.StringType, %sList)\n", indent, attr.TfsdkTag))
				sb.WriteString(fmt.Sprintf("%s\tresp.Diagnostics.Append(diags...)\n", indent))
				sb.WriteString(fmt.Sprintf("%s\tif !resp.Diagnostics.HasError() {\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t\tdata.%s = listVal\n", indent, fieldName))
				sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
				sb.WriteString(fmt.Sprintf("%s} else {\n", indent))
				sb.WriteString(fmt.Sprintf("%s\tdata.%s = types.ListNull(types.StringType)\n", indent, fieldName))
				sb.WriteString(fmt.Sprintf("%s}\n", indent))
			} else if attr.ElementType == "int64" {
				sb.WriteString(fmt.Sprintf("%sif v, ok := apiResource.Spec[\"%s\"].([]interface{}); ok && len(v) > 0 {\n", indent, jsonName))
				sb.WriteString(fmt.Sprintf("%s\tvar %sList []int64\n", indent, attr.TfsdkTag))
				sb.WriteString(fmt.Sprintf("%s\tfor _, item := range v {\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t\tif n, ok := item.(float64); ok {\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t\t\t%sList = append(%sList, int64(n))\n", indent, attr.TfsdkTag, attr.TfsdkTag))
				sb.WriteString(fmt.Sprintf("%s\t\t}\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
				sb.WriteString(fmt.Sprintf("%s\tlistVal, diags := types.ListValueFrom(ctx, types.Int64Type, %sList)\n", indent, attr.TfsdkTag))
				sb.WriteString(fmt.Sprintf("%s\tresp.Diagnostics.Append(diags...)\n", indent))
				sb.WriteString(fmt.Sprintf("%s\tif !resp.Diagnostics.HasError() {\n", indent))
				sb.WriteString(fmt.Sprintf("%s\t\tdata.%s = listVal\n", indent, fieldName))
				sb.WriteString(fmt.Sprintf("%s\t}\n", indent))
				sb.WriteString(fmt.Sprintf("%s} else {\n", indent))
				sb.WriteString(fmt.Sprintf("%s\tdata.%s = types.ListNull(types.Int64Type)\n", indent, fieldName))
				sb.WriteString(fmt.Sprintf("%s}\n", indent))
			}
		}
	}
	return sb.String()
}

// RenderNestedAttributes generates the Attributes map for nested blocks
func RenderNestedAttributes(attrs []openapi.TerraformAttribute, indent string) string {
	if len(attrs) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(indent + "Attributes: map[string]schema.Attribute{\n")

	for _, attr := range attrs {
		if attr.IsBlock {
			continue // Blocks are handled separately
		}

		attrType := "String"
		switch attr.Type {
		case "int64":
			attrType = "Int64"
		case "bool":
			attrType = "Bool"
		case "map":
			attrType = "Map"
		case "list":
			attrType = "List"
		}

		// Escape backslashes and quotes in descriptions for Go string literals
		desc := EscapeGoString(attr.Description)

		sb.WriteString(fmt.Sprintf("%s\t\"%s\": schema.%sAttribute{\n", indent, attr.TfsdkTag, attrType))
		sb.WriteString(fmt.Sprintf("%s\t\tMarkdownDescription: \"%s\",\n", indent, desc))

		if attr.DeprecationMessage != "" {
			sb.WriteString(fmt.Sprintf("%s\t\tDeprecationMessage: \"%s\",\n", indent, EscapeGoString(attr.DeprecationMessage)))
		}

		if attr.Required {
			sb.WriteString(fmt.Sprintf("%s\t\tRequired: true,\n", indent))
		} else {
			// Handle Optional and Computed flags (can be both true for fields like 'tenant')
			if attr.Optional {
				sb.WriteString(fmt.Sprintf("%s\t\tOptional: true,\n", indent))
			}
			if attr.Computed {
				sb.WriteString(fmt.Sprintf("%s\t\tComputed: true,\n", indent))
			}
		}

		if attr.PlanModifier != "" {
			typeName := "String"
			pkgName := "stringplanmodifier"
			switch attr.Type {
			case "bool":
				typeName = "Bool"
				pkgName = "boolplanmodifier"
			case "int64":
				typeName = "Int64"
				pkgName = "int64planmodifier"
			}
			sb.WriteString(fmt.Sprintf("%s\t\tPlanModifiers: []planmodifier.%s{\n", indent, typeName))
			switch attr.PlanModifier {
			case "RequiresReplace":
				sb.WriteString(fmt.Sprintf("%s\t\t\t%s.RequiresReplace(),\n", indent, pkgName))
			default:
				sb.WriteString(fmt.Sprintf("%s\t\t\t%s.UseStateForUnknown(),\n", indent, pkgName))
			}
			sb.WriteString(fmt.Sprintf("%s\t\t},\n", indent))
		}

		if attr.Type == "map" || attr.Type == "list" {
			// Map ElementType to Terraform attr type
			elementTfType := "types.StringType"
			switch attr.ElementType {
			case "int64":
				elementTfType = "types.Int64Type"
			case "bool":
				elementTfType = "types.BoolType"
			}
			sb.WriteString(fmt.Sprintf("%s\t\tElementType: %s,\n", indent, elementTfType))
		}

		// Add string validators (LengthBetween/LengthAtMost/LengthAtLeast, RegexMatches, OneOf)
		if attr.Type == "string" {
			var stringValidators []string
			if attr.MinLength > 0 && attr.MaxLength > 0 {
				stringValidators = append(stringValidators, fmt.Sprintf("stringvalidator.LengthBetween(%d, %d)", attr.MinLength, attr.MaxLength))
			} else if attr.MaxLength > 0 {
				stringValidators = append(stringValidators, fmt.Sprintf("stringvalidator.LengthAtMost(%d)", attr.MaxLength))
			} else if attr.MinLength > 0 {
				stringValidators = append(stringValidators, fmt.Sprintf("stringvalidator.LengthAtLeast(%d)", attr.MinLength))
			}
			if attr.Pattern != "" {
				stringValidators = append(stringValidators, fmt.Sprintf("stringvalidator.RegexMatches(regexp.MustCompile(%s), \"\")", RegexLiteral(attr.Pattern)))
			}
			if len(attr.EnumValues) > 0 {
				quoted := make([]string, len(attr.EnumValues))
				for i, v := range attr.EnumValues {
					quoted[i] = fmt.Sprintf("%q", v)
				}
				stringValidators = append(stringValidators, fmt.Sprintf("stringvalidator.OneOf(%s)", strings.Join(quoted, ", ")))
			}
			if len(stringValidators) > 0 {
				sb.WriteString(fmt.Sprintf("%s\t\tValidators: []validator.String{\n", indent))
				for _, sv := range stringValidators {
					sb.WriteString(fmt.Sprintf("%s\t\t\t%s,\n", indent, sv))
				}
				sb.WriteString(fmt.Sprintf("%s\t\t},\n", indent))
			}
		}

		// List/set size validators
		if (attr.Type == "list" || attr.Type == "set") && (attr.MinItems > 0 || attr.MaxItems > 0) {
			sb.WriteString(fmt.Sprintf("%s\t\tValidators: []validator.List{\n", indent))
			if attr.MinItems > 0 && attr.MaxItems > 0 {
				sb.WriteString(fmt.Sprintf("%s\t\t\tlistvalidator.SizeBetween(%d, %d),\n", indent, attr.MinItems, attr.MaxItems))
			} else if attr.MaxItems > 0 {
				sb.WriteString(fmt.Sprintf("%s\t\t\tlistvalidator.SizeAtMost(%d),\n", indent, attr.MaxItems))
			} else {
				sb.WriteString(fmt.Sprintf("%s\t\t\tlistvalidator.SizeAtLeast(%d),\n", indent, attr.MinItems))
			}
			sb.WriteString(fmt.Sprintf("%s\t\t},\n", indent))
		}

		sb.WriteString(fmt.Sprintf("%s\t},\n", indent))
	}

	sb.WriteString(indent + "},\n")
	return sb.String()
}

// CollectNestedModelTypes recursively collects all nested model type definitions
func CollectNestedModelTypes(resourceTitleCase string, attrs []openapi.TerraformAttribute, prefix string, collected *[]NestedModelInfo) {
	for _, attr := range attrs {
		if !attr.IsBlock {
			continue
		}

		// Build the type name: ResourceName + Prefix + AttributeName + Model
		currentPrefix := prefix + naming.ToResourceTypeName(attr.TfsdkTag)
		typeName := resourceTitleCase + currentPrefix + "Model"

		// Check if this block has any nested attributes (non-empty)
		hasContent := false
		for _, nested := range attr.NestedAttributes {
			if !nested.IsBlock {
				hasContent = true
				break
			}
		}

		// Also check for nested blocks
		for _, nested := range attr.NestedAttributes {
			if nested.IsBlock {
				hasContent = true
				break
			}
		}

		*collected = append(*collected, NestedModelInfo{
			TypeName:    typeName,
			Description: attr.TfsdkTag,
			Prefix:      currentPrefix, // Store the full prefix for generating nested type references
			Attributes:  attr.NestedAttributes,
			IsEmpty:     !hasContent && len(attr.NestedAttributes) == 0,
		})

		// Recursively collect nested types
		if len(attr.NestedAttributes) > 0 {
			CollectNestedModelTypes(resourceTitleCase, attr.NestedAttributes, currentPrefix, collected)
		}
	}
}

// RenderNestedModelTypes generates all nested model struct definitions for a resource
func RenderNestedModelTypes(resourceTitleCase string, attrs []openapi.TerraformAttribute) string {
	// First check if there are any blocks
	hasBlocks := false
	for _, attr := range attrs {
		if attr.IsBlock {
			hasBlocks = true
			break
		}
	}
	if !hasBlocks {
		return ""
	}

	var sb strings.Builder

	// Add empty model type for blocks with no attributes
	sb.WriteString(fmt.Sprintf("// %sEmptyModel represents empty nested blocks\n", resourceTitleCase))
	sb.WriteString(fmt.Sprintf("type %sEmptyModel struct {\n}\n\n", resourceTitleCase))

	// Collect all nested model types
	var models []NestedModelInfo
	CollectNestedModelTypes(resourceTitleCase, attrs, "", &models)

	// Generate each model type
	for _, model := range models {
		if model.IsEmpty {
			continue // Empty models use the shared EmptyModel
		}

		sb.WriteString(fmt.Sprintf("// %s represents %s block\n", model.TypeName, model.Description))
		sb.WriteString(fmt.Sprintf("type %s struct {\n", model.TypeName))

		// Generate fields for non-block attributes
		for _, attr := range model.Attributes {
			if attr.IsBlock {
				continue
			}
			goType := "String"
			switch attr.Type {
			case "int64":
				goType = "Int64"
			case "bool":
				goType = "Bool"
			case "map":
				goType = "Map"
			case "list":
				goType = "List"
			}
			sb.WriteString(fmt.Sprintf("\t%s types.%s `tfsdk:\"%s\"`\n", attr.GoName, goType, attr.TfsdkTag))
		}

		// Generate pointer fields for nested block attributes
		for _, attr := range model.Attributes {
			if !attr.IsBlock {
				continue
			}
			// Use the model's full prefix to build the nested type name
			nestedTypeName := resourceTitleCase + model.Prefix + naming.ToResourceTypeName(attr.TfsdkTag) + "Model"

			// Check if this nested block is empty
			isNestedEmpty := len(attr.NestedAttributes) == 0
			if !isNestedEmpty {
				hasNonBlockAttrs := false
				for _, nested := range attr.NestedAttributes {
					if !nested.IsBlock {
						hasNonBlockAttrs = true
						break
					}
				}
				isNestedEmpty = !hasNonBlockAttrs && len(attr.NestedAttributes) == 0
			}

			if isNestedEmpty {
				nestedTypeName = resourceTitleCase + "EmptyModel"
			}

			// For nested blocks within other models, use Go slices/pointers (not types.List)
			if attr.NestedBlockType == "list" {
				sb.WriteString(fmt.Sprintf("\t%s []%s `tfsdk:\"%s\"`\n", attr.GoName, nestedTypeName, attr.TfsdkTag))
			} else {
				sb.WriteString(fmt.Sprintf("\t%s *%s `tfsdk:\"%s\"`\n", attr.GoName, nestedTypeName, attr.TfsdkTag))
			}
		}

		sb.WriteString("}\n\n")

		// Generate AttrTypes map for this model (needed for types.List conversion)
		sb.WriteString(fmt.Sprintf("// %sAttrTypes defines the attribute types for %s\n", model.TypeName, model.TypeName))
		sb.WriteString(fmt.Sprintf("var %sAttrTypes = map[string]attr.Type{\n", model.TypeName))

		// Generate attr.Type for non-block attributes
		for _, attr := range model.Attributes {
			if attr.IsBlock {
				continue
			}
			var attrType string
			switch attr.Type {
			case "string":
				attrType = "types.StringType"
			case "int64":
				attrType = "types.Int64Type"
			case "bool":
				attrType = "types.BoolType"
			case "map":
				attrType = "types.MapType{ElemType: types.StringType}"
			case "list":
				elemType := "types.StringType"
				switch attr.ElementType {
				case "int64":
					elemType = "types.Int64Type"
				case "bool":
					elemType = "types.BoolType"
				}
				attrType = fmt.Sprintf("types.ListType{ElemType: %s}", elemType)
			default:
				attrType = "types.StringType"
			}
			sb.WriteString(fmt.Sprintf("\t\"%s\": %s,\n", attr.TfsdkTag, attrType))
		}

		// Generate attr.Type for block attributes
		for _, attr := range model.Attributes {
			if !attr.IsBlock {
				continue
			}

			isNestedEmpty := len(attr.NestedAttributes) == 0

			if isNestedEmpty {
				// Empty nested block - use inline empty map
				if attr.NestedBlockType == "list" {
					sb.WriteString(fmt.Sprintf("\t\"%s\": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{}}},\n", attr.TfsdkTag))
				} else {
					sb.WriteString(fmt.Sprintf("\t\"%s\": types.ObjectType{AttrTypes: map[string]attr.Type{}},\n", attr.TfsdkTag))
				}
			} else {
				// Non-empty nested block - reference its AttrTypes variable
				nestedAttrTypesName := resourceTitleCase + model.Prefix + naming.ToResourceTypeName(attr.TfsdkTag) + "ModelAttrTypes"
				if attr.NestedBlockType == "list" {
					sb.WriteString(fmt.Sprintf("\t\"%s\": types.ListType{ElemType: types.ObjectType{AttrTypes: %s}},\n", attr.TfsdkTag, nestedAttrTypesName))
				} else {
					sb.WriteString(fmt.Sprintf("\t\"%s\": types.ObjectType{AttrTypes: %s},\n", attr.TfsdkTag, nestedAttrTypesName))
				}
			}
		}

		sb.WriteString("}\n\n")
	}

	return sb.String()
}

// RenderBlockFields generates the block fields for the main ResourceModel struct
func RenderBlockFields(resourceTitleCase string, attrs []openapi.TerraformAttribute) string {
	var sb strings.Builder

	for _, attr := range attrs {
		if !attr.IsBlock {
			continue
		}

		// Determine the type name
		typeName := resourceTitleCase + naming.ToResourceTypeName(attr.TfsdkTag) + "Model"

		// Check if this block is empty
		isBlockEmpty := len(attr.NestedAttributes) == 0
		if !isBlockEmpty {
			hasNonBlockAttrs := false
			for _, nested := range attr.NestedAttributes {
				if !nested.IsBlock {
					hasNonBlockAttrs = true
					break
				}
			}
			// Also check for any nested blocks
			hasNestedBlocks := false
			for _, nested := range attr.NestedAttributes {
				if nested.IsBlock {
					hasNestedBlocks = true
					break
				}
			}
			isBlockEmpty = !hasNonBlockAttrs && !hasNestedBlocks
		}

		if isBlockEmpty {
			typeName = resourceTitleCase + "EmptyModel"
		}

		// For list nested blocks, use types.List to properly handle unknown values during planning.
		// For single nested blocks, use pointer type.
		if attr.NestedBlockType == "list" {
			sb.WriteString(fmt.Sprintf("\t%s types.List `tfsdk:\"%s\"`\n", attr.GoName, attr.TfsdkTag))
		} else {
			sb.WriteString(fmt.Sprintf("\t%s *%s `tfsdk:\"%s\"`\n", attr.GoName, typeName, attr.TfsdkTag))
		}
	}

	return sb.String()
}

// RenderNestedBlocks generates the Blocks map for nested blocks within a block
func RenderNestedBlocks(attrs []openapi.TerraformAttribute, indent string) string {
	var hasBlocks bool
	for _, attr := range attrs {
		if attr.IsBlock {
			hasBlocks = true
			break
		}
	}

	if !hasBlocks {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(indent + "Blocks: map[string]schema.Block{\n")

	for _, attr := range attrs {
		if !attr.IsBlock {
			continue
		}

		blockType := "SingleNestedBlock"
		if attr.NestedBlockType == "list" {
			blockType = "ListNestedBlock"
		}

		// Escape backslashes and quotes in descriptions
		desc := EscapeGoString(attr.Description)

		sb.WriteString(fmt.Sprintf("%s\t\"%s\": schema.%s{\n", indent, attr.TfsdkTag, blockType))
		sb.WriteString(fmt.Sprintf("%s\t\tMarkdownDescription: \"%s\",\n", indent, desc))

		if attr.NestedBlockType == "list" {
			sb.WriteString(fmt.Sprintf("%s\t\tNestedObject: schema.NestedBlockObject{\n", indent))
			if len(attr.NestedAttributes) > 0 {
				sb.WriteString(RenderNestedAttributes(attr.NestedAttributes, indent+"\t\t\t"))
				sb.WriteString(RenderNestedBlocks(attr.NestedAttributes, indent+"\t\t\t"))
			} else {
				sb.WriteString(fmt.Sprintf("%s\t\t\tAttributes: map[string]schema.Attribute{},\n", indent))
			}
			sb.WriteString(fmt.Sprintf("%s\t\t},\n", indent))
		} else {
			// SingleNestedBlock
			if len(attr.NestedAttributes) > 0 {
				sb.WriteString(RenderNestedAttributes(attr.NestedAttributes, indent+"\t\t"))
				sb.WriteString(RenderNestedBlocks(attr.NestedAttributes, indent+"\t\t"))
			}
		}

		sb.WriteString(fmt.Sprintf("%s\t},\n", indent))
	}

	sb.WriteString(indent + "},\n")
	return sb.String()
}
