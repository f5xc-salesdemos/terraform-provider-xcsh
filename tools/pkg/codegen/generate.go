// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

// Package codegen (generate.go) provides file generation functions that
// execute Go text/templates against ResourceTemplate data and write the
// formatted output to disk.
package codegen

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/openapi"
	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/schema"
)

// GenerateResourceFile generates the Terraform resource Go file for a single resource.
// outputDir is the directory where the file will be written (e.g. "internal/provider").
func GenerateResourceFile(resource *openapi.ResourceTemplate, outputDir string) error {
	outputPath := filepath.Join(outputDir, resource.Name+"_resource.go")

	// Create template with custom functions
	funcMap := template.FuncMap{
		"renderNestedAttrs":               RenderNestedAttributes,
		"renderNestedBlocks":              RenderNestedBlocks,
		"renderNestedModelTypes":          RenderNestedModelTypes,
		"renderBlockFields":               RenderBlockFields,
		"renderSpecStructFields":          RenderSpecStructFields,
		"renderSpecMarshalCode":           RenderSpecMarshalCode,
		"renderSpecMarshalCodeForCreate":  RenderSpecMarshalCodeForCreate,
		"renderSpecUnmarshalCode":         RenderSpecUnmarshalCode,
		"renderCreateComputedFieldsCode":  RenderCreateComputedFieldsCode,
		"renderUpdateComputedFieldsCode":  RenderUpdateComputedFieldsCode,
		"renderFetchedComputedFieldsCode": RenderFetchedComputedFieldsCode,
		"filterSpecFields":                schema.FilterSpecFields,
		"enumValuesLiteral": func(values []string) string {
			quoted := make([]string, len(values))
			for i, v := range values {
				quoted[i] = fmt.Sprintf("%q", v)
			}
			return strings.Join(quoted, ", ")
		},
		"regexLiteral": RegexLiteral,
	}

	tmpl, err := template.New("resource").Funcs(funcMap).Parse(ResourceTemplate)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	// Execute template to buffer first
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, resource); err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}

	// Format the generated code with gofmt
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// If formatting fails, write unformatted code with warning
		fmt.Printf("Warning: gofmt failed for %s: %v (writing unformatted)\n", outputPath, err)
		formatted = buf.Bytes()
	}

	return os.WriteFile(outputPath, formatted, 0644)
}

// GenerateClientTypes generates the client type Go file for a single resource.
// clientDir is the directory where the file will be written (e.g. "internal/client").
func GenerateClientTypes(resource *openapi.ResourceTemplate, clientDir string) error {
	outputPath := filepath.Join(clientDir, resource.Name+"_types.go")

	// Create template with custom functions for spec field generation
	funcMap := template.FuncMap{
		"renderSpecStructFields": func(attrs []openapi.TerraformAttribute) string {
			return RenderSpecStructFields(attrs, "\t")
		},
	}

	tmpl, err := template.New("client").Funcs(funcMap).Parse(ClientTemplate)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	// Execute template to buffer first
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, resource); err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}

	// Format the generated code with gofmt
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// If formatting fails, write unformatted code with warning
		fmt.Printf("Warning: gofmt failed for %s: %v (writing unformatted)\n", outputPath, err)
		formatted = buf.Bytes()
	}

	return os.WriteFile(outputPath, formatted, 0644)
}

// GenerateDataSource generates the Terraform data source Go file for a single resource.
// outputDir is the directory where the file will be written (e.g. "internal/provider").
func GenerateDataSource(resource *openapi.ResourceTemplate, outputDir string) error {
	outputPath := filepath.Join(outputDir, resource.Name+"_data_source.go")

	tmpl, err := template.New("datasource").Parse(DataSourceTemplate)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	// Execute template to buffer first
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, resource); err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}

	// Format the generated code with gofmt
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// If formatting fails, write unformatted code with warning
		fmt.Printf("Warning: gofmt failed for %s: %v (writing unformatted)\n", outputPath, err)
		formatted = buf.Bytes()
	}

	return os.WriteFile(outputPath, formatted, 0644)
}
