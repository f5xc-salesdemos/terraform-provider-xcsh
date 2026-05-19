// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

// Package registration generates the provider.go file that registers all
// resources and data sources, and cleans up orphan generated files that no
// longer correspond to any OpenAPI specification.
package registration

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/naming"
	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/openapi"
)

// CoreResources are resources that must always be registered in the provider,
// even if they're not present in the current OpenAPI specifications.
// These resources have working implementations that were generated previously
// but may not be included in the enriched API specs.
// Note: namespace was removed in v3.0.0 as part of backwards compatibility cleanup
var CoreResources = []string{}

// GenerateCombinedClientTypes is a no-op placeholder retained for interface
// compatibility. Individual client type files are generated per-resource
// during the main generation loop.
func GenerateCombinedClientTypes(results []openapi.GenerationResult, clientDir string) {
	// This is handled by individual client type files
}

// CleanOrphanGeneratedFiles removes generated files that no longer have
// matching resources.  Only removes files with the "DO NOT EDIT" header to
// avoid deleting manually maintained files.
func CleanOrphanGeneratedFiles(outDir, clntDir string, generatedNames map[string]bool) {
	suffixes := []struct {
		dir    string
		suffix string
	}{
		{outDir, "_resource.go"},
		{outDir, "_data_source.go"},
		{clntDir, "_types.go"},
	}

	removedCount := 0
	for _, s := range suffixes {
		matches, err := filepath.Glob(filepath.Join(s.dir, "*"+s.suffix))
		if err != nil {
			continue
		}
		for _, file := range matches {
			baseName := strings.TrimSuffix(filepath.Base(file), s.suffix)
			if generatedNames[baseName] {
				continue
			}
			// Check if file has "DO NOT EDIT" header (generated file)
			content, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			headerLen := 500
			if len(content) < headerLen {
				headerLen = len(content)
			}
			if !strings.Contains(string(content[:headerLen]), "DO NOT EDIT") {
				continue
			}
			fmt.Printf("🗑️  Removing orphan: %s\n", filepath.Base(file))
			os.Remove(file)
			removedCount++
		}
	}

	// Also clean up orphan test files whose resource implementation no longer exists.
	// Test files don't have the "DO NOT EDIT" header, but if the corresponding
	// _resource.go file does not exist, the test file is dead code.
	testMatches, err := filepath.Glob(filepath.Join(outDir, "*_resource_test.go"))
	if err == nil {
		for _, testFile := range testMatches {
			baseName := strings.TrimSuffix(filepath.Base(testFile), "_resource_test.go")
			resourceFile := filepath.Join(outDir, baseName+"_resource.go")
			if _, err := os.Stat(resourceFile); os.IsNotExist(err) {
				fmt.Printf("🗑️  Removing orphan test: %s\n", filepath.Base(testFile))
				os.Remove(testFile)
				removedCount++
			}
		}
	}

	// Clean up orphan data source test files
	dsTestMatches, err := filepath.Glob(filepath.Join(outDir, "*_data_source_test.go"))
	if err == nil {
		for _, testFile := range dsTestMatches {
			baseName := strings.TrimSuffix(filepath.Base(testFile), "_data_source_test.go")
			dataSourceFile := filepath.Join(outDir, baseName+"_data_source.go")
			if _, err := os.Stat(dataSourceFile); os.IsNotExist(err) {
				fmt.Printf("🗑️  Removing orphan test: %s\n", filepath.Base(testFile))
				os.Remove(testFile)
				removedCount++
			}
		}
	}

	if removedCount > 0 {
		fmt.Printf("🧹 Cleaned up %d orphan generated files\n", removedCount)
	}
}

// GenerateProviderRegistration generates the provider.go file that wires all
// successfully generated resources and data sources into the Terraform
// provider.  outputDir is the directory for provider files (e.g.
// internal/provider).
func GenerateProviderRegistration(results []openapi.GenerationResult, outputDir string) {
	// Collect successful resources
	var resources []string
	var dataSources []string

	// First, add core resources that must always be registered
	// These have working implementations but may not be in current OpenAPI specs
	added := make(map[string]bool)
	for _, core := range CoreResources {
		titleCase := naming.ToResourceTypeName(core)
		resources = append(resources, fmt.Sprintf("\t\tNew%sResource,", titleCase))
		dataSources = append(dataSources, fmt.Sprintf("\t\tNew%sDataSource,", titleCase))
		added[core] = true
	}

	// Then add resources from spec generation results (avoiding duplicates)
	for _, r := range results {
		if r.Success && !added[r.ResourceName] {
			titleCase := naming.ToResourceTypeName(r.ResourceName)
			resources = append(resources, fmt.Sprintf("\t\tNew%sResource,", titleCase))
			dataSources = append(dataSources, fmt.Sprintf("\t\tNew%sDataSource,", titleCase))
		}
	}

	// Sort for consistent output
	sort.Strings(resources)
	sort.Strings(dataSources)

	// Generate provider.go file
	providerPath := filepath.Join(outputDir, "provider.go")
	fmt.Printf("\n📝 Updating provider.go with %d resources and %d data sources...\n", len(resources), len(dataSources))

	providerContent := fmt.Sprintf(`// Code generated by generate-all-schemas.go. DO NOT EDIT.
// Source: F5 XC OpenAPI specification

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/f5xc-salesdemos/terraform-provider-f5xc/internal/client"
)

// Ensure F5XCProvider satisfies various provider interfaces.
var _ provider.Provider = &F5XCProvider{}

// F5XCProvider defines the provider implementation.
type F5XCProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// F5XCProviderModel describes the provider data model.
type F5XCProviderModel struct {
	APIToken     types.String `+"`"+`tfsdk:"api_token"`+"`"+`
	APIURL       types.String `+"`"+`tfsdk:"api_url"`+"`"+`
	APIP12File   types.String `+"`"+`tfsdk:"api_p12_file"`+"`"+`
	P12Password  types.String `+"`"+`tfsdk:"p12_password"`+"`"+`
	APICert      types.String `+"`"+`tfsdk:"api_cert"`+"`"+`
	APIKey       types.String `+"`"+`tfsdk:"api_key"`+"`"+`
	APICACert    types.String `+"`"+`tfsdk:"api_ca_cert"`+"`"+`
}

func (p *F5XCProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "f5xc"
	resp.Version = p.version
}

func (p *F5XCProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform provider for F5 Distributed Cloud (F5XC) enabling infrastructure as code " +
			"for load balancers, security policies, sites, and networking. Community-maintained provider " +
			"built from public F5 API documentation.",
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				MarkdownDescription: "F5 Distributed Cloud API URL. " +
					"Defaults to https://console.ves.volterra.io. " +
					"Example: https://tenant.console.ves.volterra.io. " +
					"Can also be set via F5XC_API_URL environment variable.",
				Optional: true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "F5 Distributed Cloud API Token for token-based authentication. " +
					"Can also be set via F5XC_API_TOKEN environment variable. " +
					"Either api_token or api_p12_file/api_cert must be specified.",
				Optional:  true,
				Sensitive: true,
			},
			"api_p12_file": schema.StringAttribute{
				MarkdownDescription: "Path to PKCS#12 certificate bundle file for certificate-based authentication. " +
					"Can also be set via F5XC_P12_FILE environment variable. " +
					"When using P12 authentication, p12_password must also be provided.",
				Optional:  true,
				Sensitive: false,
			},
			"p12_password": schema.StringAttribute{
				MarkdownDescription: "Password for the PKCS#12 certificate bundle. " +
					"Can also be set via F5XC_P12_PASSWORD environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"api_cert": schema.StringAttribute{
				MarkdownDescription: "Path to PEM-encoded client certificate file for certificate-based authentication. " +
					"Can also be set via F5XC_CERT environment variable. " +
					"When using certificate authentication, api_key must also be provided.",
				Optional: true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Path to PEM-encoded client private key file for certificate-based authentication. " +
					"Can also be set via F5XC_KEY environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"api_ca_cert": schema.StringAttribute{
				MarkdownDescription: "Path to PEM-encoded CA certificate file for verifying the F5XC API server. " +
					"Can also be set via F5XC_CACERT environment variable. Optional.",
				Optional: true,
			},
		},
	}
}

func (p *F5XCProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring F5XC client")

	var config F5XCProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get configuration values from environment variables first
	apiURL := os.Getenv("F5XC_API_URL")
	apiToken := os.Getenv("F5XC_API_TOKEN")
	apiP12File := os.Getenv("F5XC_P12_FILE")
	p12Password := os.Getenv("F5XC_P12_PASSWORD")
	apiCert := os.Getenv("F5XC_CERT")
	apiKey := os.Getenv("F5XC_KEY")
	apiCACert := os.Getenv("F5XC_CACERT")

	// Configuration values override environment variables
	if !config.APIURL.IsNull() {
		apiURL = config.APIURL.ValueString()
	}
	if !config.APIToken.IsNull() {
		apiToken = config.APIToken.ValueString()
	}
	if !config.APIP12File.IsNull() {
		apiP12File = config.APIP12File.ValueString()
	}
	if !config.P12Password.IsNull() {
		p12Password = config.P12Password.ValueString()
	}
	if !config.APICert.IsNull() {
		apiCert = config.APICert.ValueString()
	}
	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}
	if !config.APICACert.IsNull() {
		apiCACert = config.APICACert.ValueString()
	}

	// Set default API URL if not provided
	if apiURL == "" {
		apiURL = "https://console.ves.volterra.io"
	}

	// Normalize the API URL (removes /api suffix and trailing slashes)
	apiURL, _ = normalizeAPIURL(apiURL)

	var c *client.Client
	var err error

	// Determine authentication method
	switch {
	case apiP12File != "":
		// P12 certificate authentication
		if p12Password == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("p12_password"),
				"Missing P12 Password",
				"When using P12 certificate authentication (api_p12_file), the p12_password must be provided. "+
					"Set the p12_password value in the configuration or use the F5XC_P12_PASSWORD environment variable.",
			)
			return
		}
		c, err = client.NewClientWithP12(apiURL, apiP12File, p12Password)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Create F5XC Client",
				"Could not create F5XC client with P12 certificate: "+err.Error(),
			)
			return
		}
		tflog.Info(ctx, "Configured F5XC client with P12 certificate authentication", map[string]any{"success": true, "api_url": apiURL})

	case apiCert != "" && apiKey != "":
		// PEM certificate/key authentication
		c, err = client.NewClientWithCert(apiURL, apiCert, apiKey, apiCACert)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Create F5XC Client",
				"Could not create F5XC client with certificate: "+err.Error(),
			)
			return
		}
		tflog.Info(ctx, "Configured F5XC client with certificate authentication", map[string]any{"success": true, "api_url": apiURL})

	case apiToken != "":
		// API token authentication
		c = client.NewClient(apiURL, apiToken)
		tflog.Info(ctx, "Configured F5XC client with API token authentication", map[string]any{"success": true, "api_url": apiURL})

	default:
		resp.Diagnostics.AddError(
			"Missing Authentication Configuration",
			"The provider requires authentication. Please configure one of the following:\n"+
				"  - api_token (or F5XC_API_TOKEN environment variable) for API token authentication\n"+
				"  - api_p12_file and p12_password (or F5XC_P12_FILE and F5XC_P12_PASSWORD environment variables) for P12 certificate authentication\n"+
				"  - api_cert and api_key (or F5XC_CERT and F5XC_KEY environment variables) for PEM certificate authentication",
		)
		return
	}

	// Make the client available during DataSource and Resource type Configure methods
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *F5XCProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
%s
	}
}

func (p *F5XCProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
%s
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &F5XCProvider{
			version: version,
		}
	}
}
`, strings.Join(resources, "\n"), strings.Join(dataSources, "\n"))

	// Format the generated code with gofmt
	formatted, err := format.Source([]byte(providerContent))
	if err != nil {
		// If formatting fails, write unformatted code with warning
		fmt.Printf("⚠️  gofmt failed for %s: %v (writing unformatted)\n", providerPath, err)
		formatted = []byte(providerContent)
	}

	if err := os.WriteFile(providerPath, formatted, 0644); err != nil {
		fmt.Printf("❌ Error writing provider.go: %v\n", err)
		return
	}

	fmt.Printf("✅ Updated %s\n", providerPath)
}
