// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

//go:build ignore
// +build ignore

// This tool generates a three-level llms.txt hierarchy for AI consumption:
//   - L0: docs/llms.txt — Provider entry point with category index
//   - L1: docs/_llms-txt/<subcategory>.txt — Category indexes with resource lists
//   - L2: docs/_llms-txt/resources/<name>.txt — Self-contained per-resource docs
//
// Usage:
//
//	go run tools/generate-llms-txt.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/resource"
)

// LLMsConfig represents the llms-config.json structure.
type LLMsConfig struct {
	Deprecation struct {
		Notice    string `json:"notice"`
		Canonical struct {
			Provider string `json:"provider"`
			Registry string `json:"registry"`
			Source   string `json:"source"`
		} `json:"canonical"`
		Deprecated struct {
			Provider string `json:"provider"`
			Registry string `json:"registry"`
			Source   string `json:"source"`
		} `json:"deprecated"`
		Note string `json:"note"`
	} `json:"deprecation"`
	CustomSets []struct {
		Label       string   `json:"label"`
		Description string   `json:"description"`
		Paths       []string `json:"paths"`
	} `json:"customSets"`
	Promote []string `json:"promote"`
	Demote  []string `json:"demote"`
}

// ResourceInfo holds parsed information about a resource.
type ResourceInfo struct {
	Name            string
	FullName        string // f5xc_<name>
	Category        string
	Description     string
	RequiredFields  []string
	OneOfGroups     []OneOfGroup
	ServerDefaults  []string
	MinimalConfig   string
	Dependencies    []string
	IsPromoted      bool
}

// OneOfGroup represents a mutually exclusive group of fields.
type OneOfGroup struct {
	Parent  string   // parent block name (empty for root level)
	Options []string // list of option names
}

// CategoryInfo holds category metadata.
type CategoryInfo struct {
	Name        string
	Slug        string
	Description string
	Resources   []ResourceInfo
}

// JSONIndex is the top-level structure for terraform-llms-index.json.
type JSONIndex struct {
	Version    string                  `json:"version"`
	Provider   JSONProvider            `json:"provider"`
	Categories []JSONCategory          `json:"categories"`
	Resources  map[string]JSONResource `json:"resources"`
}

type JSONProvider struct {
	Source        string   `json:"source"`
	Registry      string   `json:"registry"`
	RequiredBlock string   `json:"required_block"`
	SyntaxRules   []string `json:"syntax_rules"`
}

type JSONCategory struct {
	Name            string   `json:"name"`
	Slug            string   `json:"slug"`
	Description     string   `json:"description"`
	ResourceCount   int      `json:"resource_count"`
	Resources       []string `json:"resources"`
	DependencyChain string   `json:"dependency_chain,omitempty"`
}

type JSONOneOfGroup struct {
	Parent string   `json:"parent,omitempty"`
	Fields []string `json:"fields"`
}

type JSONDependencies struct {
	Requires []string `json:"requires"`
	UsedBy   []string `json:"used_by,omitempty"`
}

type JSONResource struct {
	Category       string           `json:"category"`
	Description    string           `json:"description"`
	Required       []string         `json:"required"`
	OneOfGroups    []JSONOneOfGroup `json:"oneof_groups,omitempty"`
	ServerDefaults []string         `json:"server_defaults,omitempty"`
	MinimalConfig  string           `json:"minimal_config,omitempty"`
	Dependencies   JSONDependencies `json:"dependencies"`
	ImportSyntax   string           `json:"import_syntax"`
}

// categoryDescriptions provides hardcoded descriptions for each category.
var categoryDescriptions = map[string]string{
	"API Security":              "API definition, discovery, testing, and security controls for web APIs",
	"Applications":              "Application settings, types, discovery, and filtering",
	"Authentication":            "Authentication methods, cloud credentials, and secret management",
	"BIG-IP Integration":        "BIG-IP proxy, data groups, and iRules integration",
	"Certificates":              "TLS certificates, certificate chains, CRLs, and trusted CA lists",
	"Cloud Resources":           "Cloud elastic IPs, address allocators, and geo-location resources",
	"DNS":                       "DNS domains, zones, compliance checks, and DNS proxy configuration",
	"Integrations":              "External integrations including code base and ticket tracking",
	"Kubernetes":                "Container registries, workloads, and Kubernetes integrations",
	"Load Balancing":            "HTTP/TCP/UDP/CDN load balancers, origin pools, health checks, and routing",
	"Monitoring":                "Log receivers, alert policies, APM, and global logging configuration",
	"Networking":                "Virtual networks, BGP, cloud connectivity, tunnels, and network interfaces",
	"Organization":              "Namespaces, tenant configuration, and organizational settings",
	"Security":                  "WAF, bot defense, rate limiting, firewall policies, and security controls",
	"Service Mesh":              "Service mesh policies and traffic management",
	"Sites":                     "AWS/Azure/GCP VPC sites, SecureMesh, VoltStack, and site mesh groups",
	"Subscriptions":             "Cloud subscription management and metering",
	"Uncategorized":             "Resources pending categorization",
	"VPN":                       "VPN and IPSec configuration",
	"Infrastructure Protection": "DDoS protection and infrastructure security",
}

// Regex patterns
var (
	oneOfCommentRE     = regexp.MustCompile(`//\s*One of the arguments from this list "([^"]+)" must be set`)
	codeBlockRE        = regexp.MustCompile("(?s)```(?:terraform|hcl)\n(.*?)```")
	serverDefaultLineRE = regexp.MustCompile(`^#\s+-\s+(\S+)\s*$`)
	requiredFieldRE    = regexp.MustCompile(`<a[^>]*>.*?</a>.*?\[(` + "`" + `[^` + "`" + `]+` + "`" + `)\].*?-\s*Required\s+(String|Number|Bool|Block)`)
	serverDefaultTagRE = regexp.MustCompile(`Server Default|⚙️\s*\*\*Server Default\*\*`)
	resourceBlockRE    = regexp.MustCompile(`(?s)resource\s+"f5xc_\w+"\s+"[^"]+"\s+\{[^}]*\}`)
)

func main() {
	config, err := loadConfig("docs/llms-config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Parse all resources
	resources, err := parseAllResources(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing resources: %v\n", err)
		os.Exit(1)
	}

	// Group by category
	categories := groupByCategory(resources)

	// Build reverse dependency map
	reverseDeps := buildReverseDependencies(resources)

	// Create output directories
	if err := os.MkdirAll("docs/_llms-txt/resources", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directories: %v\n", err)
		os.Exit(1)
	}

	// Generate L2 files (per-resource)
	for _, res := range resources {
		if err := generateL2(res, reverseDeps); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating L2 for %s: %v\n", res.Name, err)
		}
	}
	fmt.Printf("Generated %d L2 resource files\n", len(resources))

	// Generate L1 files (per-category)
	for _, cat := range categories {
		if err := generateL1(cat); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating L1 for %s: %v\n", cat.Name, err)
		}
	}
	fmt.Printf("Generated %d L1 category files\n", len(categories))

	// Generate L0 file (entry point)
	if err := generateL0(config, categories); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating L0: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Generated L0 entry point: docs/llms.txt")

	// Generate JSON index for machine consumption
	if err := generateJSONIndex(config, categories, reverseDeps); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating JSON index: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Generated JSON index: docs/terraform-llms-index.json")

	// Summary
	totalFiles := 2 + len(categories) + len(resources)
	fmt.Printf("\nTotal: %d files generated (1 L0 + %d L1 + %d L2 + 1 JSON)\n",
		totalFiles, len(categories), len(resources))
}

func loadConfig(path string) (*LLMsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config LLMsConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func parseAllResources(config *LLMsConfig) ([]ResourceInfo, error) {
	files, err := filepath.Glob("docs/resources/*.md")
	if err != nil {
		return nil, err
	}

	var resources []ResourceInfo
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, ".md")
		if name == "index" {
			continue
		}

		res, err := parseResource(f, name, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not parse %s: %v\n", name, err)
			continue
		}
		resources = append(resources, res)
	}

	return resources, nil
}

func parseResource(path, name string, config *LLMsConfig) (ResourceInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ResourceInfo{}, err
	}
	content := string(data)

	res := ResourceInfo{
		Name:     name,
		FullName: "f5xc_" + name,
	}

	// Extract frontmatter
	res.Category = extractSubcategory(content)
	if res.Category == "" {
		res.Category = resource.GetCategory(name)
	}

	res.Description = cleanDescription(extractDescription(content))

	// Check if promoted
	res.IsPromoted = isPromoted(name, config.Promote)

	// Extract required fields
	res.RequiredFields = extractRequiredFields(content)

	// Extract OneOf groups with context
	res.OneOfGroups = extractOneOfGroupsWithContext(content)

	// Extract server defaults
	res.ServerDefaults = extractServerDefaults(content)

	// Extract minimal config
	res.MinimalConfig = extractMinimalConfig(content, name)

	// Get dependencies from resource package
	if deps, ok := resource.Dependencies[name]; ok {
		res.Dependencies = deps
	}

	return res, nil
}

func extractSubcategory(content string) string {
	lines := strings.Split(content, "\n")
	inFrontmatter := false
	for _, line := range lines {
		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			break
		}
		if !inFrontmatter {
			continue
		}
		if strings.HasPrefix(line, "subcategory:") {
			val := strings.TrimPrefix(line, "subcategory:")
			val = strings.TrimSpace(val)
			val = strings.Trim(val, `"'`)
			return val
		}
	}
	return ""
}

func extractDescription(content string) string {
	lines := strings.Split(content, "\n")
	inFrontmatter := false
	inDesc := false
	var descLines []string
	descIndent := ""

	for _, line := range lines {
		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			break
		}
		if !inFrontmatter {
			continue
		}
		if inDesc {
			if strings.HasPrefix(line, descIndent) && len(strings.TrimSpace(line)) > 0 {
				descLines = append(descLines, strings.TrimPrefix(line, descIndent))
				continue
			}
			if len(strings.TrimSpace(line)) == 0 {
				continue
			}
			break
		}
		if strings.HasPrefix(line, "description:") {
			rest := strings.TrimPrefix(line, "description:")
			rest = strings.TrimSpace(rest)
			if rest == "|-" || rest == "|" {
				inDesc = true
				descIndent = "  "
			} else {
				return strings.Trim(rest, `"'`)
			}
		}
	}

	if len(descLines) > 0 {
		return strings.Join(descLines, " ")
	}
	return ""
}

func cleanDescription(desc string) string {
	// Remove common prefixes
	prefixes := []string{
		"Manages a ",
		"Manages an ",
		"Manages ",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(desc, prefix) {
			desc = strings.TrimPrefix(desc, prefix)
			// Capitalize first letter
			if len(desc) > 0 {
				desc = strings.ToUpper(string(desc[0])) + desc[1:]
			}
			break
		}
	}

	// Remove "resource in F5 Distributed Cloud for " prefix
	if idx := strings.Index(desc, " resource in F5 Distributed Cloud for "); idx != -1 {
		desc = desc[idx+len(" resource in F5 Distributed Cloud for "):]
		if len(desc) > 0 {
			desc = strings.ToUpper(string(desc[0])) + desc[1:]
		}
	}

	// Remove " in F5 Distributed Cloud." suffix
	desc = strings.TrimSuffix(desc, " in F5 Distributed Cloud.")

	// Clean up trailing periods and whitespace
	desc = strings.TrimSpace(desc)
	desc = strings.TrimSuffix(desc, ".")

	return desc
}

func isPromoted(name string, promotePatterns []string) bool {
	for _, pattern := range promotePatterns {
		if strings.HasPrefix(pattern, "resources/") {
			pattern = strings.TrimPrefix(pattern, "resources/")
			pattern = strings.TrimSuffix(pattern, "*")
			if strings.HasPrefix(name, pattern) {
				return true
			}
		}
	}
	return false
}

func extractRequiredFields(content string) []string {
	var fields []string
	seen := make(map[string]bool)

	// Find argument reference section
	lines := strings.Split(content, "\n")
	inArgRef := false
	for _, line := range lines {
		if strings.Contains(line, "## Argument Reference") ||
			strings.Contains(line, "### Minimum Configuration") ||
			strings.Contains(line, "### Metadata Argument Reference") ||
			strings.Contains(line, "### Spec Argument Reference") {
			inArgRef = true
			continue
		}
		if strings.HasPrefix(line, "## ") && !strings.Contains(line, "Argument Reference") {
			inArgRef = false
		}
		if !inArgRef {
			continue
		}

		// Match required fields pattern
		matches := requiredFieldRE.FindStringSubmatch(line)
		if len(matches) > 1 {
			field := strings.Trim(matches[1], "`")
			if !seen[field] {
				seen[field] = true
				fields = append(fields, field)
			}
		}

		// Also check for "Required fields:" list
		if strings.HasPrefix(strings.TrimSpace(line), "- `") && strings.Contains(line, "`") {
			start := strings.Index(line, "`") + 1
			end := strings.Index(line[start:], "`")
			if end > 0 {
				field := line[start : start+end]
				if !seen[field] {
					seen[field] = true
					fields = append(fields, field)
				}
			}
		}
	}

	// Ensure name and namespace are first if present
	result := []string{}
	if seen["name"] {
		result = append(result, "name")
	}
	if seen["namespace"] {
		result = append(result, "namespace")
	}
	for _, f := range fields {
		if f != "name" && f != "namespace" {
			result = append(result, f)
		}
	}

	return result
}

func extractOneOfGroupsWithContext(content string) []OneOfGroup {
	var groups []OneOfGroup
	seen := make(map[string]bool)

	// Find the first terraform code block to extract OneOf context
	blockMatch := codeBlockRE.FindStringSubmatch(content)
	if blockMatch == nil {
		return groups
	}
	codeBlock := blockMatch[1]

	// Also build a set of block-type fields from the code example
	blockFields := findBlockFields(codeBlock)

	// Track block context using a stack of block names.
	// We maintain depth relative to the resource block itself.
	lines := strings.Split(codeBlock, "\n")
	var blockStack []string // stack of block names within the resource
	depth := 0              // depth within the resource block (0 = outside resource)
	insideResource := false
	blockNameRE := regexp.MustCompile(`^\s*(\w+)\s*\{`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		isComment := strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#")

		// Match OneOf comments first (before modifying depth)
		matches := oneOfCommentRE.FindStringSubmatch(line)
		if len(matches) > 1 {
			rawList := matches[1]
			options := strings.Fields(rawList)
			if len(options) == 0 {
				continue
			}

			// Determine parent block from stack
			parent := ""
			if len(blockStack) > 0 {
				parent = blockStack[len(blockStack)-1]
			}

			// Build dedup key
			sorted := make([]string, len(options))
			copy(sorted, options)
			sort.Strings(sorted)
			key := parent + ":" + strings.Join(sorted, ",")
			if seen[key] {
				continue
			}
			seen[key] = true

			groups = append(groups, OneOfGroup{
				Parent:  parent,
				Options: options,
			})
			continue
		}

		if isComment || len(trimmed) == 0 {
			continue
		}

		// Count braces
		opens := strings.Count(trimmed, "{")
		closes := strings.Count(trimmed, "}")

		// Process opens
		for i := 0; i < opens; i++ {
			if !insideResource {
				// Check if this is the resource block opening
				if strings.Contains(trimmed, "resource \"f5xc_") {
					insideResource = true
					depth = 1
					break // only one open brace on this line for resource
				}
				// Skip terraform/required_providers blocks
				break
			}
			depth++
			// Only push a name on the first open brace of this line
			if i == 0 {
				if m := blockNameRE.FindStringSubmatch(trimmed); m != nil {
					name := m[1]
					if name != "resource" && name != "terraform" && name != "required_providers" {
						blockStack = append(blockStack, name)
					}
				}
			}
		}

		// Process closes
		if insideResource {
			for i := 0; i < closes; i++ {
				depth--
				if depth <= 0 {
					insideResource = false
					blockStack = nil
					break
				}
				if len(blockStack) > 0 {
					blockStack = blockStack[:len(blockStack)-1]
				}
			}
		}
	}

	// blockFields is available but not needed now since isEmptyBlockField handles type
	_ = blockFields

	return groups
}

// findBlockFields scans a code block and returns a set of field names that
// are used as blocks (i.e., "fieldname {" or "fieldname {}").
func findBlockFields(codeBlock string) map[string]bool {
	blockFields := make(map[string]bool)
	blockRE := regexp.MustCompile(`^\s*(\w+)\s*\{`)
	for _, line := range strings.Split(codeBlock, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "//") || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		if m := blockRE.FindStringSubmatch(line); m != nil {
			name := m[1]
			if name != "resource" && name != "terraform" && name != "required_providers" && name != "f5xc" {
				blockFields[name] = true
			}
		}
	}
	return blockFields
}

// isEmptyBlockField returns true if the field name is typically used as an
// empty block selector (field {}), false if it's typically an attribute (field = value).
func isEmptyBlockField(name string) bool {
	// Known attribute fields (string values, not blocks)
	knownAttributes := map[string]bool{
		"host_header":               true, // string value: "example.com"
		"server_name":               true, // string value
		"api_specification":         true, // reference
		"api_testing":               true, // reference when not empty
		"captcha_challenge":         true, // config string
		"js_challenge":              true, // config string
		"policy_based_challenge":    true, // reference
		"cookie_stickiness":         true, // config
		"least_active":              true, // config
		"random":                    true, // config
		"ring_hash":                 true, // config
		"round_robin":               true, // selector (but empty block form)
		"source_ip_stickiness":      true, // config
		"http":                      true, // has nested config
		"https":                     true, // has nested config
		"https_auto_cert":           true, // has nested config
		"malware_protection_settings": true,
		"api_rate_limit":            true,
		"rate_limit":                true,
		"active_service_policies":   true,
		"service_policies_from_namespace": true, // but often empty
		"user_id_client_ip":         true,
		"user_identification":       true,
		"app_firewall":              true, // reference block with name/namespace
		"bot_defense":               true,
		"bot_defense_advanced":      true,
		"headers":                   true, // map
		"request_headers_to_remove": true, // list
	}

	if knownAttributes[name] {
		// Check if it's one of the selector variants
		// These are actually blocks when used as selectors:
		selectorsAsBlocks := map[string]bool{
			"round_robin":               true,
			"least_active":              true,
			"random":                    true,
			"cookie_stickiness":         true,
			"source_ip_stickiness":      true,
			"service_policies_from_namespace": true,
			"user_id_client_ip":         true,
		}
		if selectorsAsBlocks[name] {
			return true
		}
		return false
	}

	// Patterns that indicate empty block selectors
	blockPatterns := []string{
		"enable_", "disable_", "no_", "use_", "do_not_",
		"default_", "advertise_",
	}
	for _, prefix := range blockPatterns {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	blockSuffixes := []string{
		"_check", "_config", "_vip", "_security", "_mtls",
	}
	for _, suffix := range blockSuffixes {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}

	// If ends with _header but not "host_header", it's likely a selector block
	if strings.HasSuffix(name, "_header") && name != "host_header" {
		return true
	}

	// Default: assume attribute
	return false
}

func extractServerDefaults(content string) []string {
	// Look for server default markers in the document
	var defaults []string
	seen := make(map[string]bool)

	// First, check for commented list in code blocks
	blockMatch := codeBlockRE.FindStringSubmatch(content)
	if blockMatch != nil {
		codeBlock := blockMatch[1]
		for _, line := range strings.Split(codeBlock, "\n") {
			if m := serverDefaultLineRE.FindStringSubmatch(line); m != nil {
				field := m[1]
				if !seen[field] {
					seen[field] = true
					defaults = append(defaults, field)
				}
			}
		}
	}

	// Also check for Server Default markers in argument reference
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if serverDefaultTagRE.MatchString(line) {
			// Extract field name from the line
			if idx := strings.Index(line, "[`"); idx != -1 {
				end := strings.Index(line[idx+2:], "`]")
				if end > 0 {
					field := line[idx+2 : idx+2+end]
					if !seen[field] {
						seen[field] = true
						defaults = append(defaults, field)
					}
				}
			}
		}
	}

	return defaults
}

func extractMinimalConfig(content, name string) string {
	// Find the first resource block in the example
	blockMatch := codeBlockRE.FindStringSubmatch(content)
	if blockMatch == nil {
		return ""
	}

	codeBlock := blockMatch[1]

	// Find the resource block
	lines := strings.Split(codeBlock, "\n")
	var configLines []string
	inResource := false
	braceDepth := 0
	resourcePattern := fmt.Sprintf(`resource "f5xc_%s"`, name)

	for _, line := range lines {
		if !inResource && strings.Contains(line, resourcePattern) {
			inResource = true
			configLines = append(configLines, line)
			braceDepth = 1
			continue
		}
		if inResource {
			configLines = append(configLines, line)
			braceDepth += strings.Count(line, "{") - strings.Count(line, "}")
			if braceDepth <= 0 {
				break
			}
		}
	}

	if len(configLines) == 0 {
		return ""
	}

	// Clean up the config - remove comments except OneOf comments
	var cleanLines []string
	for _, line := range configLines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		cleanLines = append(cleanLines, line)
	}

	result := strings.Join(cleanLines, "\n")

	// Limit size — truncate at line boundary and close unclosed braces
	if len(result) > 2000 {
		cut := strings.LastIndex(result[:2000], "\n")
		if cut <= 0 {
			cut = 2000
		}
		result = result[:cut]
		opens := strings.Count(result, "{")
		closes := strings.Count(result, "}")
		for i := 0; i < opens-closes; i++ {
			result += "\n}"
		}
	}

	return result
}

func groupByCategory(resources []ResourceInfo) []CategoryInfo {
	catMap := make(map[string]*CategoryInfo)

	for _, res := range resources {
		cat, ok := catMap[res.Category]
		if !ok {
			cat = &CategoryInfo{
				Name:        res.Category,
				Slug:        categorySlug(res.Category),
				Description: categoryDescriptions[res.Category],
			}
			if cat.Description == "" {
				cat.Description = "Resources for " + res.Category
			}
			catMap[res.Category] = cat
		}
		cat.Resources = append(cat.Resources, res)
	}

	// Convert to slice and sort
	var categories []CategoryInfo
	for _, cat := range catMap {
		// Sort resources: promoted first, then alphabetically
		sort.Slice(cat.Resources, func(i, j int) bool {
			if cat.Resources[i].IsPromoted != cat.Resources[j].IsPromoted {
				return cat.Resources[i].IsPromoted
			}
			return cat.Resources[i].Name < cat.Resources[j].Name
		})
		categories = append(categories, *cat)
	}

	// Sort categories by resource count descending
	sort.Slice(categories, func(i, j int) bool {
		return len(categories[i].Resources) > len(categories[j].Resources)
	})

	return categories
}

func truncateAtWord(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	cut := strings.LastIndex(s[:maxLen-3], " ")
	if cut <= 0 {
		cut = maxLen - 3
	}
	return s[:cut] + "..."
}

func categorySlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "/", "-")
	return slug
}

func buildReverseDependencies(resources []ResourceInfo) map[string][]string {
	reverse := make(map[string][]string)
	for _, res := range resources {
		for _, dep := range res.Dependencies {
			reverse[dep] = append(reverse[dep], res.Name)
		}
	}
	// Sort each list
	for k := range reverse {
		sort.Strings(reverse[k])
	}
	return reverse
}

// generateL0 creates docs/llms.txt
func generateL0(config *LLMsConfig, categories []CategoryInfo) error {
	var sb strings.Builder

	// Header
	sb.WriteString("# F5 Distributed Cloud Terraform Provider\n\n")
	sb.WriteString(fmt.Sprintf("> Provider: %s | Registry: <%s>\n\n",
		config.Deprecation.Canonical.Provider,
		config.Deprecation.Canonical.Registry))
	sb.WriteString(fmt.Sprintf("DEPRECATED — do not use: %s\n\n",
		config.Deprecation.Deprecated.Provider))

	// Required Block
	sb.WriteString("## Required Block\n\n")
	sb.WriteString("    terraform {\n")
	sb.WriteString("      required_providers {\n")
	sb.WriteString("        f5xc = {\n")
	sb.WriteString("          source = \"f5xc-salesdemos/f5xc\"\n")
	sb.WriteString("        }\n")
	sb.WriteString("      }\n")
	sb.WriteString("    }\n\n")

	// Syntax Rules
	sb.WriteString("## Syntax Rules\n\n")
	sb.WriteString("- OneOf selectors: use empty block `field {}`, never `field = true`\n")
	sb.WriteString("- Cross-resource refs: block with name + namespace attributes\n")
	sb.WriteString("- Boolean attributes: use `= true` / `= false`\n")
	sb.WriteString("- Fields marked \"Server applies default when omitted\" can be safely omitted\n\n")

	// Guides
	sb.WriteString("## Guides\n\n")
	guides := []struct{ name, desc string }{
		{"authentication", "API token, P12 certificate, and PEM auth methods"},
		{"http-loadbalancer", "Deploy production load balancers with security"},
	}
	for _, g := range guides {
		sb.WriteString(fmt.Sprintf("- [%s](guides/%s) : %s\n", g.name, g.name, g.desc))
	}
	sb.WriteString("\n")

	// Functions
	sb.WriteString("## Functions\n\n")
	functions := []struct{ name, desc string }{
		{"blindfold()", "Encrypt secrets using local public key encryption"},
		{"blindfold_file()", "Encrypt file contents directly"},
	}
	for _, f := range functions {
		fname := strings.TrimSuffix(f.name, "()")
		sb.WriteString(fmt.Sprintf("- [%s](functions/%s) : %s\n", f.name, fname, f.desc))
	}
	sb.WriteString("\n")

	// Resource Categories
	sb.WriteString("## Resource Categories\n\n")
	for _, cat := range categories {
		sb.WriteString(fmt.Sprintf("- [%s](_llms-txt/%s.txt) (%d resources): %s\n",
			cat.Name, cat.Slug, len(cat.Resources), cat.Description))
	}
	sb.WriteString("\n")

	// Data Sources note
	sb.WriteString("## Data Sources\n\n")
	sb.WriteString("Read-only data sources for querying existing F5 XC objects across all resource categories.\n")

	return os.WriteFile("docs/llms.txt", []byte(sb.String()), 0644)
}

// generateL1 creates docs/_llms-txt/<subcategory>.txt
func generateL1(cat CategoryInfo) error {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("# %s\n\n", cat.Name))
	sb.WriteString(fmt.Sprintf("%s\n\n", cat.Description))

	// Resources
	sb.WriteString("## Resources\n\n")
	for _, res := range cat.Resources {
		desc := truncateAtWord(res.Description, 80)
		sb.WriteString(fmt.Sprintf("- [%s](_llms-txt/resources/%s.txt) : %s\n",
			res.Name, res.Name, desc))
	}

	// Dependency chain (only if there are known dependencies)
	if len(cat.Resources) > 0 {
		deps := buildCategoryDependencyChain(cat.Resources)
		if deps != "" {
			sb.WriteString("\n## Dependency Chain\n\n")
			sb.WriteString(deps + "\n")
		}
	}

	path := fmt.Sprintf("docs/_llms-txt/%s.txt", cat.Slug)
	return os.WriteFile(path, []byte(sb.String()), 0644)
}

func buildCategoryDependencyChain(resources []ResourceInfo) string {
	// Build a simple dependency chain showing common patterns
	depMap := make(map[string][]string)
	names := make(map[string]bool)

	for _, res := range resources {
		names[res.Name] = true
		for _, dep := range res.Dependencies {
			depMap[dep] = append(depMap[dep], res.Name)
		}
	}

	// Find roots (resources that nothing depends on within this category)
	// and build chains
	var chains []string

	// Common patterns
	if names["http_loadbalancer"] && names["origin_pool"] && names["healthcheck"] {
		chains = append(chains, "namespace → healthcheck → origin_pool → http_loadbalancer")
	}
	if names["tcp_loadbalancer"] && names["origin_pool"] {
		chains = append(chains, "namespace → origin_pool → tcp_loadbalancer")
	}
	if names["app_firewall"] {
		chains = append(chains, "namespace → app_firewall → http_loadbalancer")
	}

	if len(chains) == 0 {
		return ""
	}
	return strings.Join(chains, "\n")
}

// generateL2 creates docs/_llms-txt/resources/<name>.txt
func generateL2(res ResourceInfo, reverseDeps map[string][]string) error {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("# %s\n\n", res.FullName))
	sb.WriteString(fmt.Sprintf("Category: %s\n", res.Category))
	sb.WriteString(fmt.Sprintf("%s\n\n", res.Description))

	// Required fields
	if len(res.RequiredFields) > 0 {
		sb.WriteString("## Required\n\n")
		for _, f := range res.RequiredFields {
			sb.WriteString(fmt.Sprintf("- %s\n", f))
		}
		sb.WriteString("\n")
	}

	// OneOf Groups
	if len(res.OneOfGroups) > 0 {
		sb.WriteString("## OneOf Groups\n\n")
		for _, group := range res.OneOfGroups {
			if group.Parent != "" {
				sb.WriteString(fmt.Sprintf("Within %s, pick exactly one:\n", group.Parent))
			} else {
				sb.WriteString("Pick exactly one:\n")
			}
			for _, opt := range group.Options {
				// Determine if block or attribute based on naming conventions
				// Empty-block patterns: enable_*, disable_*, no_*, use_*, *_check, *_config,
				// *_policy (no args), *_vip, *_security, *_header (selector), *_mtls
				// Attribute patterns: host_header (string), specific known attributes
				isBlock := isEmptyBlockField(opt)

				if isBlock {
					sb.WriteString(fmt.Sprintf("  %s {}\n", opt))
				} else {
					sb.WriteString(fmt.Sprintf("  %s = \"...\"\n", opt))
				}
			}
			sb.WriteString("\n")
		}
	}

	// Server Defaults
	if len(res.ServerDefaults) > 0 {
		sb.WriteString("## Server Defaults (safe to omit)\n\n")
		for _, f := range res.ServerDefaults {
			sb.WriteString(fmt.Sprintf("- %s\n", f))
		}
		sb.WriteString("\n")
	}

	// Minimal Valid Config
	if res.MinimalConfig != "" {
		sb.WriteString("## Minimal Valid Config\n\n")
		sb.WriteString("```terraform\n")
		sb.WriteString(res.MinimalConfig)
		if !strings.HasSuffix(res.MinimalConfig, "\n") {
			sb.WriteString("\n")
		}
		sb.WriteString("```\n\n")
	}

	// Dependencies
	sb.WriteString("## Dependencies\n\n")
	if len(res.Dependencies) > 0 {
		sb.WriteString("Requires: " + strings.Join(res.Dependencies, ", ") + "\n")
	} else {
		sb.WriteString("Requires: none\n")
	}
	if usedBy, ok := reverseDeps[res.Name]; ok && len(usedBy) > 0 {
		sb.WriteString("Used by: " + strings.Join(usedBy, ", ") + "\n")
	}

	path := fmt.Sprintf("docs/_llms-txt/resources/%s.txt", res.Name)
	content := sb.String()

	// Check size
	if len(content) > 8192 {
		fmt.Printf("Warning: %s exceeds 8KB (%d bytes)\n", res.Name, len(content))
	}

	return os.WriteFile(path, []byte(content), 0644)
}

func generateJSONIndex(config *LLMsConfig, categories []CategoryInfo, reverseDeps map[string][]string) error {
	idx := JSONIndex{
		Version: "0.1.0",
		Provider: JSONProvider{
			Source:   config.Deprecation.Canonical.Provider,
			Registry: config.Deprecation.Canonical.Registry,
			RequiredBlock: `terraform {
  required_providers {
    f5xc = {
      source = "f5xc-salesdemos/f5xc"
    }
  }
}`,
			SyntaxRules: []string{
				"OneOf selectors: use empty block `field {}`, never `field = true`",
				"Cross-resource refs: block with name + namespace attributes",
				"Boolean attributes: use `= true` / `= false`",
				`Fields marked "Server applies default when omitted" can be safely omitted`,
			},
		},
		Resources: make(map[string]JSONResource),
	}

	for _, cat := range categories {
		jcat := JSONCategory{
			Name:          cat.Name,
			Slug:          cat.Slug,
			Description:   cat.Description,
			ResourceCount: len(cat.Resources),
		}
		for _, res := range cat.Resources {
			jcat.Resources = append(jcat.Resources, res.Name)
		}
		jcat.DependencyChain = buildCategoryDependencyChain(cat.Resources)
		idx.Categories = append(idx.Categories, jcat)

		for _, res := range cat.Resources {
			var oneOfGroups []JSONOneOfGroup
			for _, g := range res.OneOfGroups {
				oneOfGroups = append(oneOfGroups, JSONOneOfGroup{
					Parent: g.Parent,
					Fields: g.Options,
				})
			}
			deps := JSONDependencies{Requires: res.Dependencies}
			if deps.Requires == nil {
				deps.Requires = []string{}
			}
			if usedBy, ok := reverseDeps[res.Name]; ok {
				deps.UsedBy = usedBy
			}
			idx.Resources[res.Name] = JSONResource{
				Category:       cat.Slug,
				Description:    res.Description,
				Required:       res.RequiredFields,
				OneOfGroups:    oneOfGroups,
				ServerDefaults: res.ServerDefaults,
				MinimalConfig:  res.MinimalConfig,
				Dependencies:   deps,
				ImportSyntax:   fmt.Sprintf("terraform import f5xc_%s.example namespace/name", res.Name),
			}
		}
	}

	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON index: %w", err)
	}
	data = append(data, '\n')
	return os.WriteFile("docs/terraform-llms-index.json", data, 0644)
}
