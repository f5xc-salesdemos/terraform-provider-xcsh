// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package constraints

import "regexp"

// Parsed represents extracted x-f5xc-constraints data.
type Parsed struct {
	MinLength int
	MaxLength int
	Pattern   string
	MinItems  int
	MaxItems  int
	Minimum   int
	Maximum   int
}

// Parse extracts constraint data from x-f5xc-constraints map.
// Returns nil for nil input, low-confidence, or non-deterministic constraints.
// Accepts deterministic at top level (actual spec format) or in metadata (legacy).
// Only uses constraints where deterministic==true OR metadata.confidence >= 0.9.
func Parse(raw map[string]interface{}) *Parsed {
	if raw == nil {
		return nil
	}

	deterministic := false
	confidence := 0.0

	// Check top-level deterministic (actual spec format since v2.1.80+)
	if d, ok := raw["deterministic"].(bool); ok {
		deterministic = d
	}

	// Also check metadata for confidence and legacy deterministic
	if meta, ok := raw["metadata"].(map[string]interface{}); ok {
		if d, ok := meta["deterministic"].(bool); ok && d {
			deterministic = true
		}
		if c, ok := meta["confidence"].(float64); ok {
			confidence = c
		}
	}

	if !deterministic && confidence < 0.9 {
		return nil
	}

	p := &Parsed{}

	// String constraints
	if v, ok := raw["minLength"].(float64); ok {
		p.MinLength = int(v)
	}
	if v, ok := raw["maxLength"].(float64); ok {
		p.MaxLength = int(v)
	}
	if v, ok := raw["pattern"].(string); ok {
		if _, err := regexp.Compile(v); err == nil {
			p.Pattern = v
		}
	}

	// List/array constraints
	if v, ok := raw["minItems"].(float64); ok {
		p.MinItems = int(v)
	}
	if v, ok := raw["maxItems"].(float64); ok {
		p.MaxItems = int(v)
	}

	// Numeric constraints
	if v, ok := raw["minimum"].(float64); ok {
		p.Minimum = int(v)
	}
	if v, ok := raw["maximum"].(float64); ok {
		p.Maximum = int(v)
	}

	return p
}
