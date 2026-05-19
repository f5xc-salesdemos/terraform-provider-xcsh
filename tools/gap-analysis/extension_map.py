# ruff: noqa: INP001
# Copyright (c) 2026 Robin Mordasiewicz. MIT License.

"""Extension Consumption Map Builder.

Parses Go source files and enriched API specs to classify each x-f5xc-*
extension into one of six statuses:

    CONSUMED          - Parsed from spec AND rendered into Terraform output
    PARSED_NOT_RENDERED - Parsed into Go structs but not rendered in templates
    DOMAIN_ONLY       - Only meaningful at the domain/index level, not per-field
    DEFINED_UNUSED    - Defined in Go struct but never assigned or rendered
    UNKNOWN           - Registered in extension_constants.py but not in Go struct
    NOT_EMITTED       - Registered but never appears in any output spec JSON
"""

from __future__ import annotations

import json
import re
from pathlib import Path

# =============================================================================
# Hardcoded classification sets
# =============================================================================

CONSUMED_EXTENSIONS = frozenset(
    {
        # SP-1: Parsed into openapi.Schema struct fields
        # SP-4: Constraints pipeline fixed — validators emitted in generated code
        "x-f5xc-constraints",
        "x-f5xc-conflicts-with",
        # SP-6: Required fields + server defaults wired into schema generation
        "x-f5xc-required-for",
        "x-f5xc-server-default",
        # Pre-existing: consumed since initial generator
        "x-f5xc-description-medium",
        "x-f5xc-description-short",
        "x-f5xc-description",
        "x-f5xc-complexity",
        "x-f5xc-example",
        "x-f5xc-requires-tier",
        "x-f5xc-category",
        "x-f5xc-is-preview",
        # SP-7: Operation extensions injected into docs
        "x-f5xc-danger-level",
        "x-f5xc-side-effects",
        "x-f5xc-confirmation-required",
        "x-f5xc-best-practices",
    }
)

PARSED_NOT_RENDERED = frozenset(
    {
        # Parsed into openapi.Schema but not yet rendered into output
        "x-f5xc-recommended-value",
        "x-f5xc-recommended-oneof-variant",
        "x-f5xc-minimum-configuration",
        "x-f5xc-validation",
        "x-f5xc-defaults",
        "x-f5xc-conditions",
        "x-f5xc-deprecated",
        "x-f5xc-completion",
        "x-f5xc-display-name",
        "x-f5xc-examples",
        "x-f5xc-required-for-operations",
        "x-f5xc-displayorder",
        "x-f5xc-uniqueness",
        "x-f5xc-terraform-resource",
        "x-f5xc-operation-metadata",
        "x-f5xc-required-fields",
    }
)

DOMAIN_ONLY = frozenset(
    {
        "x-f5xc-use-cases",
        "x-f5xc-related-domains",
        "x-f5xc-icon",
        "x-f5xc-logo-svg",
        "x-f5xc-doc-section",
        "x-f5xc-cli-domain",
        "x-f5xc-cli-metadata",
        "x-f5xc-glossary",
        "x-f5xc-guided-workflows",
        "x-f5xc-acronyms",
        "x-f5xc-critical-resources",
        "x-f5xc-primary-resources",
    }
)

# Regex to find json:"x-f5xc-..." tags in Go struct definitions
_JSON_TAG_RE = re.compile(r'json:"(x-f5xc-[^"]+)"')

# Regex to find TerraformAttribute struct literal field assignments
# Matches patterns like:  FieldName: value  or  FieldName:  value
_STRUCT_FIELD_ASSIGN_RE = re.compile(
    r"^\s+([A-Z][A-Za-z0-9]+)\s*:\s+\S",
)

# Regex to find fields referenced inside Go template strings ({{.FieldName}})
_TEMPLATE_FIELD_RE = re.compile(r"\{\{[^}]*\.([A-Z][A-Za-z0-9]+)")


# =============================================================================
# Step 2: Extract x-f5xc-* json tags from Go struct definitions
# =============================================================================


def extract_go_struct_xf5xc_fields(go_file: Path) -> set[str]:
    """Extract x-f5xc-* extension names from json struct tags in a Go file.

    Scans for ``json:"x-f5xc-..."`` tags in Go struct definitions and
    returns the set of extension names found.

    Args:
        go_file: Path to the Go source file.

    Returns:
        Set of x-f5xc-* extension name strings.
    """
    text = go_file.read_text(encoding="utf-8")
    return set(_JSON_TAG_RE.findall(text))


# =============================================================================
# Step 4: Detect TerraformAttribute field usage
# =============================================================================


def find_terraform_attribute_field_usage(
    gen_file: Path,
) -> tuple[set[str], set[str]]:
    """Find TerraformAttribute fields assigned in struct literals vs templates.

    Parses the Go generator file to discover:
    1. Fields assigned in ``TerraformAttribute{ ... }`` struct literals
    2. Fields referenced inside Go template strings via ``{{.FieldName}}``

    Args:
        gen_file: Path to the Go generator file.

    Returns:
        Tuple of (assigned_fields, rendered_fields).
    """
    text = gen_file.read_text(encoding="utf-8")

    # --- Find assigned fields in TerraformAttribute struct literals ---
    # We look for the struct literal block and extract field names
    assigned: set[str] = set()
    in_struct = False
    brace_depth = 0

    for line in text.splitlines():
        # Detect start of TerraformAttribute struct literal
        if "TerraformAttribute{" in line:
            in_struct = True
            brace_depth = 0
            # Count braces on this line
            brace_depth += line.count("{") - line.count("}")
            continue

        if in_struct:
            brace_depth += line.count("{") - line.count("}")
            if brace_depth <= 0:
                in_struct = False
                continue

            m = _STRUCT_FIELD_ASSIGN_RE.match(line)
            if m:
                assigned.add(m.group(1))

    # --- Find rendered fields in Go template strings ---
    rendered: set[str] = set()

    # Go raw string literals use backticks
    # We look for template directives inside them
    for m in _TEMPLATE_FIELD_RE.finditer(text):
        rendered.add(m.group(1))

    return assigned, rendered


# =============================================================================
# Helper: Get emitted extensions from spec JSONs
# =============================================================================


def get_emitted_extensions(specs_dir: Path) -> set[str]:
    """Scan output spec JSON files for x-f5xc-* keys actually emitted.

    Args:
        specs_dir: Directory containing enriched API spec JSON files.

    Returns:
        Set of x-f5xc-* extension names found in the spec files.
    """
    pattern = re.compile(r'"(x-f5xc-[^"]+)"')
    extensions: set[str] = set()

    for json_file in sorted(specs_dir.glob("*.json")):
        text = json_file.read_text(encoding="utf-8")
        extensions.update(pattern.findall(text))

    return extensions


# =============================================================================
# Helper: Get registered extensions from extension_constants.py
# =============================================================================


def get_registered_extensions(specs_root: Path) -> set[str]:
    """Read VALID_X_F5XC_EXTENSIONS from extension_constants.py.

    Parses the Python source file to extract all x-f5xc-* string literals
    assigned as constants (lines matching ``X_F5XC_... = "x-f5xc-..."``).

    Args:
        specs_root: Root of the api-specs-enriched repository.

    Returns:
        Set of registered x-f5xc-* extension names.
    """
    constants_file = specs_root / "scripts" / "utils" / "extension_constants.py"
    text = constants_file.read_text(encoding="utf-8")

    # Match constant definitions like: X_F5XC_FOO = "x-f5xc-foo"
    pattern = re.compile(r'=\s*"(x-f5xc-[^"]+)"')
    return set(pattern.findall(text))


# =============================================================================
# Step 6: Build the full extension map
# =============================================================================


def build_extension_map(
    provider_root: Path,
    specs_root: Path,
) -> dict[str, str]:
    """Classify every known x-f5xc-* extension into a consumption status.

    The classification logic:
    1. If the extension is in CONSUMED_EXTENSIONS -> CONSUMED
    2. If the extension is in PARSED_NOT_RENDERED -> PARSED_NOT_RENDERED
    3. If the extension is in DOMAIN_ONLY -> DOMAIN_ONLY
    4. If the extension is defined in Go structs but not in any of the above
       hardcoded sets -> DEFINED_UNUSED
    5. If the extension is registered but not emitted in specs -> NOT_EMITTED
    6. Otherwise -> UNKNOWN

    Args:
        provider_root: Root of the terraform-provider-f5xc repository.
        specs_root: Root of the api-specs-enriched repository.

    Returns:
        Dict mapping extension name to status string.
    """
    go_file = provider_root / "tools" / "generate-all-schemas.go"
    types_file = provider_root / "tools" / "pkg" / "openapi" / "types.go"
    transform_file = provider_root / "tools" / "transform-docs.go"
    specs_dir = specs_root / "docs" / "specifications" / "api"

    # Gather data — scan multiple Go files since SP-2 modularized the generator
    go_fields = extract_go_struct_xf5xc_fields(go_file)
    if types_file.exists():
        go_fields |= extract_go_struct_xf5xc_fields(types_file)
    if transform_file.exists():
        go_fields |= extract_go_struct_xf5xc_fields(transform_file)
    registered = get_registered_extensions(specs_root)
    emitted = get_emitted_extensions(specs_dir)

    # Union of all known extensions
    all_extensions = go_fields | registered | emitted

    result: dict[str, str] = {}

    for ext in sorted(all_extensions):
        if ext in CONSUMED_EXTENSIONS:
            result[ext] = "CONSUMED"
        elif ext in PARSED_NOT_RENDERED:
            result[ext] = "PARSED_NOT_RENDERED"
        elif ext in DOMAIN_ONLY:
            result[ext] = "DOMAIN_ONLY"
        elif ext in go_fields:
            # Defined in Go struct but not in any hardcoded consumption set
            result[ext] = "DEFINED_UNUSED"
        elif ext in registered and ext not in emitted:
            # Registered in extension_constants.py but never emitted
            result[ext] = "NOT_EMITTED"
        else:
            # Registered and/or emitted but not consumed by the provider
            result[ext] = "UNKNOWN"

    return result


# =============================================================================
# CLI entry point
# =============================================================================

if __name__ == "__main__":
    import sys

    _MIN_ARGS_FOR_SPECS = 3
    provider = (
        Path(sys.argv[1])
        if len(sys.argv) > 1
        else Path("/workspace/terraform-provider-f5xc")
    )
    specs = (
        Path(sys.argv[2])
        if len(sys.argv) >= _MIN_ARGS_FOR_SPECS
        else Path("/workspace/api-specs-enriched")
    )

    ext_map = build_extension_map(provider, specs)

    print(json.dumps(ext_map, indent=2))  # noqa: T201
