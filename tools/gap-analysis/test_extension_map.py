# ruff: noqa: INP001, S101, D102, PLR2004
# Copyright (c) 2026 Robin Mordasiewicz. MIT License.

"""Tests for extension_map.py - Extension Consumption Map Builder.

Tests run against real files in the terraform-provider-f5xc and
api-specs-enriched repositories.
"""

from __future__ import annotations

from pathlib import Path

import pytest
from extension_map import (
    CONSUMED_EXTENSIONS,
    DOMAIN_ONLY,
    PARSED_NOT_RENDERED,
    build_extension_map,
    extract_go_struct_xf5xc_fields,
    find_terraform_attribute_field_usage,
    get_emitted_extensions,
    get_registered_extensions,
)

PROVIDER_ROOT = Path("/workspace/terraform-provider-f5xc")
SPECS_ROOT = Path("/workspace/api-specs-enriched")
GO_FILE = PROVIDER_ROOT / "tools" / "generate-all-schemas.go"
SPECS_DIR = SPECS_ROOT / "docs" / "specifications" / "api"


# =========================================================================
# Step 1: Test Go struct field extraction
# =========================================================================


class TestExtractSchemaDefinitionFields:
    """Verify extract_go_struct_xf5xc_fields finds x-f5xc-* json tags."""

    def test_extract_finds_known_extensions(self) -> None:
        """Known x-f5xc-* fields present in SchemaDefinition are returned."""
        fields = extract_go_struct_xf5xc_fields(GO_FILE)

        expected = {
            "x-f5xc-constraints",
            "x-f5xc-conflicts-with",
            "x-f5xc-server-default",
            "x-f5xc-description-medium",
            "x-f5xc-required-for",
            "x-f5xc-recommended-value",
            "x-f5xc-minimum-configuration",
            "x-f5xc-recommended-oneof-variant",
        }
        for ext in expected:
            assert ext in fields, f"Expected {ext} in extracted fields"

    def test_extract_excludes_non_f5xc_extensions(self) -> None:
        """x-displayname and x-ves-example are NOT x-f5xc-* extensions."""
        fields = extract_go_struct_xf5xc_fields(GO_FILE)

        assert "x-displayname" not in fields
        assert "x-ves-example" not in fields

    def test_extract_returns_set(self) -> None:
        """Return type is a set of strings."""
        fields = extract_go_struct_xf5xc_fields(GO_FILE)
        assert isinstance(fields, set)
        assert all(isinstance(f, str) for f in fields)

    def test_extract_all_start_with_prefix(self) -> None:
        """Every returned field starts with x-f5xc-."""
        fields = extract_go_struct_xf5xc_fields(GO_FILE)
        for field in fields:
            assert field.startswith("x-f5xc-"), f"{field} missing prefix"


# =========================================================================
# Step 3: Test TerraformAttribute field usage detection
# =========================================================================


class TestFindTerraformAttributeUsage:
    """Verify detection of fields assigned in struct literals vs templates."""

    def test_assigned_fields_include_known(self) -> None:
        """ConflictsWith, MinimumConfigRequired, RecommendedValue are assigned."""
        assigned, _rendered = find_terraform_attribute_field_usage(GO_FILE)

        assert "ConflictsWith" in assigned
        assert "MinimumConfigRequired" in assigned
        assert "RecommendedValue" in assigned

    def test_rendered_fields_include_description(self) -> None:
        """Description is used inside Go template strings."""
        _assigned, rendered = find_terraform_attribute_field_usage(GO_FILE)

        assert "Description" in rendered

    def test_returns_two_sets(self) -> None:
        """Return value is a tuple of two sets."""
        result = find_terraform_attribute_field_usage(GO_FILE)
        assert isinstance(result, tuple)
        assert len(result) == 2
        assigned, rendered = result
        assert isinstance(assigned, set)
        assert isinstance(rendered, set)


# =========================================================================
# Step 5: Test full extension status classification
# =========================================================================


class TestBuildExtensionMap:
    """Verify the full classification pipeline."""

    @pytest.fixture
    def ext_map(self) -> dict[str, str]:
        return build_extension_map(PROVIDER_ROOT, SPECS_ROOT)

    def test_consumed_description_medium(self, ext_map: dict[str, str]) -> None:
        assert ext_map.get("x-f5xc-description-medium") == "CONSUMED"

    def test_consumed_server_default(self, ext_map: dict[str, str]) -> None:
        assert ext_map.get("x-f5xc-server-default") == "CONSUMED"

    def test_defined_unused_constraints(self, ext_map: dict[str, str]) -> None:
        assert ext_map.get("x-f5xc-constraints") == "DEFINED_UNUSED"

    def test_parsed_not_rendered_conflicts_with(self, ext_map: dict[str, str]) -> None:
        assert ext_map.get("x-f5xc-conflicts-with") == "PARSED_NOT_RENDERED"

    def test_unknown_danger_level(self, ext_map: dict[str, str]) -> None:
        assert ext_map.get("x-f5xc-danger-level") == "UNKNOWN"

    def test_not_emitted_namespace_scope(self, ext_map: dict[str, str]) -> None:
        assert ext_map.get("x-f5xc-namespace-scope") == "NOT_EMITTED"

    def test_all_values_are_valid_statuses(self, ext_map: dict[str, str]) -> None:
        valid = {
            "CONSUMED",
            "PARSED_NOT_RENDERED",
            "DOMAIN_ONLY",
            "DEFINED_UNUSED",
            "UNKNOWN",
            "NOT_EMITTED",
        }
        for ext, status in ext_map.items():
            assert status in valid, f"{ext} has invalid status {status}"

    def test_map_is_not_empty(self, ext_map: dict[str, str]) -> None:
        assert len(ext_map) > 0


# =========================================================================
# Helper function tests
# =========================================================================


class TestGetEmittedExtensions:
    """Verify scanning of output spec JSONs."""

    def test_returns_set(self) -> None:
        result = get_emitted_extensions(SPECS_DIR)
        assert isinstance(result, set)

    def test_finds_known_emitted(self) -> None:
        result = get_emitted_extensions(SPECS_DIR)
        # We know these exist in the spec JSONs from our exploration
        assert "x-f5xc-server-default" in result
        assert "x-f5xc-description-medium" in result

    def test_all_start_with_prefix(self) -> None:
        result = get_emitted_extensions(SPECS_DIR)
        for ext in result:
            assert ext.startswith("x-f5xc-"), f"{ext} missing prefix"


class TestGetRegisteredExtensions:
    """Verify reading the extension registry."""

    def test_returns_set(self) -> None:
        result = get_registered_extensions(SPECS_ROOT)
        assert isinstance(result, set)

    def test_finds_known_registered(self) -> None:
        result = get_registered_extensions(SPECS_ROOT)
        assert "x-f5xc-server-default" in result
        assert "x-f5xc-description-medium" in result
        assert "x-f5xc-danger-level" in result

    def test_all_start_with_prefix(self) -> None:
        result = get_registered_extensions(SPECS_ROOT)
        for ext in result:
            assert ext.startswith("x-f5xc-"), f"{ext} missing prefix"


class TestConstants:
    """Verify the hardcoded classification sets are consistent."""

    def test_no_overlap_consumed_parsed(self) -> None:
        assert CONSUMED_EXTENSIONS.isdisjoint(PARSED_NOT_RENDERED)

    def test_no_overlap_consumed_domain(self) -> None:
        assert CONSUMED_EXTENSIONS.isdisjoint(DOMAIN_ONLY)

    def test_no_overlap_parsed_domain(self) -> None:
        assert PARSED_NOT_RENDERED.isdisjoint(DOMAIN_ONLY)
