// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

//go:build ignore
// +build ignore

// This tool generates docs/resources/*.ai.json files from existing documentation.
// These JSON files are machine-readable metadata consumed by AI assistants.
//
// Usage:
//
//	go run tools/generate-ai-metadata.go
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
	pkgschema "github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/schema"
)

// AIMetadata is the JSON structure written to docs/resources/<name>.ai.json.
type AIMetadata struct {
	Resource           string            `json:"resource"`
	Category           string            `json:"category"`
	Description        string            `json:"description"`
	Dependencies       []string          `json:"dependencies"`
	SyntaxRules        map[string]string `json:"syntax_rules"`
	RequiredFields     []string          `json:"required_fields"`
	OneOfGroups        []OneOfGroup      `json:"oneof_groups"`
	ServerDefaults     []string          `json:"server_defaults"`
	MinimalValidConfig string            `json:"minimal_valid_config,omitempty"`
}

// OneOfGroup represents a mutually exclusive group of fields.
type OneOfGroup struct {
	Key     string        `json:"key"`
	Options []OneOfOption `json:"options"`
}

// OneOfOption is one member of a OneOf group.
type OneOfOption struct {
	Name string `json:"name"`
	Type string `json:"type"` // "block" or "attribute"
}

// terraformSchema holds the optional cross-check schema loaded at startup.
var terraformSchema *pkgschema.TerraformSchema

// oneOfCommentRE matches:  // One of the arguments from this list "a b c" must be set
var oneOfCommentRE = regexp.MustCompile(`//\s*One of the arguments from this list "([^"]+)" must be set`)

// serverDefaultLineRE matches: # - fieldname
var serverDefaultLineRE = regexp.MustCompile(`^#\s+-\s+(\S+)\s*$`)

// codeBlockRE captures the first ```terraform ... ``` block.
var codeBlockRE = regexp.MustCompile("(?s)```(?:terraform|hcl)\n(.*?)```")

func main() {
	// Try to load cached terraform schema (optional).
	if data, err := os.ReadFile("tools/terraform-schema.json"); err == nil {
		ts, err := pkgschema.ParseTerraformSchema(data)
		if err == nil {
			terraformSchema = ts
			fmt.Printf("Loaded Terraform schema (%d resources)\n", len(ts.Resources))
		} else {
			fmt.Printf("Note: Could not parse terraform schema: %v\n", err)
		}
	} else {
		fmt.Println("Note: tools/terraform-schema.json not found — generating metadata without schema cross-check")
	}

	docsDir := "docs/resources"
	files, err := filepath.Glob(filepath.Join(docsDir, "*.md"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding doc files: %v\n", err)
		os.Exit(1)
	}

	succeeded := 0
	failed := 0
	for _, f := range files {
		if err := processFile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", f, err)
			failed++
		} else {
			succeeded++
		}
	}

	fmt.Printf("AI metadata generation complete: %d succeeded, %d failed\n", succeeded, failed)
	if failed > 0 {
		os.Exit(1)
	}
}

// processFile reads a single .md file and writes the corresponding .ai.json.
func processFile(mdPath string) error {
	base := filepath.Base(mdPath)
	if strings.Contains(base, "_nested_blocks") {
		return nil // skip old nested-block pages
	}
	name := strings.TrimSuffix(base, ".md")

	data, err := os.ReadFile(mdPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", mdPath, err)
	}
	content := string(data)

	meta := buildMetadata(name, content)

	out, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", name, err)
	}

	outPath := filepath.Join("docs/resources", name+".ai.json")
	if err := os.WriteFile(outPath, append(out, '\n'), 0644); err != nil {
		return fmt.Errorf("write %s: %w", outPath, err)
	}
	fmt.Printf("Generated: %s\n", outPath)
	return nil
}

// buildMetadata constructs the AIMetadata for one resource.
func buildMetadata(name, content string) AIMetadata {
	resourceFull := "f5xc_" + name

	// --- description from frontmatter ---
	description := extractDescription(content)

	// --- category ---
	category := resource.GetCategory(name)

	// --- dependencies ---
	deps := resource.Dependencies[name]
	if deps == nil {
		deps = []string{}
	}

	// --- required fields ---
	requiredFields := buildRequiredFields(name)

	// --- OneOf groups ---
	oneofGroups := extractOneOfGroups(name, content)

	// --- server defaults ---
	serverDefaults := extractServerDefaults(content)

	// --- minimal valid config (first terraform code block) ---
	minConfig := extractMinimalConfig(content)

	meta := AIMetadata{
		Resource:     resourceFull,
		Category:     category,
		Description:  description,
		Dependencies: deps,
		SyntaxRules: map[string]string{
			"oneof_block_syntax": "Use empty block {} for OneOf selectors, never = true",
			"reference_syntax":   "Cross-resource refs use a block with name + namespace",
		},
		RequiredFields: requiredFields,
		OneOfGroups:    oneofGroups,
		ServerDefaults: serverDefaults,
	}
	if len(minConfig) > 0 && len(minConfig) <= 2000 {
		meta.MinimalValidConfig = minConfig
	}

	return meta
}

// extractDescription pulls the description from YAML frontmatter.
// Frontmatter looks like:
//
//	description: |-
//	  First line.
//	  Second line.
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
			break // end of frontmatter
		}
		if !inFrontmatter {
			continue
		}
		if inDesc {
			// Continue collecting indented description lines.
			if strings.HasPrefix(line, descIndent) && len(strings.TrimSpace(line)) > 0 {
				descLines = append(descLines, strings.TrimPrefix(line, descIndent))
				continue
			}
			// Non-indented line after desc block — done.
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
				descIndent = "  " // YAML block scalars indent by 2
			} else {
				// Inline description
				return strings.Trim(rest, `"'`)
			}
		}
	}

	if len(descLines) > 0 {
		return strings.Join(descLines, " ")
	}
	return ""
}

// buildRequiredFields returns the required field list for a resource.
// Always includes name and namespace; adds resource-specific extras.
func buildRequiredFields(name string) []string {
	fields := []string{"name", "namespace"}
	if name == "http_loadbalancer" || name == "tcp_loadbalancer" || name == "udp_loadbalancer" {
		fields = append(fields, "domains")
	}
	return fields
}

// extractOneOfGroups parses OneOf comment lines from the markdown content.
// Deduplicates by the sorted option-set key and resolves block vs attribute
// using the Terraform schema when available.
func extractOneOfGroups(name, content string) []OneOfGroup {
	seen := make(map[string]bool)
	var groups []OneOfGroup

	matches := oneOfCommentRE.FindAllStringSubmatch(content, -1)
	for _, m := range matches {
		rawList := m[1] // e.g. "advertise_custom advertise_on_public advertise_on_public_default_vip do_not_advertise"
		options := strings.Fields(rawList)
		if len(options) == 0 {
			continue
		}

		// Build a stable dedup key from sorted options.
		sorted := make([]string, len(options))
		copy(sorted, options)
		sort.Strings(sorted)
		key := strings.Join(sorted, ",")
		if seen[key] {
			continue
		}
		seen[key] = true

		// Build the option list with type resolution.
		var opts []OneOfOption
		for _, opt := range options {
			kind := "unknown"
			if terraformSchema != nil {
				resourceFull := "f5xc_" + name
				if rs, ok := terraformSchema.Resources[resourceFull]; ok {
					switch rs.FieldKind(opt) {
					case pkgschema.FieldKindBlock:
						kind = "block"
					case pkgschema.FieldKindAttribute:
						kind = "attribute"
					}
				}
			}
			opts = append(opts, OneOfOption{Name: opt, Type: kind})
		}

		groups = append(groups, OneOfGroup{
			Key:     options[0], // Use first option as the group key (matches original order)
			Options: opts,
		})
	}

	return groups
}

// extractServerDefaults parses the "# - fieldname" comment block from the
// first terraform code block in the document.
func extractServerDefaults(content string) []string {
	// Find the first code block.
	blockMatch := codeBlockRE.FindStringSubmatch(content)
	if blockMatch == nil {
		return []string{}
	}
	codeBlock := blockMatch[1]

	var defaults []string
	for _, line := range strings.Split(codeBlock, "\n") {
		if m := serverDefaultLineRE.FindStringSubmatch(line); m != nil {
			defaults = append(defaults, m[1])
		}
	}
	return defaults
}

// extractMinimalConfig returns the content of the first terraform/hcl code block.
func extractMinimalConfig(content string) string {
	m := codeBlockRE.FindStringSubmatch(content)
	if m == nil {
		return ""
	}
	return strings.TrimSpace(m[1])
}
