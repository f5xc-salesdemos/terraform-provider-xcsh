// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package schema

import (
	"testing"

	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/openapi"
)

func TestHasNestedModelsWithAttrTypes(t *testing.T) {
	// No blocks
	attrs := []openapi.TerraformAttribute{
		{Name: "field", IsBlock: false},
	}
	if HasNestedModelsWithAttrTypes(attrs) {
		t.Error("Should return false when no blocks exist")
	}

	// Block with nested attributes
	attrs = []openapi.TerraformAttribute{
		{Name: "block", IsBlock: true, NestedAttributes: []openapi.TerraformAttribute{
			{Name: "sub", Type: "string"},
		}},
	}
	if !HasNestedModelsWithAttrTypes(attrs) {
		t.Error("Should return true when block has nested attributes")
	}

	// Block without nested attributes
	attrs = []openapi.TerraformAttribute{
		{Name: "empty_block", IsBlock: true},
	}
	if HasNestedModelsWithAttrTypes(attrs) {
		t.Error("Should return false when block has no nested attributes")
	}
}

func TestHasMaxLengthValidatorsAny(t *testing.T) {
	if HasMaxLengthValidatorsAny(nil) {
		t.Error("Should return false for nil")
	}

	attrs := []openapi.TerraformAttribute{
		{Name: "a", MaxLength: 0},
		{Name: "b", MaxLength: 255},
	}
	if !HasMaxLengthValidatorsAny(attrs) {
		t.Error("Should return true when MaxLength > 0")
	}

	// Nested
	attrs = []openapi.TerraformAttribute{
		{Name: "parent", NestedAttributes: []openapi.TerraformAttribute{
			{Name: "child", MaxLength: 100},
		}},
	}
	if !HasMaxLengthValidatorsAny(attrs) {
		t.Error("Should return true for nested MaxLength > 0")
	}
}

func TestHasEnumValidatorsAny(t *testing.T) {
	attrs := []openapi.TerraformAttribute{
		{Name: "a"},
	}
	if HasEnumValidatorsAny(attrs) {
		t.Error("Should return false when no enums")
	}

	attrs = []openapi.TerraformAttribute{
		{Name: "status", EnumValues: []string{"ACTIVE", "INACTIVE"}},
	}
	if !HasEnumValidatorsAny(attrs) {
		t.Error("Should return true when EnumValues present")
	}
}

func TestHasPatternValidatorsAny(t *testing.T) {
	attrs := []openapi.TerraformAttribute{
		{Name: "a"},
	}
	if HasPatternValidatorsAny(attrs) {
		t.Error("Should return false when no patterns")
	}

	attrs = []openapi.TerraformAttribute{
		{Name: "ip", Pattern: `^\d+\.\d+\.\d+\.\d+$`},
	}
	if !HasPatternValidatorsAny(attrs) {
		t.Error("Should return true when Pattern present")
	}
}

func TestHasListSizeValidatorsAny(t *testing.T) {
	// Block attributes should be excluded
	attrs := []openapi.TerraformAttribute{
		{Name: "block", IsBlock: true, Type: "list", MinItems: 1},
	}
	if HasListSizeValidatorsAny(attrs) {
		t.Error("Should return false for block attributes")
	}

	// Non-block list with constraints
	attrs = []openapi.TerraformAttribute{
		{Name: "items", IsBlock: false, Type: "list", MaxItems: 10},
	}
	if !HasListSizeValidatorsAny(attrs) {
		t.Error("Should return true for non-block list with MaxItems")
	}
}

func TestCollectConflictAttrs(t *testing.T) {
	attrs := []openapi.TerraformAttribute{
		{TfsdkTag: "field_a", GoName: "FieldA", ConflictsWith: []string{"field_b"}, IsBlock: false},
		{TfsdkTag: "field_b", GoName: "FieldB", ConflictsWith: []string{"field_a"}, IsBlock: false},
		{TfsdkTag: "block_c", GoName: "BlockC", ConflictsWith: []string{"field_a"}, IsBlock: true}, // excluded
	}
	result, lookup := CollectConflictAttrs(attrs)
	if len(result) != 2 {
		t.Fatalf("CollectConflictAttrs returned %d attrs, want 2", len(result))
	}
	if lookup["field_a"] != "FieldA" {
		t.Errorf("lookup[field_a] = %q, want %q", lookup["field_a"], "FieldA")
	}
	if _, ok := lookup["block_c"]; ok {
		t.Error("Block attributes should not be in lookup")
	}
}

func TestScanPlanModifierUsage(t *testing.T) {
	attrs := []openapi.TerraformAttribute{
		{Type: "string", PlanModifier: "UseStateForUnknown"},
		{Type: "bool", PlanModifier: "RequiresReplace"},
		{Type: "int64", PlanModifier: ""},
		{NestedAttributes: []openapi.TerraformAttribute{
			{Type: "int64", PlanModifier: "UseStateForUnknown"},
		}},
	}
	usesBool, usesInt64, usesString := ScanPlanModifierUsage(attrs)
	if !usesBool {
		t.Error("Should detect bool plan modifier")
	}
	if !usesInt64 {
		t.Error("Should detect int64 plan modifier from nested")
	}
	if !usesString {
		t.Error("Should detect string plan modifier")
	}
}
