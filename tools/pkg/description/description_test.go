// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package description

import (
	"strings"
	"testing"
)

func TestClean(t *testing.T) {
	tests := []struct {
		name     string
		desc     string
		path     string
		expected string
	}{
		{
			name:     "empty input",
			desc:     "",
			path:     "",
			expected: "",
		},
		{
			name:     "removes example section",
			desc:     "A description Example: some example text",
			path:     "",
			expected: "A description",
		},
		{
			name:     "removes validation rules",
			desc:     "A description Validation Rules: must be > 0",
			path:     "",
			expected: "A description",
		},
		{
			name:     "removes x-example annotation",
			desc:     `A description x-example: "test-value"`,
			path:     "",
			expected: "A description",
		},
		{
			name:     "removes x-required annotation",
			desc:     "A description x-required ",
			path:     "",
			expected: "A description",
		},
		{
			name:     "removes ves.io annotations",
			desc:     "A field ves.io.schema.rules.string.max_len: 256",
			path:     "",
			expected: "A field",
		},
		{
			name:     "removes Required YES/NO",
			desc:     "A field Required: YES more text",
			path:     "",
			expected: "A field more text",
		},
		{
			name:     "normalizes empty message description",
			desc:     "Empty. This can be used for messages where no values are needed",
			path:     "",
			expected: "Enable this option",
		},
		{
			name:     "normalizes shape of specification",
			desc:     "Shape of the foobar specification",
			path:     "",
			expected: "Configuration for foobar",
		},
		{
			name:     "normalizes whitespace",
			desc:     "A   description\n\nwith   spaces",
			path:     "",
			expected: "A description with spaces",
		},
		{
			name:     "escapes quotes",
			desc:     `A "quoted" description`,
			path:     "",
			expected: "A 'quoted' description",
		},
		{
			name:     "removes double periods",
			desc:     "A sentence. . More text",
			path:     "",
			expected: "A sentence. More text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Clean(tt.desc, tt.path)
			if result != tt.expected {
				t.Errorf("Clean(%q, %q) = %q, want %q", tt.desc, tt.path, result, tt.expected)
			}
		})
	}
}

func TestExtractDefault(t *testing.T) {
	tests := []struct {
		name         string
		desc         string
		wantDefault  interface{}
		wantContains string // substring the cleaned desc should contain (or empty)
	}{
		{
			name:        "empty input",
			desc:        "",
			wantDefault: nil,
		},
		{
			name:        "defaults to numeric",
			desc:        "Timeout value. Defaults to 30000ms.",
			wantDefault: "30000ms",
		},
		{
			name:        "defaults to boolean",
			desc:        "Enable feature. Defaults to true.",
			wantDefault: "true",
		},
		{
			name:        "default value is string",
			desc:        "Path prefix. Default value is /graphql.",
			wantDefault: "/graphql",
		},
		{
			name:        "default is numeric",
			desc:        "Retry count. Default is 10.",
			wantDefault: "10",
		},
		{
			name:        "skips NONE",
			desc:        "Mode. Defaults to NONE.",
			wantDefault: nil,
		},
		{
			name:        "skips INVALID",
			desc:        "Mode. Default value is INVALID.",
			wantDefault: nil,
		},
		{
			name:        "no default present",
			desc:        "A plain description with no default.",
			wantDefault: nil,
		},
		{
			name:         "removes default text from description",
			desc:         "Timeout value. Defaults to 5s. More info.",
			wantDefault:  "5s",
			wantContains: "More info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDefault, gotDesc := ExtractDefault(tt.desc)
			if tt.wantDefault == nil && gotDefault != nil {
				t.Errorf("ExtractDefault(%q): want nil default, got %v", tt.desc, gotDefault)
			} else if tt.wantDefault != nil {
				if gotDefault == nil {
					t.Errorf("ExtractDefault(%q): want default %v, got nil", tt.desc, tt.wantDefault)
				} else if gotDefault != tt.wantDefault {
					t.Errorf("ExtractDefault(%q): want default %v, got %v", tt.desc, tt.wantDefault, gotDefault)
				}
			}
			if tt.wantContains != "" && !strings.Contains(gotDesc, tt.wantContains) {
				t.Errorf("ExtractDefault(%q): cleaned desc %q should contain %q", tt.desc, gotDesc, tt.wantContains)
			}
		})
	}
}

func TestFormatDefault(t *testing.T) {
	tests := []struct {
		name     string
		desc     string
		defVal   interface{}
		expected string
	}{
		{
			name:     "nil default",
			desc:     "A description",
			defVal:   nil,
			expected: "A description",
		},
		{
			name:     "string default",
			desc:     "A description",
			defVal:   "true",
			expected: "A description. Defaults to `true`.",
		},
		{
			name:     "numeric default",
			desc:     "Timeout.",
			defVal:   "30s",
			expected: "Timeout. Defaults to `30s`.",
		},
		{
			name:     "skips INVALID",
			desc:     "Mode",
			defVal:   "INVALID",
			expected: "Mode",
		},
		{
			name:     "skips zero",
			desc:     "Count",
			defVal:   "0",
			expected: "Count",
		},
		{
			name:     "desc already ends with period",
			desc:     "A description.",
			defVal:   "foo",
			expected: "A description. Defaults to `foo`.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDefault(tt.desc, tt.defVal)
			if result != tt.expected {
				t.Errorf("FormatDefault(%q, %v) = %q, want %q", tt.desc, tt.defVal, result, tt.expected)
			}
		})
	}
}

func TestFormatEnum(t *testing.T) {
	tests := []struct {
		name     string
		desc     string
		enums    []interface{}
		contains []string // substrings the result must contain
	}{
		{
			name:  "empty enums",
			desc:  "A description",
			enums: nil,
			contains: []string{"A description"},
		},
		{
			name:  "single enum",
			desc:  "Mode.",
			enums: []interface{}{"ACTIVE"},
			contains: []string{
				"[Enum: ACTIVE]",
				"The only possible value is `ACTIVE`.",
			},
		},
		{
			name:  "multiple enums",
			desc:  "Security level",
			enums: []interface{}{"LOW", "MEDIUM", "HIGH"},
			contains: []string{
				"[Enum: LOW|MEDIUM|HIGH]",
				"Possible values are `LOW`, `MEDIUM`, `HIGH`.",
				"Security level",
			},
		},
		{
			name:  "empty desc with enums",
			desc:  "",
			enums: []interface{}{"A", "B"},
			contains: []string{
				"[Enum: A|B]",
				"Possible values are `A`, `B`.",
			},
		},
		{
			name:  "skips empty string values",
			desc:  "Test",
			enums: []interface{}{"", "VALID"},
			contains: []string{
				"[Enum: VALID]",
				"The only possible value is `VALID`.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatEnum(tt.desc, tt.enums)
			for _, sub := range tt.contains {
				if !strings.Contains(result, sub) {
					t.Errorf("FormatEnum(%q, %v) = %q, missing %q", tt.desc, tt.enums, result, sub)
				}
			}
		})
	}
}

func TestStartsWithVerb(t *testing.T) {
	tests := []struct {
		desc     string
		expected bool
	}{
		{"Manages a resource", true},
		{"Creates a new item", true},
		{"Configures the settings", true},
		{"Defines the schema", true},
		{"Deploys the application", true},
		{"Enables the feature", true},
		{"Allows access", true},
		{"Provides support", true},
		{"A plain description", false},
		{"The resource name", false},
		{"", false},
		{"Some random text", false},
		{"manage resources", true},
		{"create items", true},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := StartsWithVerb(tt.desc)
			if result != tt.expected {
				t.Errorf("StartsWithVerb(%q) = %v, want %v", tt.desc, result, tt.expected)
			}
		})
	}
}

func TestDetermineOneOfDefault(t *testing.T) {
	tests := []struct {
		name      string
		groupName string
		fields    []string
		expected  string
	}{
		{
			name:      "known pattern advertise_choice",
			groupName: "advertise_choice",
			fields:    []string{"advertise_on_public", "advertise_on_public_default_vip", "do_not_advertise"},
			expected:  "advertise_on_public_default_vip",
		},
		{
			name:      "known pattern not in fields",
			groupName: "waf_choice",
			fields:    []string{"enable_waf", "custom_waf"},
			expected:  "",
		},
		{
			name:      "heuristic default variant",
			groupName: "unknown_choice",
			fields:    []string{"custom_option", "default_option", "advanced_option"},
			expected:  "default_option",
		},
		{
			name:      "heuristic no_ prefix",
			groupName: "unknown_choice",
			fields:    []string{"enable_feature", "no_feature"},
			expected:  "no_feature",
		},
		{
			name:      "heuristic disable_ prefix",
			groupName: "unknown_choice",
			fields:    []string{"enable_it", "disable_it"},
			expected:  "disable_it",
		},
		{
			name:      "no clear default",
			groupName: "unknown_choice",
			fields:    []string{"option_a", "option_b"},
			expected:  "",
		},
		{
			name:      "empty group name with default heuristic",
			groupName: "",
			fields:    []string{"custom", "default_value"},
			expected:  "default_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetermineOneOfDefault(tt.groupName, tt.fields)
			if result != tt.expected {
				t.Errorf("DetermineOneOfDefault(%q, %v) = %q, want %q", tt.groupName, tt.fields, result, tt.expected)
			}
		})
	}
}

func TestAddOneOfConstraint(t *testing.T) {
	tests := []struct {
		name     string
		desc     string
		fields   []string
		contains []string
	}{
		{
			name:     "fewer than 2 fields returns unchanged",
			desc:     "Original",
			fields:   []string{"only_one"},
			contains: []string{"Original"},
		},
		{
			name:   "two fields adds constraint prefix",
			desc:   "A description.",
			fields: []string{"option_a", "option_b"},
			contains: []string{
				"[OneOf: option_a, option_b]",
				"A description.",
			},
		},
		{
			name:   "empty desc with fields",
			desc:   "",
			fields: []string{"a", "b"},
			contains: []string{
				"[OneOf: a, b]",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddOneOfConstraint(tt.desc, tt.fields)
			for _, sub := range tt.contains {
				if !strings.Contains(result, sub) {
					t.Errorf("AddOneOfConstraint(%q, %v) = %q, missing %q", tt.desc, tt.fields, result, sub)
				}
			}
		})
	}
}

func TestAddOneOfConstraintWithGroup(t *testing.T) {
	tests := []struct {
		name      string
		desc      string
		groupName string
		fields    []string
		contains  []string
	}{
		{
			name:      "with known group and default",
			desc:      "Choose advertise method.",
			groupName: "advertise_choice",
			fields:    []string{"advertise_on_public", "advertise_on_public_default_vip"},
			contains: []string{
				"[OneOf: advertise_on_public, advertise_on_public_default_vip; Default: advertise_on_public_default_vip]",
				"Choose advertise method.",
			},
		},
		{
			name:      "fewer than 2 fields unchanged",
			desc:      "Original",
			groupName: "test",
			fields:    []string{"single"},
			contains:  []string{"Original"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddOneOfConstraintWithGroup(tt.desc, tt.groupName, tt.fields)
			for _, sub := range tt.contains {
				if !strings.Contains(result, sub) {
					t.Errorf("AddOneOfConstraintWithGroup(%q, %q, %v) = %q, missing %q",
						tt.desc, tt.groupName, tt.fields, result, sub)
				}
			}
		})
	}
}

func TestExtractCapability(t *testing.T) {
	tests := []struct {
		name     string
		desc     string
		expected string
	}{
		{
			name:     "empty input",
			desc:     "",
			expected: "",
		},
		{
			name:     "shape of prefix",
			desc:     "Shape of the widget specification",
			expected: "widget configuration",
		},
		{
			name:     "represents prefix with suffix",
			desc:     "Represents the network connector object",
			expected: "network connector configuration",
		},
		{
			name:     "too short after stripping",
			desc:     "The foo",
			expected: "",
		},
		{
			name:     "contains shape after stripping",
			desc:     "Shape of shape",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractCapability(tt.desc)
			if result != tt.expected {
				t.Errorf("ExtractCapability(%q) = %q, want %q", tt.desc, result, tt.expected)
			}
		})
	}
}

func TestTransformResourceDescription(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		rawDesc  string
		contains string
	}{
		{
			name:     "empty description",
			resource: "http_loadbalancer",
			rawDesc:  "",
			contains: "Manages a",
		},
		{
			name:     "shape of pattern",
			resource: "origin_pool",
			rawDesc:  "Shape of the origin_pool specification",
			contains: "Manages a",
		},
		{
			name:     "create verb pattern",
			resource: "test_resource",
			rawDesc:  "Creates a test endpoint for validation",
			contains: "Manages",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TransformResourceDescription(tt.resource, tt.rawDesc)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("TransformResourceDescription(%q, %q) = %q, want to contain %q",
					tt.resource, tt.rawDesc, result, tt.contains)
			}
		})
	}
}
