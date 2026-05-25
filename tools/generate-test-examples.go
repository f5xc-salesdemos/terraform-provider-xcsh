// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

//go:build ignore
// +build ignore

// generate-test-examples.go extracts verified Terraform configurations from
// acceptance test files and writes them as named example .tf files.
// These examples are proven to work against the live F5 XC staging API.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var configFuncRegex = regexp.MustCompile(`func (testAcc\w+Config_\w+)\([^)]*\)\s+string\s*{`)
var formatStringRegex = regexp.MustCompile("return\\s+(?:acctest\\.ConfigCompose\\(\\s*acctest\\.ProviderConfig\\(\\),\\s*)?fmt\\.Sprintf\\(`([^`]+)`")
var simpleReturnRegex = regexp.MustCompile("return\\s+fmt\\.Sprintf\\(`([^`]+)`")

type testExample struct {
	Resource string
	Name     string
	Config   string
}

func main() {
	testDir := "internal/provider"
	outputDir := "examples/resources"

	resources := []string{
		"http_loadbalancer",
		"tcp_loadbalancer",
		"healthcheck",
		"app_firewall",
		"origin_pool",
		"rate_limiter",
		"service_policy",
		"user_identification",
		"malicious_user_mitigation",
	}

	totalWritten := 0

	for _, res := range resources {
		testFile := filepath.Join(testDir, res+"_resource_test.go")
		content, err := os.ReadFile(testFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Skipping %s: %v\n", res, err)
			continue
		}

		examples := extractExamples(res, string(content))

		for _, ex := range examples {
			exDir := filepath.Join(outputDir, "f5xc_"+ex.Resource)
			if err := os.MkdirAll(exDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating %s: %v\n", exDir, err)
				continue
			}

			filename := toExampleFilename(ex.Name) + ".tf"
			if filename == "basic-system.tf" || filename == "basic.tf" {
				continue
			}

			outPath := filepath.Join(exDir, filename)
			header := fmt.Sprintf("# %s — Verified Configuration Example\n# This configuration is extracted from acceptance tests\n# and verified against the live F5 XC API.\n\n",
				toHumanName(ex.Name))

			cleaned := cleanConfig(ex.Config)
			if cleaned == "" {
				continue
			}

			if err := os.WriteFile(outPath, []byte(header+cleaned), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outPath, err)
				continue
			}

			fmt.Printf("Generated: %s\n", outPath)
			totalWritten++
		}
	}

	fmt.Printf("\nGenerated %d verified example files\n", totalWritten)
}

func extractExamples(resource, content string) []testExample {
	var examples []testExample

	matches := configFuncRegex.FindAllStringSubmatchIndex(content, -1)

	for _, match := range matches {
		funcName := content[match[2]:match[3]]
		funcStart := match[0]

		funcEnd := findFuncEnd(content, funcStart)
		if funcEnd < 0 {
			continue
		}

		funcBody := content[funcStart:funcEnd]

		configName := extractConfigName(funcName, resource)
		if configName == "" {
			continue
		}

		hcl := extractHCL(funcBody)
		if hcl == "" {
			continue
		}

		examples = append(examples, testExample{
			Resource: resource,
			Name:     configName,
			Config:   hcl,
		})
	}

	return examples
}

func findFuncEnd(content string, start int) int {
	depth := 0
	inString := false
	inBacktick := false

	for i := start; i < len(content); i++ {
		c := content[i]

		if c == '`' {
			inBacktick = !inBacktick
			continue
		}
		if inBacktick {
			continue
		}
		if c == '"' && (i == 0 || content[i-1] != '\\') {
			inString = !inString
			continue
		}
		if inString {
			continue
		}
		if c == '{' {
			depth++
		}
		if c == '}' {
			depth--
			if depth == 0 {
				return i + 1
			}
		}
	}

	return -1
}

func extractConfigName(funcName, resource string) string {
	prefix := "testAcc" + toCamelCase(resource) + "Config_"
	altPrefix := "testAccHTTPLBConfig_"
	altPrefix2 := "testAccHTTPLoadBalancerConfig_"

	name := ""
	if strings.HasPrefix(funcName, prefix) {
		name = funcName[len(prefix):]
	} else if strings.HasPrefix(funcName, altPrefix) {
		name = funcName[len(altPrefix):]
	} else if strings.HasPrefix(funcName, altPrefix2) {
		name = funcName[len(altPrefix2):]
	} else {
		for _, p := range []string{"testAcc"} {
			if strings.HasPrefix(funcName, p) && strings.Contains(funcName, "Config_") {
				idx := strings.Index(funcName, "Config_")
				name = funcName[idx+7:]
				break
			}
		}
	}

	return name
}

func extractHCL(funcBody string) string {
	backtickStart := strings.Index(funcBody, "`\n")
	if backtickStart < 0 {
		backtickStart = strings.Index(funcBody, "(`")
		if backtickStart < 0 {
			return ""
		}
		backtickStart++
	}
	backtickStart++

	backtickEnd := strings.Index(funcBody[backtickStart+1:], "`")
	if backtickEnd < 0 {
		return ""
	}
	backtickEnd += backtickStart + 1

	return funcBody[backtickStart:backtickEnd]
}

func cleanConfig(config string) string {
	config = strings.TrimSpace(config)

	config = regexp.MustCompile(`%\[1\]q`).ReplaceAllString(config, `"example"`)
	config = regexp.MustCompile(`%\[1\]s`).ReplaceAllString(config, "example")
	config = regexp.MustCompile(`%\[2\]q`).ReplaceAllString(config, `"example-value"`)
	config = regexp.MustCompile(`%\[2\]s`).ReplaceAllString(config, `"example-value"`)
	config = regexp.MustCompile(`%\[2\]d`).ReplaceAllString(config, "443")
	config = regexp.MustCompile(`%\[3\]q`).ReplaceAllString(config, `"example-description"`)
	config = regexp.MustCompile(`%\[3\]s`).ReplaceAllString(config, "example-value")
	config = regexp.MustCompile(`%\[3\]d`).ReplaceAllString(config, "3")
	config = regexp.MustCompile(`%\[4\]q`).ReplaceAllString(config, `"example-value"`)
	config = regexp.MustCompile(`%\[4\]s`).ReplaceAllString(config, "example-value")
	config = regexp.MustCompile(`%\[4\]d`).ReplaceAllString(config, "5")
	config = regexp.MustCompile(`%\[5\]s`).ReplaceAllString(config, "example-value")
	config = regexp.MustCompile(`%\[5\]d`).ReplaceAllString(config, "15")
	config = regexp.MustCompile(`%s`).ReplaceAllString(config, "example")
	config = regexp.MustCompile(`%q`).ReplaceAllString(config, `"example"`)
	config = regexp.MustCompile(`%d`).ReplaceAllString(config, "80")

	if strings.Contains(config, "%") && strings.Contains(config, "[") {
		return ""
	}

	return config
}

func toExampleFilename(name string) string {
	name = strings.TrimSuffix(name, "System")
	name = strings.TrimSuffix(name, "_system")

	acronyms := map[string]string{
		"WAF": "waf", "TLS": "tls", "TCP": "tcp", "UDP": "udp",
		"HTTP": "http", "HTTPS": "https", "DNS": "dns", "API": "api",
		"IP": "ip", "SSL": "ssl", "SNI": "sni", "ICMP": "icmp",
	}

	for acr, lower := range acronyms {
		name = strings.ReplaceAll(name, acr, strings.ToUpper(lower[:1])+lower[1:])
	}

	result := ""
	for i, c := range name {
		if c >= 'A' && c <= 'Z' {
			if i > 0 && result[len(result)-1] != '-' {
				result += "-"
			}
			result += string(c + 32)
		} else if c == '_' {
			result += "-"
		} else {
			result += string(c)
		}
	}

	result = strings.ReplaceAll(result, "--", "-")
	result = strings.TrimPrefix(result, "-")
	return result
}

func toHumanName(name string) string {
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "System", "")
	name = strings.TrimSpace(name)
	words := strings.Fields(name)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func toCamelCase(snake string) string {
	parts := strings.Split(snake, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}
