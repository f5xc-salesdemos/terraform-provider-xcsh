// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package schema

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// FieldKind distinguishes block-type fields (use {}) from scalar attributes (use = value).
type FieldKind int

const (
	FieldKindUnknown   FieldKind = iota
	FieldKindAttribute           // scalar: string, bool, number, list, map → use = value
	FieldKindBlock               // nested block → use {}
)

func (k FieldKind) String() string {
	switch k {
	case FieldKindAttribute:
		return "attribute"
	case FieldKindBlock:
		return "block"
	default:
		return "unknown"
	}
}

// TerraformSchema holds parsed schema data from `terraform providers schema -json`.
type TerraformSchema struct {
	Resources map[string]*ResourceSchema
}

// ResourceSchema holds the field-kind lookup for one resource.
type ResourceSchema struct {
	Fields map[string]FieldKind // dot-separated path → kind
}

// FieldKind returns the kind for a dot-separated field path, or FieldKindUnknown.
func (rs *ResourceSchema) FieldKind(path string) FieldKind {
	if k, ok := rs.Fields[path]; ok {
		return k
	}
	return FieldKindUnknown
}

// rawProviderSchema mirrors the JSON structure from `terraform providers schema -json`.
type rawProviderSchema struct {
	FormatVersion   string                      `json:"format_version"`
	ProviderSchemas map[string]rawProviderEntry `json:"provider_schemas"`
}

type rawProviderEntry struct {
	ResourceSchemas map[string]rawResourceSchema `json:"resource_schemas"`
}

type rawResourceSchema struct {
	Version int      `json:"version"`
	Block   rawBlock `json:"block"`
}

type rawBlock struct {
	Attributes map[string]json.RawMessage `json:"attributes"`
	BlockTypes map[string]rawBlockType    `json:"block_types"`
}

type rawBlockType struct {
	NestingMode string   `json:"nesting_mode"`
	Block       rawBlock `json:"block"`
}

// ParseTerraformSchema parses the JSON output of `terraform providers schema -json`.
func ParseTerraformSchema(data []byte) (*TerraformSchema, error) {
	var raw rawProviderSchema
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse terraform schema: %w", err)
	}

	ts := &TerraformSchema{Resources: make(map[string]*ResourceSchema)}

	for _, provider := range raw.ProviderSchemas {
		for name, res := range provider.ResourceSchemas {
			rs := &ResourceSchema{Fields: make(map[string]FieldKind)}
			walkBlock(rs, "", res.Block)
			ts.Resources[name] = rs
		}
	}

	return ts, nil
}

// walkBlock recursively populates the field-kind map.
func walkBlock(rs *ResourceSchema, prefix string, block rawBlock) {
	for name := range block.Attributes {
		path := joinPath(prefix, name)
		rs.Fields[path] = FieldKindAttribute
	}
	for name, bt := range block.BlockTypes {
		path := joinPath(prefix, name)
		rs.Fields[path] = FieldKindBlock
		walkBlock(rs, path, bt.Block)
	}
}

func joinPath(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "." + name
}

// LoadTerraformSchema parses schema JSON bytes.
func LoadTerraformSchema(schemaJSON []byte) (*TerraformSchema, error) {
	return ParseTerraformSchema(schemaJSON)
}

// ListBlockFields returns all top-level fields of kind FieldKindBlock for a resource.
func (rs *ResourceSchema) ListBlockFields() []string {
	var blocks []string
	for path, kind := range rs.Fields {
		if kind == FieldKindBlock && !strings.Contains(path, ".") {
			blocks = append(blocks, path)
		}
	}
	sort.Strings(blocks)
	return blocks
}
