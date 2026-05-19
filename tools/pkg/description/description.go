// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

// Package description provides text transformation functions for cleaning,
// formatting, and enriching OpenAPI descriptions into Terraform-quality
// documentation strings.
package description

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/naming"
	"github.com/f5xc-salesdemos/terraform-provider-f5xc/tools/pkg/resource"
)

// Clean removes example annotations, validation rules, internal vendor extensions,
// and normalises whitespace in a raw OpenAPI description.
func Clean(desc string, fieldPath string) string {
	// Remove example and validation rules sections
	desc = regexp.MustCompile(`\s*Example:.*`).ReplaceAllString(desc, "")
	desc = regexp.MustCompile(`\s*Validation Rules:.*`).ReplaceAllString(desc, "")

	// Remove x-example annotations (OpenAPI 2.0 vendor extension for Swagger UI examples)
	// Pattern: x-example: "value" or x-example: 'value' embedded in description text
	desc = regexp.MustCompile(`\s*x-example:\s*["']?[^"'\n]*["']?`).ReplaceAllString(desc, "")
	// Also handle x-required annotations
	desc = regexp.MustCompile(`\s*x-required\s*`).ReplaceAllString(desc, "")

	// Remove ves.io validation annotations (common pattern in F5 XC specs)
	// These are internal protobuf validation rules that leaked into OpenAPI descriptions
	// Pattern: ves.io.schema.rules.xxx.yyy: value or ves.io.schema.xxx: value
	desc = regexp.MustCompile(`\s*ves\.io\.schema[^\s]*:\s*\S+`).ReplaceAllString(desc, "")
	desc = regexp.MustCompile(`\s*ves\.io\.[^\s]*:\s*\[.*?\]`).ReplaceAllString(desc, "")

	// Remove "Required: YES" or "Required: NO" annotations
	desc = regexp.MustCompile(`\s*Required:\s*(YES|NO)\s*`).ReplaceAllString(desc, " ")
	// Note: "Exclusive with [xxx]" patterns are intentionally preserved.
	// This text is the only documentation hint about conflicting fields until
	// ConflictsWith validators are rendered. Stripping it removes valuable
	// user-facing documentation.

	// Normalize generic empty message descriptions to user-friendly text
	// "Empty. This can be used for messages where no values are needed" → "Enable this option"
	desc = regexp.MustCompile(`(?i)Empty\.?\s*This can be used for messages where no values are needed\.?`).ReplaceAllString(desc, "Enable this option")
	// Also handle variations
	desc = regexp.MustCompile(`(?i)This can be used for messages where no values are needed\.?`).ReplaceAllString(desc, "Enable this option")

	// Normalize "Shape of the X specification" to "Configuration for X"
	// This converts internal F5 terminology to user-friendly Terraform terminology
	desc = regexp.MustCompile(`(?i)Shape of the ([^\s]+) specification`).ReplaceAllString(desc, "Configuration for $1")
	desc = regexp.MustCompile(`(?i)Shape of ([^\s]+) specification`).ReplaceAllString(desc, "Configuration for $1")

	// Remove escaped quotes and backslashes from raw spec data
	desc = strings.ReplaceAll(desc, `\"`, `"`)
	desc = strings.ReplaceAll(desc, `\\`, `\`)
	// Normalize whitespace
	desc = regexp.MustCompile(`[\n\r]+`).ReplaceAllString(desc, " ")
	desc = regexp.MustCompile(`\s+`).ReplaceAllString(desc, " ")
	// Escape quotes for Go string literals
	desc = strings.ReplaceAll(desc, `"`, "'")
	desc = strings.TrimSpace(desc)
	// Remove trailing periods that were left after cleanup
	desc = regexp.MustCompile(`\.\s*\.`).ReplaceAllString(desc, ".")
	// Normalize example names from F5 internal conventions ("my-*") to provider standard ("example-*")
	desc = naming.NormalizeExampleNames(desc)
	return desc
}

// ExtractDefault attempts to extract a default value from description text.
// It returns the extracted default (or nil) and the cleaned description.
func ExtractDefault(desc string) (interface{}, string) {
	if desc == "" {
		return nil, desc
	}

	// Patterns to match defaults mentioned in description text
	// Pattern 1: "Defaults to X" or "Default to X" with optional units (including "or Ys" alternative)
	// Pattern 2: "Default value is X" or "Default value: X"
	// Pattern 3: "defaults to X" (lowercase)
	patterns := []string{
		// Match "Defaults to 30000ms or 30s" or "Defaults to 30000ms" or "Defaults to true"
		// Captures first value, removes optional " or Xs" alternative
		`[Dd]efaults?\s+to\s+(\d+(?:ms|s|%)?|\d+\.\d+|true|false)(?:\s+or\s+\d+(?:ms|s)?)?`,
		// Match "Default value is /graphql" or "Default value is true"
		`[Dd]efault\s+value\s+(?:is|:)\s+([^\s.,]+)`,
		// Match "default is 10" or "default: 10"
		`[Dd]efault\s+(?:is|:)\s+(\d+(?:ms|s|%)?|\d+\.\d+|true|false)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(desc)
		if len(match) >= 2 {
			defaultVal := strings.TrimSpace(match[1])
			// Skip if it looks like a placeholder or invalid value
			upperVal := strings.ToUpper(defaultVal)
			if upperVal == "NONE" || upperVal == "INVALID" || upperVal == "UNKNOWN" || upperVal == "UNSPECIFIED" {
				continue
			}
			// Remove the matched text from description to avoid duplication
			cleanedDesc := re.ReplaceAllString(desc, "")
			cleanedDesc = strings.TrimSpace(cleanedDesc)
			// Clean up any trailing punctuation/whitespace artifacts
			cleanedDesc = regexp.MustCompile(`\s+\.`).ReplaceAllString(cleanedDesc, ".")
			cleanedDesc = regexp.MustCompile(`\.\s*\.`).ReplaceAllString(cleanedDesc, ".")
			return defaultVal, cleanedDesc
		}
	}

	return nil, desc
}

// FormatDefault appends a default value to a description per HashiCorp standards.
// Format: "Defaults to `value`."
func FormatDefault(desc string, defaultValue interface{}) string {
	if defaultValue == nil {
		return desc
	}

	// Convert default value to string
	defaultStr := fmt.Sprintf("%v", defaultValue)
	if defaultStr == "" || defaultStr == "<nil>" {
		return desc
	}

	// Skip invalid/placeholder defaults using EXACT match only
	// These are sentinel values that indicate "no value" rather than actual defaults
	invalidDefaults := map[string]bool{
		"INVALID":     true,
		"NONE":        true,
		"UNKNOWN":     true,
		"UNSPECIFIED": true,
		"0":           true, // Integer zero is often just the zero value
	}

	// Check if default is EXACTLY an invalid/placeholder value
	// Use exact matching to preserve valid defaults like XFCC_NONE, MTLS_NONE, etc.
	upperDefault := strings.ToUpper(defaultStr)
	if invalidDefaults[upperDefault] {
		return desc
	}

	// Ensure description ends properly before adding default info
	desc = strings.TrimSpace(desc)
	if desc != "" && !strings.HasSuffix(desc, ".") && !strings.HasSuffix(desc, ":") {
		desc += "."
	}

	return fmt.Sprintf("%s Defaults to `%s`.", desc, defaultStr)
}

// FormatEnum adds AI-parseable enum metadata and human-readable values to description.
// Format: "[Enum: val1|val2|val3] Human description. Possible values are `val1`, `val2`, `val3`."
// The [Enum: ...] prefix enables AI tools to deterministically extract valid values from
// `terraform providers schema -json` output without parsing natural language.
func FormatEnum(desc string, enumValues []interface{}) string {
	if len(enumValues) == 0 {
		return desc
	}

	// Convert enum values to strings (raw values for AI prefix, backtick for human)
	var rawValues []string
	var formattedValues []string
	for _, v := range enumValues {
		str := fmt.Sprintf("%v", v)
		// Skip empty or very long values
		if str == "" || len(str) > 50 {
			continue
		}
		rawValues = append(rawValues, str)
		formattedValues = append(formattedValues, fmt.Sprintf("`%s`", str))
	}

	if len(rawValues) == 0 {
		return desc
	}

	// Build AI-parseable prefix: [Enum: val1|val2|val3]
	aiPrefix := fmt.Sprintf("[Enum: %s]", strings.Join(rawValues, "|"))

	// Ensure description ends properly before adding enum info
	desc = strings.TrimSpace(desc)
	if desc != "" && !strings.HasSuffix(desc, ".") && !strings.HasSuffix(desc, ":") {
		desc += "."
	}

	// Format based on number of values per HashiCorp standards
	var humanSuffix string
	if len(formattedValues) == 1 {
		humanSuffix = fmt.Sprintf("The only possible value is %s.", formattedValues[0])
	} else {
		humanSuffix = fmt.Sprintf("Possible values are %s.", strings.Join(formattedValues, ", "))
	}

	// Combine: [AI prefix] Human description. Human enum list.
	if desc == "" {
		return fmt.Sprintf("%s %s", aiPrefix, humanSuffix)
	}
	return fmt.Sprintf("%s %s %s", aiPrefix, desc, humanSuffix)
}

// GetResourceAIMetadata generates AI-parseable metadata prefix for a resource.
// Format: [Category: X] [Namespace: required|optional] [DependsOn: res1, res2]
func GetResourceAIMetadata(resourceName string) string {
	var parts []string

	// Add category if known
	if category, ok := resource.AICategories[resourceName]; ok {
		parts = append(parts, fmt.Sprintf("[Category: %s]", category))
	}

	// Add namespace requirement
	// Most F5 XC resources require a namespace except for system-level resources
	systemResources := map[string]bool{
		"namespace": true, "tenant": true, "role": true, "allowed_tenant": true,
		"dns_zone": true, "dns_domain": true, "virtual_site": true,
		"cloud_credentials": true, "certificate": true, "trusted_ca_list": true,
	}
	if systemResources[resourceName] {
		parts = append(parts, "[Namespace: not_required]")
	} else {
		parts = append(parts, "[Namespace: required]")
	}

	// Add dependencies if known
	if deps, ok := resource.Dependencies[resourceName]; ok && len(deps) > 0 {
		parts = append(parts, fmt.Sprintf("[DependsOn: %s]", strings.Join(deps, ", ")))
	}

	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " ")
}

// TransformResourceDescription converts technical API descriptions into user-friendly
// Terraform resource descriptions following HashiCorp best practices.
// Pattern: "Manages a [Resource] in F5 Distributed Cloud [for purpose/capability]."
func TransformResourceDescription(resourceName, rawDescription string) string {
	humanName := naming.ToHumanReadableName(resourceName)

	// Clean and normalize the raw description first
	// Pass empty fieldPath since this is the resource-level description
	desc := Clean(rawDescription, "")
	desc = strings.TrimSpace(desc)

	// Generate the human-readable description
	var humanDesc string

	// If empty, use default
	if desc == "" {
		humanDesc = fmt.Sprintf("Manages a %s resource in F5 Distributed Cloud.", humanName)
	} else {
		// Detect and transform common technical description patterns
		lowerDesc := strings.ToLower(desc)

		// Pattern 1: "Shape of the X specification" -> extract X and make user-friendly
		if strings.Contains(lowerDesc, "shape of") {
			humanDesc = GenerateCapabilityDescriptionOnly(resourceName, humanName, desc)
		} else if strings.HasSuffix(lowerDesc, " object") || strings.HasSuffix(lowerDesc, " configuration") ||
			strings.HasSuffix(lowerDesc, " spec") || strings.HasSuffix(lowerDesc, " specification") {
			// Pattern 2: "X object" or "X configuration" - technical object reference
			humanDesc = GenerateCapabilityDescriptionOnly(resourceName, humanName, desc)
		} else {
			// Pattern 3: Already starts with a verb like "Create", "Configure", "Define"
			// Transform to "Manages" for consistency
			actionVerbs := []string{"create", "configure", "define", "set up", "establish", "provision"}
			matched := false
			for _, verb := range actionVerbs {
				if strings.HasPrefix(lowerDesc, verb) {
					// Replace the action verb with "Manages" for Terraform consistency
					remainder := desc[len(verb):]
					remainder = strings.TrimPrefix(remainder, "s") // handle "Creates" -> "Create"
					remainder = strings.TrimSpace(remainder)
					if remainder != "" {
						// Clean up articles
						remainder = strings.TrimPrefix(remainder, "a ")
						remainder = strings.TrimPrefix(remainder, "an ")
						remainder = strings.TrimPrefix(remainder, "the ")
						humanDesc = fmt.Sprintf("Manages %s in F5 Distributed Cloud.", remainder)
						matched = true
						break
					}
				}
			}

			if !matched {
				// Pattern 4: Description is already decent but needs "Manages" prefix
				// If it doesn't start with a verb, add "Manages a X resource" prefix
				if !StartsWithVerb(desc) {
					// Use the description as the capability explanation
					capability := ExtractCapability(desc)
					if capability != "" {
						humanDesc = fmt.Sprintf("Manages a %s resource in F5 Distributed Cloud for %s.", humanName, capability)
					} else {
						humanDesc = fmt.Sprintf("Manages a %s resource in F5 Distributed Cloud. %s", humanName, desc)
					}
				} else {
					// If description already looks good, just ensure it ends properly
					if !strings.HasSuffix(desc, ".") {
						desc = desc + "."
					}
					humanDesc = desc
				}
			}
		}
	}

	// Return human-readable description only
	// Note: AI metadata was previously added here for the terraform-schema-ai tool,
	// but that tool was removed in commit 0c45bb3 as unused. The metadata was
	// polluting Terraform Registry documentation for human users.
	return humanDesc
}

// GenerateCapabilityDescriptionOnly generates a capability-based description for a resource.
// It returns only the description without AI metadata (metadata is added by the caller).
func GenerateCapabilityDescriptionOnly(resourceName, humanName, rawDesc string) string {
	// Resource-specific capability mappings for common F5 XC resources
	capabilities := map[string]string{
		// Sites
		"securemesh_site":    "deploying secure mesh edge sites with distributed security capabilities",
		"securemesh_site_v2": "deploying secure mesh edge sites with enhanced security and networking features",
		"aws_vpc_site":       "deploying F5 sites within AWS VPC environments",
		"azure_vnet_site":    "deploying F5 sites within Azure Virtual Network environments",
		"gcp_vpc_site":       "deploying F5 sites within Google Cloud VPC environments",
		"aws_tgw_site":       "deploying F5 sites connected via AWS Transit Gateway",
		"voltstack_site":     "deploying Volterra stack sites for edge computing",
		"virtual_site":       "creating logical groupings of sites based on labels and selectors",

		// Load Balancing
		"http_loadbalancer": "load balancing HTTP/HTTPS traffic with advanced routing and security",
		"tcp_loadbalancer":  "load balancing TCP traffic across origin pools",
		"udp_loadbalancer":  "load balancing UDP traffic across origin pools",
		"dns_load_balancer": "intelligent DNS-based load balancing across multiple endpoints",
		"cdn_loadbalancer":  "content delivery and edge caching with load balancing",
		"origin_pool":       "defining backend server pools for load balancer targets",
		"healthcheck":       "monitoring backend server health and availability",
		"route":             "defining traffic routing rules for load balancers",

		// Security
		"app_firewall":                   "web application firewall (WAF) protection",
		"service_policy":                 "defining service-level access control and security policies",
		"network_firewall":               "network-level firewall rules and security controls",
		"rate_limiter":                   "protecting services from traffic spikes and DDoS attacks",
		"bot_defense_app_infrastructure": "bot detection and mitigation capabilities",
		"malicious_user_mitigation":      "identifying and blocking malicious user behavior",
		"waf_exclusion_policy":           "excluding specific requests from WAF inspection",

		// Networking
		"network_connector": "connecting networks across sites and cloud providers",
		"virtual_network":   "creating isolated virtual network segments",
		"cloud_connect":     "establishing connectivity to cloud provider networks",
		"cloud_link":        "linking F5 sites to cloud provider infrastructure",
		"bgp":               "BGP routing configuration for network connectivity",
		"ip_prefix_set":     "defining IP address prefix lists for network policies",
		"network_interface": "configuring network interfaces on sites",

		// DNS
		"dns_zone":              "DNS zone management and configuration",
		"dns_domain":            "DNS domain registration and management",
		"dns_lb_pool":           "DNS load balancer endpoint pools",
		"dns_lb_health_check":   "health monitoring for DNS load balanced endpoints",
		"dns_compliance_checks": "DNS security and compliance verification",

		// Kubernetes
		"k8s_cluster":              "Kubernetes cluster integration and management",
		"virtual_k8s":              "virtual Kubernetes cluster deployment",
		"k8s_cluster_role":         "Kubernetes RBAC cluster role definitions",
		"k8s_cluster_role_binding": "Kubernetes RBAC cluster role bindings",
		"k8s_pod_security_policy":  "Kubernetes pod security policy enforcement",
		"container_registry":       "container image registry configuration",

		// Authentication & Secrets
		"authentication":    "authentication methods and identity provider integration",
		"cloud_credentials": "cloud provider credential management for site deployment",
		"api_credential":    "API credential management for service authentication",
		"token":             "API token generation and management",
		"secret_policy":     "secret access policies and controls",

		// Certificates
		"certificate":       "TLS/SSL certificate management",
		"certificate_chain": "certificate chain configuration for TLS",
		"trusted_ca_list":   "trusted certificate authority list management",

		// Monitoring & Logging
		"log_receiver":        "log collection and forwarding configuration",
		"global_log_receiver": "global log aggregation settings",
		"alert_policy":        "alerting rules and notification policies",
		"alert_receiver":      "alert notification endpoints",

		// API Security
		"api_definition": "API schema and endpoint definitions for security",
		"api_discovery":  "automatic API endpoint discovery and inventory",
		"api_testing":    "API testing and validation capabilities",
		"api_crawler":    "API endpoint crawling and discovery",

		// Organization
		"namespace":      "logical namespace isolation for resources",
		"tenant":         "tenant configuration and management",
		"role":           "role-based access control definitions",
		"allowed_tenant": "tenant access permissions and restrictions",
	}

	// Check if we have a specific capability mapping
	if capability, ok := capabilities[resourceName]; ok {
		return fmt.Sprintf("Manages a %s resource in F5 Distributed Cloud for %s.", humanName, capability)
	}

	// Try to extract capability from raw description
	capability := ExtractCapability(rawDesc)
	if capability != "" {
		return fmt.Sprintf("Manages a %s resource in F5 Distributed Cloud for %s.", humanName, capability)
	}

	// Default fallback
	return fmt.Sprintf("Manages a %s resource in F5 Distributed Cloud.", humanName)
}

// ExtractCapability tries to extract a meaningful capability phrase
// from technical descriptions.
func ExtractCapability(desc string) string {
	lowerDesc := strings.ToLower(desc)

	// Remove common technical prefixes
	prefixes := []string{
		"shape of the ",
		"shape of ",
		"specification for ",
		"configuration for ",
		"defines the ",
		"defines a ",
		"represents the ",
		"represents a ",
		"the ",
		"a ",
		"an ",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(lowerDesc, prefix) {
			desc = desc[len(prefix):]
			lowerDesc = strings.ToLower(desc)
			break
		}
	}

	// Remove common technical suffixes
	suffixes := []string{
		" specification",
		" spec",
		" configuration",
		" config",
		" object",
		" definition",
	}
	for _, suffix := range suffixes {
		if strings.HasSuffix(lowerDesc, suffix) {
			desc = desc[:len(desc)-len(suffix)]
			break
		}
	}

	desc = strings.TrimSpace(desc)

	// If what remains is meaningful (more than just a name), use it
	if len(desc) > 5 && !strings.Contains(strings.ToLower(desc), "shape") {
		// Convert to lowercase capability phrase
		return strings.ToLower(desc) + " configuration"
	}

	return ""
}

// StartsWithVerb checks if the description starts with an action verb.
func StartsWithVerb(desc string) bool {
	verbs := []string{
		"manages", "creates", "configures", "defines", "sets", "establishes",
		"provisions", "deploys", "enables", "allows", "provides", "supports",
		"manage", "create", "configure", "define", "set", "establish",
		"provision", "deploy", "enable", "allow", "provide", "support",
	}
	lowerDesc := strings.ToLower(desc)
	for _, verb := range verbs {
		if strings.HasPrefix(lowerDesc, verb+" ") || strings.HasPrefix(lowerDesc, verb+"s ") {
			return true
		}
	}
	return false
}

// AddOneOfConstraint adds a OneOf constraint hint to the description with recommended default.
// The format is AI-friendly: [OneOf: field1, field2; Default: recommended_field]
func AddOneOfConstraint(desc string, oneOfFields []string) string {
	return AddOneOfConstraintWithGroup(desc, "", oneOfFields)
}

// AddOneOfConstraintWithGroup adds a OneOf constraint with group name for better AI parsing.
func AddOneOfConstraintWithGroup(desc string, groupName string, oneOfFields []string) string {
	if len(oneOfFields) < 2 {
		return desc
	}

	// Format fields
	quotedFields := make([]string, len(oneOfFields))
	for i, f := range oneOfFields {
		quotedFields[i] = f
	}

	// Determine recommended default choice using AI-friendly heuristics
	defaultChoice := DetermineOneOfDefault(groupName, oneOfFields)

	// Build AI-friendly constraint marker
	var constraint string
	if defaultChoice != "" {
		constraint = fmt.Sprintf("[OneOf: %s; Default: %s]", strings.Join(quotedFields, ", "), defaultChoice)
	} else {
		constraint = fmt.Sprintf("[OneOf: %s]", strings.Join(quotedFields, ", "))
	}

	// Add constraint at the beginning of description
	if desc == "" {
		return constraint
	}
	return constraint + " " + desc
}

// DetermineOneOfDefault determines the recommended default for a OneOf group.
// This helps AI agents make informed decisions without trial-and-error.
func DetermineOneOfDefault(groupName string, fields []string) string {
	// Common patterns for F5 XC resources - these are the safe, recommended defaults
	defaultPatterns := map[string]string{
		"advertise_choice":           "advertise_on_public_default_vip",
		"loadbalancer_type":          "https_auto_cert",
		"hash_policy_choice":         "round_robin",
		"waf_choice":                 "disable_waf",
		"challenge_type":             "no_challenge",
		"rate_limit_choice":          "disable_rate_limit",
		"service_policy_choice":      "no_service_policies",
		"tls_choice":                 "no_tls",
		"bot_defense_choice":         "disable_bot_defense",
		"api_definition_choice":      "disable_api_definition",
		"api_discovery_choice":       "disable_api_discovery",
		"ip_reputation_choice":       "disable_ip_reputation",
		"malware_protection":         "disable_malware_protection",
		"client_side_defense_choice": "disable_client_side_defense",
	}

	// Check if we have a known pattern for this group
	if groupName != "" {
		if defaultVal, ok := defaultPatterns[groupName]; ok {
			// Verify the default is actually in the fields
			for _, f := range fields {
				if f == defaultVal {
					return defaultVal
				}
			}
		}
	}

	// Heuristics for common patterns
	for _, f := range fields {
		// Prefer "default" variants (e.g., advertise_on_public_default_vip)
		if strings.Contains(f, "default") {
			return f
		}
	}

	for _, f := range fields {
		// Prefer "no_" or "disable_" for optional security features
		if strings.HasPrefix(f, "no_") || strings.HasPrefix(f, "disable_") {
			return f
		}
	}

	// No clear default - don't recommend one
	return ""
}
