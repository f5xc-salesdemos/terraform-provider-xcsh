# ruff: noqa: INP001
# Copyright (c) 2026 Robin Mordasiewicz. MIT License.

"""Spec Coverage Auditor.

Audits each registered resource against its enriched API spec to measure
coverage of key x-f5xc-* extensions per field and per operation.

Reads index.json to discover primary resources, then inspects each domain
spec file to count extension presence across schema properties and API
operations.
"""

from __future__ import annotations

import json
from pathlib import Path

# HTTP methods that represent API operations
_HTTP_METHODS = frozenset({"get", "post", "put", "patch", "delete", "head", "options"})


# =============================================================================
# Function 1: Build resource-to-domain map from index.json
# =============================================================================


def build_resource_domain_map(specs_dir: Path) -> dict[str, dict]:
    """Read index.json and map each primary resource to its domain metadata.

    Args:
        specs_dir: Directory containing the enriched API spec JSON files
                   and index.json.

    Returns:
        Dict mapping resource name to a dict with:
            domain_file, domain_name, category, schema_components,
            api_paths, tier.
    """
    index_path = specs_dir / "index.json"
    with index_path.open(encoding="utf-8") as f:
        index_data = json.load(f)

    resource_map: dict[str, dict] = {}

    for spec in index_data.get("specifications", []):
        domain_file = spec.get("file", "")
        domain_name = spec.get("domain", "")
        category = spec.get("x-f5xc-category", "")

        for resource in spec.get("x-f5xc-primary-resources", []):
            name = resource.get("name", "")
            if not name:
                continue

            resource_map[name] = {
                "domain_file": domain_file,
                "domain_name": domain_name,
                "category": category,
                "schema_components": resource.get("schema_components", []),
                "api_paths": resource.get("api_paths", []),
                "tier": resource.get("tier", ""),
            }

    return resource_map


# =============================================================================
# Function 2: Audit schema coverage for a resource
# =============================================================================


def audit_resource_coverage(
    domain_file: Path,
    schema_components: list[str],
) -> dict[str, int]:
    """Count extension coverage across schema properties for given components.

    For each schema component prefix, finds all matching schemas in the
    domain spec file and counts how many properties have each extension.

    A schema component like ``views.http_loadbalancer`` matches all schema
    keys starting with ``viewshttp_loadbalancer`` (dots are removed).

    Args:
        domain_file: Path to the domain spec JSON file.
        schema_components: List of schema component identifiers from
                          index.json (e.g., ["views.http_loadbalancer"]).

    Returns:
        Dict with total_fields, fields_with_constraints,
        fields_with_description_medium, fields_with_required_for,
        fields_with_server_default, fields_with_conflicts_with,
        fields_with_enum.
    """
    with domain_file.open(encoding="utf-8") as f:
        spec_data = json.load(f)

    schemas = spec_data.get("components", {}).get("schemas", {})

    # Convert schema component identifiers to prefixes (remove dots)
    prefixes = [comp.replace(".", "") for comp in schema_components]

    # Collect all properties from matching schemas
    total_fields = 0
    fields_with_constraints = 0
    fields_with_description_medium = 0
    fields_with_required_for = 0
    fields_with_server_default = 0
    fields_with_conflicts_with = 0
    fields_with_enum = 0

    for schema_key, schema_def in schemas.items():
        # Check if this schema matches any of the component prefixes
        if not any(schema_key.startswith(prefix) for prefix in prefixes):
            continue

        properties = schema_def.get("properties", {})
        for prop_def in properties.values():
            total_fields += 1

            if "x-f5xc-constraints" in prop_def:
                fields_with_constraints += 1

            if (
                "x-f5xc-description-medium" in prop_def
                or "x-f5xc-description-short" in prop_def
            ):
                fields_with_description_medium += 1

            if "x-f5xc-required-for" in prop_def:
                fields_with_required_for += 1

            if "x-f5xc-server-default" in prop_def:
                fields_with_server_default += 1

            if "x-f5xc-conflicts-with" in prop_def:
                fields_with_conflicts_with += 1

            if "enum" in prop_def:
                fields_with_enum += 1

    return {
        "total_fields": total_fields,
        "fields_with_constraints": fields_with_constraints,
        "fields_with_description_medium": fields_with_description_medium,
        "fields_with_required_for": fields_with_required_for,
        "fields_with_server_default": fields_with_server_default,
        "fields_with_conflicts_with": fields_with_conflicts_with,
        "fields_with_enum": fields_with_enum,
    }


# =============================================================================
# Function 3: Audit operation-level extensions
# =============================================================================


def audit_operation_extensions(domain_file: Path) -> dict[str, int]:
    """Count operation-level extension coverage in a domain spec file.

    Iterates all API paths and their HTTP methods to count which
    operations have specific x-f5xc-* extensions.

    Args:
        domain_file: Path to the domain spec JSON file.

    Returns:
        Dict with total_operations, ops_with_operation_metadata,
        ops_with_danger_level, ops_with_confirmation_required,
        ops_with_side_effects, ops_with_required_fields.
    """
    with domain_file.open(encoding="utf-8") as f:
        spec_data = json.load(f)

    paths = spec_data.get("paths", {})

    total_operations = 0
    ops_with_operation_metadata = 0
    ops_with_danger_level = 0
    ops_with_confirmation_required = 0
    ops_with_side_effects = 0
    ops_with_required_fields = 0

    for methods in paths.values():
        for method, operation in methods.items():
            if method.lower() not in _HTTP_METHODS:
                continue
            if not isinstance(operation, dict):
                continue

            total_operations += 1

            if "x-f5xc-operation-metadata" in operation:
                ops_with_operation_metadata += 1

            if "x-f5xc-danger-level" in operation:
                ops_with_danger_level += 1

            if "x-f5xc-confirmation-required" in operation:
                ops_with_confirmation_required += 1

            if "x-f5xc-side-effects" in operation:
                ops_with_side_effects += 1

            if "x-f5xc-required-fields" in operation:
                ops_with_required_fields += 1

    return {
        "total_operations": total_operations,
        "ops_with_operation_metadata": ops_with_operation_metadata,
        "ops_with_danger_level": ops_with_danger_level,
        "ops_with_confirmation_required": ops_with_confirmation_required,
        "ops_with_side_effects": ops_with_side_effects,
        "ops_with_required_fields": ops_with_required_fields,
    }


# =============================================================================
# Function 4: Audit all resources
# =============================================================================


def audit_all_resources(specs_dir: Path) -> list[dict]:
    """Audit coverage for all primary resources in the index.

    Iterates the resource-to-domain map, calls audit_resource_coverage
    for each resource, and returns a combined list.

    Args:
        specs_dir: Directory containing the enriched API spec JSON files
                   and index.json.

    Returns:
        List of dicts, each with resource_name, domain_file, domain_name,
        category, tier, plus all coverage fields from
        audit_resource_coverage.
    """
    resource_map = build_resource_domain_map(specs_dir)
    results: list[dict] = []

    for resource_name, info in sorted(resource_map.items()):
        domain_file_path = specs_dir / info["domain_file"]

        if not domain_file_path.exists():
            continue

        coverage = audit_resource_coverage(
            domain_file_path,
            info["schema_components"],
        )

        result = {
            "resource_name": resource_name,
            "domain_file": info["domain_file"],
            "domain_name": info["domain_name"],
            "category": info["category"],
            "tier": info["tier"],
        }
        result.update(coverage)
        results.append(result)

    return results


# =============================================================================
# CLI entry point
# =============================================================================

if __name__ == "__main__":
    import sys

    specs = (
        Path(sys.argv[1])
        if len(sys.argv) > 1
        else Path("/workspace/api-specs-enriched/docs/specifications/api")
    )

    results = audit_all_resources(specs)
    print(json.dumps(results, indent=2))  # noqa: T201
