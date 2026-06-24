// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package codegen

import (
	"strings"
	"testing"

	"github.com/f5xc-salesdemos/terraform-provider-xcsh/tools/pkg/openapi"
)

func TestEscapeGoString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no special characters",
			input:    "simple string",
			expected: "simple string",
		},
		{
			name:     "with double quotes",
			input:    `say "hello"`,
			expected: `say \"hello\"`,
		},
		{
			name:     "with backslash",
			input:    `path\to\file`,
			expected: `path\\to\\file`,
		},
		{
			name:     "with both",
			input:    `path\to\"file"`,
			expected: `path\\to\\\"file\"`,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EscapeGoString(tt.input)
			if got != tt.expected {
				t.Errorf("EscapeGoString(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestRegexLiteral(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple pattern uses backticks",
			input:    `^[a-z]+$`,
			expected: "`^[a-z]+$`",
		},
		{
			name:     "pattern with backslash uses backticks",
			input:    `^\d{3}-\d{4}$`,
			expected: "`" + `^\d{3}-\d{4}$` + "`",
		},
		{
			name:     "pattern with backtick uses quoted string",
			input:    "pattern`with`backtick",
			expected: `"pattern` + "`" + `with` + "`" + `backtick"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RegexLiteral(tt.input)
			if got != tt.expected {
				t.Errorf("RegexLiteral(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGetGoClientType(t *testing.T) {
	tests := []struct {
		name     string
		attr     openapi.TerraformAttribute
		expected string
	}{
		{
			name:     "string type",
			attr:     openapi.TerraformAttribute{Type: "string"},
			expected: "string",
		},
		{
			name:     "int64 type",
			attr:     openapi.TerraformAttribute{Type: "int64"},
			expected: "int64",
		},
		{
			name:     "bool type",
			attr:     openapi.TerraformAttribute{Type: "bool"},
			expected: "bool",
		},
		{
			name:     "list of strings",
			attr:     openapi.TerraformAttribute{Type: "list", ElementType: "string"},
			expected: "[]string",
		},
		{
			name:     "list of int64",
			attr:     openapi.TerraformAttribute{Type: "list", ElementType: "int64"},
			expected: "[]int64",
		},
		{
			name:     "list of unknown",
			attr:     openapi.TerraformAttribute{Type: "list", ElementType: "object"},
			expected: "[]interface{}",
		},
		{
			name:     "map type",
			attr:     openapi.TerraformAttribute{Type: "map"},
			expected: "map[string]string",
		},
		{
			name:     "block single",
			attr:     openapi.TerraformAttribute{IsBlock: true, NestedBlockType: "single"},
			expected: "map[string]interface{}",
		},
		{
			name:     "block list",
			attr:     openapi.TerraformAttribute{IsBlock: true, NestedBlockType: "list"},
			expected: "[]map[string]interface{}",
		},
		{
			name:     "unknown type",
			attr:     openapi.TerraformAttribute{Type: "unknown"},
			expected: "interface{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetGoClientType(tt.attr)
			if got != tt.expected {
				t.Errorf("GetGoClientType() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestRenderSpecStructFields(t *testing.T) {
	tests := []struct {
		name     string
		attrs    []openapi.TerraformAttribute
		indent   string
		contains []string
		empty    bool
	}{
		{
			name:  "empty attrs",
			attrs: nil,
			empty: true,
		},
		{
			name: "metadata fields are excluded",
			attrs: []openapi.TerraformAttribute{
				{GoName: "Name", TfsdkTag: "name", JsonName: "name", Type: "string"},
			},
			empty: true,
		},
		{
			name: "string field with omitempty",
			attrs: []openapi.TerraformAttribute{
				{GoName: "Domain", TfsdkTag: "domain", JsonName: "domain", Type: "string", IsSpecField: true},
			},
			indent: "\t",
			contains: []string{
				`Domain string ` + "`" + `json:"domain,omitempty"` + "`",
			},
		},
		{
			name: "block field without omitempty",
			attrs: []openapi.TerraformAttribute{
				{GoName: "Config", TfsdkTag: "config", JsonName: "config", Type: "object", IsBlock: true, NestedBlockType: "single", IsSpecField: true},
			},
			indent: "\t",
			contains: []string{
				`Config map[string]interface{} ` + "`" + `json:"config"` + "`",
			},
		},
		{
			name: "uses tfsdk tag when json name empty",
			attrs: []openapi.TerraformAttribute{
				{GoName: "Port", TfsdkTag: "port", JsonName: "", Type: "int64", IsSpecField: true},
			},
			indent: "\t",
			contains: []string{
				`Port int64 ` + "`" + `json:"port,omitempty"` + "`",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderSpecStructFields(tt.attrs, tt.indent)
			if tt.empty {
				if got != "" {
					t.Errorf("expected empty string, got %q", got)
				}
				return
			}
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("RenderSpecStructFields() missing %q in:\n%s", want, got)
				}
			}
		})
	}
}

func TestRenderNestedAttributes_Empty(t *testing.T) {
	got := RenderNestedAttributes(nil, "\t")
	if got != "" {
		t.Errorf("expected empty string for nil attrs, got %q", got)
	}
}

func TestRenderNestedBlocks_NoBlocks(t *testing.T) {
	attrs := []openapi.TerraformAttribute{
		{GoName: "Name", TfsdkTag: "name", Type: "string"},
	}
	got := RenderNestedBlocks(attrs, "\t")
	if got != "" {
		t.Errorf("expected empty string when no blocks, got %q", got)
	}
}

func TestCollectNestedModelTypes(t *testing.T) {
	attrs := []openapi.TerraformAttribute{
		{
			GoName:          "Config",
			TfsdkTag:        "config",
			IsBlock:         true,
			NestedBlockType: "single",
			NestedAttributes: []openapi.TerraformAttribute{
				{GoName: "Port", TfsdkTag: "port", Type: "int64"},
			},
		},
	}

	var models []NestedModelInfo
	CollectNestedModelTypes("Test", attrs, "", &models)

	if len(models) != 1 {
		t.Fatalf("expected 1 model, got %d", len(models))
	}
	if models[0].TypeName != "TestConfigModel" {
		t.Errorf("expected type name TestConfigModel, got %s", models[0].TypeName)
	}
	if models[0].IsEmpty {
		t.Error("expected non-empty model")
	}
}
