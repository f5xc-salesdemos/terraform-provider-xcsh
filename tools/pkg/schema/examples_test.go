// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package schema

import (
	"strings"
	"testing"

	"github.com/f5xc-salesdemos/terraform-provider-xcsh/tools/pkg/openapi"
)

func TestGenerateExampleUsage_Basic(t *testing.T) {
	attrs := []openapi.TerraformAttribute{
		{Name: "name", Type: "string", Required: true},
		{Name: "namespace", Type: "string", Required: true},
		{Name: "port", Type: "int64", Optional: true},
	}
	result := GenerateExampleUsage("http_loadbalancer", attrs)
	if !strings.Contains(result, `resource "xcsh_http_loadbalancer" "example"`) {
		t.Error("Should contain resource declaration")
	}
	if !strings.Contains(result, `name      = "example"`) {
		t.Error("Should contain name attribute")
	}
	if !strings.Contains(result, "namespace") {
		t.Error("Should contain namespace attribute")
	}
	if !strings.Contains(result, "}") {
		t.Error("Should contain closing brace")
	}
}

func TestGenerateExampleUsage_WithBlock(t *testing.T) {
	attrs := []openapi.TerraformAttribute{
		{Name: "name", Type: "string", Required: true},
		{Name: "spec", Type: "object", IsBlock: true, NestedAttributes: []openapi.TerraformAttribute{
			{Name: "host", Type: "string"},
			{Name: "port", Type: "int64"},
		}},
	}
	result := GenerateExampleUsage("origin_pool", attrs)
	if !strings.Contains(result, "spec {") {
		t.Error("Should contain block example")
	}
	if !strings.Contains(result, `host = "example"`) {
		t.Error("Should contain nested string example")
	}
}

func TestGenerateExampleUsage_SkipsComputed(t *testing.T) {
	attrs := []openapi.TerraformAttribute{
		{Name: "id", Type: "string", Computed: true},
		{Name: "custom", Type: "string", Optional: true},
	}
	result := GenerateExampleUsage("test_resource", attrs)
	if strings.Contains(result, "id =") {
		t.Error("Should not include computed attributes in example")
	}
}
