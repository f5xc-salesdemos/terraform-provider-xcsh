// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package schema

import (
	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/conflicts"
	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/openapi"
)

// HasNestedModelsWithAttrTypes checks recursively if any nested blocks would generate AttrTypes.
// This is needed to determine if the attr import is required.
// AttrTypes are generated for any nested model that has ANY nested attributes (block or non-block).
func HasNestedModelsWithAttrTypes(attributes []openapi.TerraformAttribute) bool {
	for _, attr := range attributes {
		if attr.IsBlock {
			// If this block has ANY nested attributes, AttrTypes will be generated for it
			if len(attr.NestedAttributes) > 0 {
				return true
			}
			// Note: Even if NestedAttributes is empty, we don't need to recurse
			// because empty blocks use EmptyModel which doesn't have AttrTypes
		}
	}
	return false
}

// HasMaxLengthValidatorsAny recursively checks if any attribute (top-level or nested)
// has a non-zero MaxLength. This determines whether stringvalidator must be imported.
func HasMaxLengthValidatorsAny(attributes []openapi.TerraformAttribute) bool {
	for _, attr := range attributes {
		if attr.MaxLength > 0 {
			return true
		}
		if len(attr.NestedAttributes) > 0 {
			if HasMaxLengthValidatorsAny(attr.NestedAttributes) {
				return true
			}
		}
	}
	return false
}

// HasEnumValidatorsAny recursively checks if any attribute (top-level or nested)
// has EnumValues. This determines whether stringvalidator must be imported for OneOf.
func HasEnumValidatorsAny(attributes []openapi.TerraformAttribute) bool {
	for _, attr := range attributes {
		if len(attr.EnumValues) > 0 {
			return true
		}
		if HasEnumValidatorsAny(attr.NestedAttributes) {
			return true
		}
	}
	return false
}

// HasPatternValidatorsAny recursively checks if any attribute (top-level or nested)
// has a Pattern regex. This determines whether regexp must be imported.
func HasPatternValidatorsAny(attributes []openapi.TerraformAttribute) bool {
	for _, attr := range attributes {
		if attr.Pattern != "" {
			return true
		}
		if HasPatternValidatorsAny(attr.NestedAttributes) {
			return true
		}
	}
	return false
}

// HasListSizeValidatorsAny recursively checks if any non-block list/set attribute (top-level or nested)
// has MinItems or MaxItems constraints. This determines whether listvalidator must be imported.
// Only non-block attributes are checked because blocks use schema.ListNestedBlock which
// does not support the same validator interface.
func HasListSizeValidatorsAny(attributes []openapi.TerraformAttribute) bool {
	for _, attr := range attributes {
		if !attr.IsBlock && (attr.MinItems > 0 || attr.MaxItems > 0) && (attr.Type == "list" || attr.Type == "set") {
			return true
		}
		if HasListSizeValidatorsAny(attr.NestedAttributes) {
			return true
		}
	}
	return false
}

// CollectConflictAttrs collects top-level non-block attributes that have ConflictsWith relationships.
// Block attributes are excluded because their Go types are pointers to nested models (not framework
// value types) and do not have the IsNull() method needed for conflict checking.
func CollectConflictAttrs(attributes []openapi.TerraformAttribute) ([]conflicts.Attr, map[string]string) {
	// Build a lookup map from tfsdk tag to Go field name for non-block attributes only.
	// This is used to resolve conflict target names and to filter out block targets.
	goNameLookup := make(map[string]string)
	for _, attr := range attributes {
		if attr.IsBlock {
			continue
		}
		goNameLookup[attr.TfsdkTag] = attr.GoName
	}

	var result []conflicts.Attr
	for _, attr := range attributes {
		if attr.IsBlock {
			continue
		}
		if len(attr.ConflictsWith) > 0 {
			result = append(result, conflicts.Attr{
				TfsdkTag:      attr.TfsdkTag,
				GoName:        attr.GoName,
				ConflictsWith: attr.ConflictsWith,
			})
		}
	}
	return result, goNameLookup
}

// HasInt64RangeValidatorsAny returns true if any attribute or nested attribute has Minimum or Maximum set.
func HasInt64RangeValidatorsAny(attributes []openapi.TerraformAttribute) bool {
	for _, attr := range attributes {
		if attr.Minimum > 0 || attr.Maximum > 0 {
			return true
		}
		for _, nested := range attr.NestedAttributes {
			if nested.Minimum > 0 || nested.Maximum > 0 {
				return true
			}
		}
	}
	return false
}

// ScanPlanModifierUsage recursively scans attributes to determine which plan modifier imports are needed.
func ScanPlanModifierUsage(attributes []openapi.TerraformAttribute) (usesBool, usesInt64, usesString, usesList, usesMap bool) {
	for _, attr := range attributes {
		if attr.PlanModifier != "" && !attr.IsBlock {
			switch attr.Type {
			case "bool":
				usesBool = true
			case "int64":
				usesInt64 = true
			case "string":
				usesString = true
			case "list", "set":
				usesList = true
			case "map":
				usesMap = true
			}
		}
		if len(attr.NestedAttributes) > 0 {
			nb, ni, ns, nl, nm := ScanPlanModifierUsage(attr.NestedAttributes)
			usesBool = usesBool || nb
			usesInt64 = usesInt64 || ni
			usesString = usesString || ns
			usesList = usesList || nl
			usesMap = usesMap || nm
		}
	}
	return
}
