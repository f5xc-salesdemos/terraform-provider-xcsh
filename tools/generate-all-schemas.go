// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

//go:build ignore
// +build ignore

// generate-all-schemas.go - Batch generator for all F5 XC Terraform resources
// This tool processes all OpenAPI spec files and generates comprehensive Terraform schemas.
//
// CI/CD Integration:
//   Changes to this file trigger the generate.yml workflow which regenerates all
//   provider resources from the latest OpenAPI specifications.
//
// Pipeline Verification: 2026-05-23T05:45Z - Verified against api-specs-enriched v2.1.104
//
// Usage: go run tools/generate-all-schemas.go [--spec-dir=/path/to/specs] [--dry-run]
//
// Environment Variables:
//   F5XC_SPEC_DIR - Directory containing OpenAPI spec files (default: /tmp)

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/codegen"
	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/namespace"
	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/naming"
	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/openapi"
	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/registration"
	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/schema"
)

// Configuration
var (
	specDir   string
	dryRun    bool
	outputDir string
	clientDir string
	verbose   bool
)

// Type aliases — canonical definitions live in tools/pkg/openapi.
type OpenAPI3Spec = openapi.Spec
type Components = openapi.Components
type SchemaDefinition = openapi.Schema
type TerraformAttribute = openapi.TerraformAttribute
type ResourceTemplate = openapi.ResourceTemplate
type GenerationResult = openapi.GenerationResult

// Global maps from index.json for resource metadata enrichment
var (
	resourceTierMap         = make(map[string]string)                        // resourceName -> tier
	resourceDependencyMap   = make(map[string]*openapi.ResourceDependencies) // resourceName -> dependencies
	resourceReferencedByMap = make(map[string][]string)                      // resourceName -> resources that depend on it
	resourceCategoryMap     = make(map[string]string)                        // resourceName -> category
)

// schemaCache and rawSpecCache are aliases for the canonical caches in the schema package.
// They remain here as local vars for backward compatibility with code that writes to them.
var schemaCache = schema.SchemaCache
var rawSpecCache = schema.RawSpecCache

func init() {
	flag.StringVar(&specDir, "spec-dir", "", "Directory containing OpenAPI spec files")
	flag.BoolVar(&dryRun, "dry-run", false, "Show what would be generated without writing files")
	flag.StringVar(&outputDir, "output-dir", "internal/provider", "Output directory for provider files")
	flag.StringVar(&clientDir, "client-dir", "internal/client", "Output directory for client files")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
}

func main() {
	flag.Parse()

	// Check for spec directory
	if specDir == "" {
		specDir = os.Getenv("F5XC_SPEC_DIR")
	}
	if specDir == "" {
		specDir = "docs/specifications/api"
	}

	fmt.Println("🔨 F5XC Terraform Provider - Batch Schema Generator")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("📁 Spec Directory: %s\n", specDir)
	fmt.Printf("📁 Output Directory: %s\n", outputDir)
	if dryRun {
		fmt.Println("🔍 DRY RUN MODE - No files will be written")
	}
	fmt.Println()

	// Detect spec version (expects v2 format)
	specVersion := openapi.GetSpecVersion(specDir)
	fmt.Printf("🔍 Detected spec version: %s\n\n", specVersion)

	var results []GenerationResult
	var successCount, failCount int

	switch specVersion {
	case openapi.SpecVersionV2:
		results, successCount, failCount = processV2Specs(specDir)
	default:
		fmt.Printf("❌ Unknown spec format in directory: %s\n", specDir)
		fmt.Println("💡 Expected v2 spec format: index.json + domains/*.json structure")
		os.Exit(1)
	}

	// Generate combined client types file
	if !dryRun {
		registration.GenerateCombinedClientTypes(results, clientDir)
	}

	// Print summary
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("📊 Generation Summary")
	fmt.Println(strings.Repeat("=", 60))
	readOnlyCount := 0
	for _, r := range results {
		if r.Success && r.IsReadOnly {
			readOnlyCount++
		}
	}
	fmt.Printf("✅ Successfully generated: %d resources + %d read-only data sources\n", successCount-readOnlyCount, readOnlyCount)
	fmt.Printf("⏭️  Skipped (no schema): %d\n", len(results)-successCount-failCount)
	fmt.Printf("❌ Failed: %d\n", failCount)

	if failCount > 0 {
		fmt.Println("\n❌ Failed resources:")
		for _, r := range results {
			if !r.Success && r.Error != "" {
				fmt.Printf("   - %s: %s\n", r.ResourceName, r.Error)
			}
		}
	}

	// Generate provider registration
	if !dryRun {
		registration.GenerateProviderRegistration(results, outputDir)
	}

	// Clean up orphan generated files that no longer have matching resources
	if !dryRun {
		generatedNames := make(map[string]bool)
		for _, r := range results {
			if r.Success {
				generatedNames[r.ResourceName] = true
			}
		}
		// Also include core resources that are always registered
		for _, core := range registration.CoreResources {
			generatedNames[core] = true
		}
		registration.CleanOrphanGeneratedFiles(outputDir, clientDir, generatedNames)
	}

	fmt.Println("\n🎉 Batch generation complete!")
}

// processV2Specs processes v2 format specs (domain-organized files from f5xc-api-enriched)
func processV2Specs(specDir string) ([]GenerationResult, int, int) {
	// Parse the index.json to get domain information
	index, err := openapi.ParseIndexFromDir(specDir)
	if err != nil {
		fmt.Printf("❌ Error parsing index.json: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📋 Spec version: %s\n", index.Version)
	fmt.Printf("📋 Generated at: %s\n", index.Timestamp)
	fmt.Printf("📄 Found %d domain specifications (v2 format)\n\n", len(index.Specifications))

	// Build global maps from index.json for metadata enrichment
	resourceTierMap = openapi.BuildResourceTierMap(index)
	resourceDependencyMap = openapi.BuildResourceDependencyMap(index)
	resourceReferencedByMap = openapi.BuildReferencedByMap(resourceDependencyMap)
	resourceCategoryMap = openapi.BuildResourceCategoryMap(index)
	if verbose {
		fmt.Printf("📊 Loaded metadata: %d resources with tier, %d with dependencies\n",
			len(resourceTierMap), len(resourceDependencyMap))
	}

	// Find all domain spec files
	domainFiles, err := openapi.FindDomainSpecFiles(specDir)
	if err != nil {
		fmt.Printf("❌ Error finding domain spec files: %v\n", err)
		os.Exit(1)
	}

	results := []GenerationResult{}
	successCount := 0
	failCount := 0
	skipCount := 0

	// Track processed resources to avoid duplicates across domain files
	// Some resources (like service_policy) appear in multiple domain specs
	processedResources := make(map[string]bool)

	// Build a map of domain metadata from index for quick lookup
	domainMetadata := make(map[string]openapi.DomainMetadata)
	for _, dm := range index.Specifications {
		domainMetadata[dm.Name] = dm
	}

	// Process each domain file
	for _, domainFile := range domainFiles {
		domainName := strings.TrimSuffix(filepath.Base(domainFile), ".json")
		fmt.Printf("🔄 Processing domain: %s\n", domainName)

		// Get domain metadata from index
		dm, hasMeta := domainMetadata[domainName]
		if hasMeta && verbose {
			fmt.Printf("   Category: %s, Tier: %s\n", dm.Category, dm.RequiresTier)
		}

		// Extract resources from the domain spec
		domainInfo, err := openapi.ExtractResourcesFromDomain(domainFile)
		if err != nil {
			fmt.Printf("   ⚠️  Error parsing domain: %v\n", err)
			results = append(results, GenerationResult{
				ResourceName: domainName,
				Success:      false,
				Error:        err.Error(),
			})
			failCount++
			continue
		}

		if len(domainInfo.Resources) == 0 {
			fmt.Printf("   ⏭️  No resources found in domain\n")
			continue
		}

		fmt.Printf("   📦 Found %d resources\n", len(domainInfo.Resources))

		// Process each resource in the domain
		for _, resource := range domainInfo.Resources {
			// Skip duplicate resources that appear in multiple domain specs
			if processedResources[resource.Name] {
				if verbose {
					fmt.Printf("      ⏭️  Skipping duplicate: %s (already processed)\n", resource.Name)
				}
				skipCount++
				continue
			}
			processedResources[resource.Name] = true

			// Create a virtual spec file path for compatibility with existing processing
			// The v2 domain spec contains all resources, so we process each individually
			result := processV2Resource(domainFile, resource, domainInfo)
			results = append(results, result)
			if result.Success {
				successCount++
			} else if result.Error != "" {
				failCount++
			}
		}
	}

	// Log duplicates if any were skipped
	if skipCount > 0 {
		fmt.Printf("\n⏭️  Skipped %d duplicate resources across domain files\n", skipCount)
	}

	return results, successCount, failCount
}

// processV2Resource processes a single resource from a v2 domain spec
func processV2Resource(domainFile string, resource openapi.ExtractedResource, domainInfo *openapi.DomainSpecInfo) GenerationResult {
	if verbose {
		fmt.Printf("      Processing resource: %s (category: %s, tier: %s)\n",
			resource.Name, resource.Category, resource.RequiresTier)
	}

	// If the spec declares x-f5xc-namespace-scope, record it so namespace.ForResource
	// returns the spec-derived value instead of the hardcoded default.
	if domainInfo.Spec != nil {
		// Check domain-level scope first
		scope := domainInfo.Spec.XF5XCNamespaceScope
		// Check info-level scope (overrides domain-level)
		if domainInfo.Spec.Info.XF5XCNamespaceScope != "" {
			scope = domainInfo.Spec.Info.XF5XCNamespaceScope
		}
		if scope != "" {
			namespace.SetSpecScope(resource.Name, scope)
			if verbose {
				fmt.Printf("      Namespace scope override: %s -> %s\n", resource.Name, scope)
			}
		}
	}

	// Use the existing processing with the domain file
	// The schema extraction will find the right Object schema based on resource name
	result := processSpecFileForResource(domainFile, resource.Name, resource.Category, resource.RequiresTier)

	return result
}

// processSpecFileForResource processes a spec file targeting a specific resource
// This is used for v2 specs where multiple resources exist in one domain file
func processSpecFileForResource(specFile string, resourceName string, category string, requiresTier string) GenerationResult {
	// Parse the spec file
	spec, err := parseOpenAPISpec(specFile)
	if err != nil {
		return GenerationResult{ResourceName: resourceName, Success: false, Error: err.Error()}
	}

	// Cache all schemas from the spec
	for name, schema := range spec.Components.Schemas {
		schemaCache[name] = schema
	}

	// Try to extract the resource schema
	schema, schemaName, isReadOnly := extractResourceSchemaByName(spec, resourceName)
	if schema == nil {
		return GenerationResult{ResourceName: resourceName, Success: false, Error: ""}
	}

	// Extract API path from spec (or construct from resource name)
	apiPath := extractAPIPathForResource(spec, resourceName)
	if apiPath == "" {
		apiPath = fmt.Sprintf("/api/config/namespaces/{namespace}/%ss", resourceName)
	}

	if isReadOnly {
		return generateReadOnlyDataSource(resourceName, schemaName, schema, apiPath, specFile, category, requiresTier)
	}
	return generateResourceFromSchema(resourceName, schemaName, schema, apiPath, specFile, category, requiresTier)
}

// extractResourceSchemaByName extracts a specific resource schema by name from a spec.
// Returns (schema, schemaName, isReadOnly). isReadOnly is true when only GetSpecType was found.
func extractResourceSchemaByName(spec *OpenAPI3Spec, resourceName string) (*SchemaDefinition, string, bool) {
	// Try CreateSpecType first (CRUD resources)
	createPatterns := []string{
		fmt.Sprintf("%sCreateSpecType", resourceName),
		fmt.Sprintf("schema%sCreateSpecType", resourceName),
		fmt.Sprintf("%sReplaceSpecType", resourceName),
		fmt.Sprintf("schema%sReplaceSpecType", resourceName),
	}
	for _, pattern := range createPatterns {
		if schema, ok := spec.Components.Schemas[pattern]; ok {
			return &schema, pattern, false
		}
	}

	// Legacy patterns (ves.io.schema format)
	legacyPatterns := []string{
		fmt.Sprintf("ves.io.schema.%s.Object", resourceName),
		fmt.Sprintf("%sType", naming.ToResourceTypeName(resourceName)),
		resourceName,
	}
	for _, pattern := range legacyPatterns {
		if schema, ok := spec.Components.Schemas[pattern]; ok {
			return &schema, pattern, false
		}
	}

	// Sort schema names for deterministic fallback matching
	sortedNames := make([]string, 0, len(spec.Components.Schemas))
	for name := range spec.Components.Schemas {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)

	lowerName := strings.ToLower(resourceName)

	// Fallback: case-insensitive match for legacy .Object suffix
	for _, name := range sortedNames {
		if strings.Contains(strings.ToLower(name), lowerName) && strings.HasSuffix(name, ".Object") {
			schema := spec.Components.Schemas[name]
			return &schema, name, false
		}
	}

	// Fallback: case-insensitive match for CreateSpecType
	for _, name := range sortedNames {
		if strings.Contains(strings.ToLower(name), lowerName) && strings.HasSuffix(name, "CreateSpecType") {
			schema := spec.Components.Schemas[name]
			return &schema, name, false
		}
	}

	// No CreateSpecType found — try GetSpecType (read-only resources)
	getPatterns := []string{
		fmt.Sprintf("%sGetSpecType", resourceName),
		fmt.Sprintf("schema%sGetSpecType", resourceName),
		fmt.Sprintf("views%sGetSpecType", resourceName),
	}
	for _, pattern := range getPatterns {
		if schema, ok := spec.Components.Schemas[pattern]; ok {
			return &schema, pattern, true
		}
	}

	// Fallback: case-insensitive match for GetSpecType
	for _, name := range sortedNames {
		if strings.Contains(strings.ToLower(name), lowerName) && strings.HasSuffix(name, "GetSpecType") {
			schema := spec.Components.Schemas[name]
			return &schema, name, true
		}
	}

	return nil, "", false
}

// extractAPIPathForResource extracts the API path for a specific resource from a spec
func extractAPIPathForResource(spec *OpenAPI3Spec, resourceName string) string {
	plural := resourceName + "s"
	// Handle special pluralization
	if strings.HasSuffix(resourceName, "y") {
		plural = strings.TrimSuffix(resourceName, "y") + "ies"
	}

	// Sort paths for deterministic matching (map order is random in Go)
	paths := make([]string, 0, len(spec.Paths))
	for path := range spec.Paths {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		if strings.Contains(path, "/"+plural) {
			return path
		}
	}
	return ""
}

// generateReadOnlyDataSource generates only a data source for a read-only resource (GetSpecType only).
func generateReadOnlyDataSource(resourceName string, schemaName string, schemaDef *SchemaDefinition, apiPath string, specFile string, category string, requiresTier string) GenerationResult {
	if verbose {
		fmt.Printf("      Generating read-only data source: %s (schema: %s)\n", resourceName, schemaName)
	}

	spec, err := parseOpenAPISpec(specFile)
	if err != nil {
		return GenerationResult{ResourceName: resourceName, Success: false, Error: err.Error()}
	}

	tmpl, err := schema.ExtractReadOnlyResourceSchema(spec, resourceName, extractAPIPath)
	if err != nil {
		if verbose {
			fmt.Printf("  ⏭️  Skipping read-only %s: %v\n", resourceName, err)
		}
		return GenerationResult{ResourceName: resourceName, Success: false}
	}

	if !dryRun {
		if err := codegen.GenerateReadOnlyDataSource(tmpl, outputDir); err != nil {
			return GenerationResult{ResourceName: resourceName, Success: false, Error: err.Error()}
		}
		if err := codegen.GenerateReadOnlyClientTypes(tmpl, clientDir); err != nil {
			return GenerationResult{ResourceName: resourceName, Success: false, Error: err.Error()}
		}
	}

	return GenerationResult{
		ResourceName: resourceName,
		Success:      true,
		IsReadOnly:   true,
		AttrCount:    len(tmpl.Attributes),
	}
}

// generateResourceFromSchema generates a resource file from an extracted schema
// This is factored out from processSpecFile to support both v1 and v2 processing
func generateResourceFromSchema(resourceName string, schemaName string, schemaDef *SchemaDefinition, apiPath string, specFile string, category string, requiresTier string) GenerationResult {
	if verbose {
		fmt.Printf("      Generating: %s (schema: %s)\n", resourceName, schemaName)
		if category != "" {
			fmt.Printf("      Category: %s, Tier: %s\n", category, requiresTier)
		}
	}

	// Skip internal/utility schemas
	skipPatterns := []string{
		"object", "status", "spec", "metadata", "types", "common",
		"refs", "crudapi", "public", "private", "api", "empty",
	}
	for _, skip := range skipPatterns {
		if resourceName == skip {
			return GenerationResult{ResourceName: resourceName, Success: false}
		}
	}

	// Parse spec to get full schema information
	spec, err := parseOpenAPISpec(specFile)
	if err != nil {
		return GenerationResult{ResourceName: resourceName, Success: false, Error: err.Error()}
	}

	// Extract resource schema using the resource name we have
	resource, err := schema.ExtractResourceSchema(spec, resourceName, extractAPIPath)
	if err != nil {
		if verbose {
			fmt.Printf("  ⏭️  Skipping %s: %v\n", resourceName, err)
		}
		return GenerationResult{ResourceName: resourceName, Success: false}
	}

	// Count attributes and blocks
	attrCount := 0
	blockCount := 0
	for _, attr := range resource.Attributes {
		if attr.IsBlock {
			blockCount++
		} else {
			attrCount++
		}
	}

	if !dryRun {
		// Generate resource file
		if err := codegen.GenerateResourceFile(resource, outputDir); err != nil {
			return GenerationResult{ResourceName: resourceName, Success: false, Error: err.Error()}
		}

		// Generate client types
		if err := codegen.GenerateClientTypes(resource, clientDir); err != nil {
			return GenerationResult{ResourceName: resourceName, Success: false, Error: err.Error()}
		}

		// Generate data source
		if err := codegen.GenerateDataSource(resource, outputDir); err != nil {
			return GenerationResult{ResourceName: resourceName, Success: false, Error: err.Error()}
		}

	}

	fmt.Printf("✅ %s: %d attrs, %d blocks\n", resourceName, attrCount, blockCount)
	return GenerationResult{
		ResourceName: resourceName,
		Success:      true,
		AttrCount:    attrCount,
		BlockCount:   blockCount,
	}
}

func parseOpenAPISpec(specFile string) (*OpenAPI3Spec, error) {
	data, err := os.ReadFile(specFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec: %w", err)
	}

	var spec OpenAPI3Spec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse spec: %w", err)
	}

	// Cache schemas
	for name, schema := range spec.Components.Schemas {
		schemaCache[name] = schema
	}

	// Also parse raw JSON to extract x-ves-oneof-field annotations
	var rawSpec map[string]interface{}
	if err := json.Unmarshal(data, &rawSpec); err == nil {
		if components, ok := rawSpec["components"].(map[string]interface{}); ok {
			if schemas, ok := components["schemas"].(map[string]interface{}); ok {
				for name, schema := range schemas {
					if schemaMap, ok := schema.(map[string]interface{}); ok {
						rawSpecCache[name] = schemaMap
					}
				}
			}
		}
	}

	return &spec, nil
}

// extractAPIPath extracts the correct API path for CRUD operations from the OpenAPI spec
// It looks for paths containing POST (create) and returns the base path pattern
// Returns: basePath (for create/list), itemPath (for get/update/delete), hasNamespace (whether path has {namespace} segment)
func extractAPIPath(spec *OpenAPI3Spec, resourceName string) (basePath string, itemPath string, hasNamespace bool) {
	resourcePlural := resourceName + "s"

	// Sort path keys for deterministic matching (map order is random in Go)
	sortedPaths := make([]string, 0, len(spec.Paths))
	for path := range spec.Paths {
		sortedPaths = append(sortedPaths, path)
	}
	sort.Strings(sortedPaths)

	// Look for CRUD paths in the spec
	for _, path := range sortedPaths {
		pathObj := spec.Paths[path]
		pathMap, ok := pathObj.(map[string]interface{})
		if !ok {
			continue
		}

		// Check if this is a CRUD endpoint (has POST for create or GET for list)
		_, hasPost := pathMap["post"]
		_, hasGet := pathMap["get"]

		// Look for the base path (list/create endpoint) - ends with plural resource name
		// Pattern: /api/.../resource_names or /api/.../resource_names (with namespace)
		if (hasPost || hasGet) && strings.HasSuffix(path, "/"+resourcePlural) {
			// Check if path contains namespace segment
			hasNamespace = strings.Contains(path, "{namespace}") || strings.Contains(path, "{metadata.namespace}")

			// Convert {metadata.namespace} to %s for namespace substitution
			// Convert {metadata.name} or {name} for item paths
			basePath = path
			if hasNamespace {
				basePath = strings.ReplaceAll(basePath, "{metadata.namespace}", "%s")
				basePath = strings.ReplaceAll(basePath, "{namespace}", "%s")
			}

			// Item path is base path + /{name}
			itemPath = path + "/{name}"
			itemPath = strings.ReplaceAll(itemPath, "{metadata.namespace}", "%s")
			itemPath = strings.ReplaceAll(itemPath, "{namespace}", "%s")
			itemPath = strings.ReplaceAll(itemPath, "{metadata.name}", "%s")
			itemPath = strings.ReplaceAll(itemPath, "{name}", "%s")

			return basePath, itemPath, hasNamespace
		}
	}

	// Fallback to default pattern if no path found
	return fmt.Sprintf("/api/config/namespaces/%%s/%s", resourcePlural),
		fmt.Sprintf("/api/config/namespaces/%%s/%s/%%s", resourcePlural),
		true
}


// extractResourceSchema, hasNestedModelsWithAttrTypes, hasMaxLengthValidatorsAny,
// hasEnumValidatorsAny, hasPatternValidatorsAny, hasListSizeValidatorsAny,
// collectConflictAttrs, scanPlanModifierUsage, generateExampleUsage,
// convertToTerraformAttribute, convertToTerraformAttributeWithDepth,
// extractNestedAttributes, resolveRef, mapSchemaType, mapSchemaTypeToGo,
// filterOptional, parseMinConfigRequiredFields, promoteMinConfigRequired,
// extractOneOfGroups, isMetadataField, filterSpecFields, and maxNestedDepth
// have been moved to tools/pkg/schema/ (SP-2 Task 3).

// Code rendering functions (renderSpecStructFields, getGoClientType, renderSpecMarshalCode*,
// renderComputedFieldsCode, renderSpecUnmarshalCode, renderNestedAttributes, renderNestedBlocks,
// renderNestedModelTypes, renderBlockFields, collectNestedModelTypes, escapeGoString, regexLiteral,
// and the NestedModelInfo type) have been moved to tools/pkg/codegen/ (SP-2 Task 4).


// Provider registration (generateProviderRegistration, cleanOrphanGeneratedFiles,
// generateCombinedClientTypes, and CoreResources) have been moved to
// tools/pkg/registration/ (SP-2 Task 6).
