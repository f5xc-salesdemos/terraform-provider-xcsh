# ruff: noqa: INP001, S101, D102
# Copyright (c) 2026 Robin Mordasiewicz. MIT License.

"""Tests for spec_coverage.py - Spec Coverage Auditor.

Tests run against real files in the api-specs-enriched repository.
"""

from __future__ import annotations

from pathlib import Path

import pytest
from spec_coverage import (
    audit_all_resources,
    audit_operation_extensions,
    audit_resource_coverage,
    build_resource_domain_map,
)

SPECS_DIR = Path("/workspace/api-specs-enriched/docs/specifications/api")


# =========================================================================
# Test 1: build_resource_domain_map
# =========================================================================


class TestBuildResourceDomainMap:
    """Verify build_resource_domain_map reads index.json correctly."""

    @pytest.fixture
    def resource_map(self) -> dict[str, dict]:
        return build_resource_domain_map(SPECS_DIR)

    def test_http_loadbalancer_present(self, resource_map: dict[str, dict]) -> None:
        """http_loadbalancer should map to a .json file."""
        assert "http_loadbalancer" in resource_map
        info = resource_map["http_loadbalancer"]
        assert info["domain_file"].endswith(".json")

    def test_namespace_role_present(self, resource_map: dict[str, dict]) -> None:
        """namespace_role resource should exist in the map."""
        assert "namespace_role" in resource_map

    def test_http_loadbalancer_has_required_fields(
        self, resource_map: dict[str, dict]
    ) -> None:
        """http_loadbalancer entry has domain_file, domain_name, category, schema_components, api_paths, tier."""
        info = resource_map["http_loadbalancer"]
        assert "domain_file" in info
        assert "domain_name" in info
        assert "category" in info
        assert "schema_components" in info
        assert "api_paths" in info
        assert "tier" in info

    def test_schema_components_is_list(self, resource_map: dict[str, dict]) -> None:
        """schema_components should be a list."""
        info = resource_map["http_loadbalancer"]
        assert isinstance(info["schema_components"], list)

    def test_returns_dict(self, resource_map: dict[str, dict]) -> None:
        """Return type is a dict mapping resource names to info dicts."""
        assert isinstance(resource_map, dict)
        assert len(resource_map) > 0

    def test_http_loadbalancer_domain_is_virtual(
        self, resource_map: dict[str, dict]
    ) -> None:
        """http_loadbalancer should be in the virtual domain."""
        info = resource_map["http_loadbalancer"]
        assert info["domain_name"] == "virtual"


# =========================================================================
# Test 2: audit_resource_coverage
# =========================================================================


class TestAuditResourceCoverage:
    """Verify audit_resource_coverage counts extensions per schema component."""

    @pytest.fixture
    def coverage(self) -> dict[str, int]:
        domain_file = SPECS_DIR / "virtual.json"
        # Use views.http_loadbalancer which maps to viewshttp_loadbalancer* schemas
        return audit_resource_coverage(domain_file, ["views.http_loadbalancer"])

    def test_total_fields_greater_than_zero(self, coverage: dict[str, int]) -> None:
        """total_fields should be > 0 for http_loadbalancer."""
        assert coverage["total_fields"] > 0

    def test_returns_required_keys(self, coverage: dict[str, int]) -> None:
        """Coverage dict has all required keys."""
        expected_keys = {
            "total_fields",
            "fields_with_constraints",
            "fields_with_description_medium",
            "fields_with_required_for",
            "fields_with_server_default",
            "fields_with_conflicts_with",
            "fields_with_enum",
        }
        for key in expected_keys:
            assert key in coverage, f"Missing key: {key}"

    def test_all_values_are_ints(self, coverage: dict[str, int]) -> None:
        """All coverage values should be integers."""
        for key, value in coverage.items():
            assert isinstance(value, int), f"{key} is not an int"

    def test_constraints_count_nonnegative(self, coverage: dict[str, int]) -> None:
        """fields_with_constraints should be >= 0."""
        assert coverage["fields_with_constraints"] >= 0

    def test_description_medium_count_nonnegative(
        self, coverage: dict[str, int]
    ) -> None:
        """fields_with_description_medium should be >= 0."""
        assert coverage["fields_with_description_medium"] >= 0


# =========================================================================
# Test 3: audit_operation_extensions
# =========================================================================


class TestAuditOperationExtensions:
    """Verify audit_operation_extensions counts operation-level extensions."""

    @pytest.fixture
    def op_coverage(self) -> dict[str, int]:
        domain_file = SPECS_DIR / "virtual.json"
        return audit_operation_extensions(domain_file)

    def test_total_operations_greater_than_zero(
        self, op_coverage: dict[str, int]
    ) -> None:
        """total_operations should be > 0 for virtual.json."""
        assert op_coverage["total_operations"] > 0

    def test_returns_required_keys(self, op_coverage: dict[str, int]) -> None:
        """Operation coverage dict has all required keys."""
        expected_keys = {
            "total_operations",
            "ops_with_operation_metadata",
            "ops_with_danger_level",
            "ops_with_confirmation_required",
            "ops_with_side_effects",
            "ops_with_required_fields",
        }
        for key in expected_keys:
            assert key in op_coverage, f"Missing key: {key}"

    def test_all_values_are_ints(self, op_coverage: dict[str, int]) -> None:
        """All operation coverage values should be integers."""
        for key, value in op_coverage.items():
            assert isinstance(value, int), f"{key} is not an int"

    def test_operation_metadata_count_nonnegative(
        self, op_coverage: dict[str, int]
    ) -> None:
        """ops_with_operation_metadata should be >= 0."""
        assert op_coverage["ops_with_operation_metadata"] >= 0


# =========================================================================
# Test 4: audit_all_resources
# =========================================================================


class TestAuditAllResources:
    """Verify audit_all_resources iterates and aggregates coverage."""

    @pytest.fixture
    def all_coverage(self) -> list[dict]:
        return audit_all_resources(SPECS_DIR)

    def test_returns_list(self, all_coverage: list[dict]) -> None:
        """Return type is a list."""
        assert isinstance(all_coverage, list)

    def test_result_length_greater_than_zero(self, all_coverage: list[dict]) -> None:
        """Result list should have at least one entry."""
        assert len(all_coverage) > 0

    def test_first_result_has_required_keys(self, all_coverage: list[dict]) -> None:
        """First result has resource_name, domain_file, total_fields."""
        first = all_coverage[0]
        assert "resource_name" in first
        assert "domain_file" in first
        assert "total_fields" in first

    def test_all_results_have_resource_name(self, all_coverage: list[dict]) -> None:
        """Every result dict has a resource_name."""
        for result in all_coverage:
            assert "resource_name" in result

    def test_all_results_have_domain_name(self, all_coverage: list[dict]) -> None:
        """Every result dict has a domain_name."""
        for result in all_coverage:
            assert "domain_name" in result

    def test_all_results_have_category(self, all_coverage: list[dict]) -> None:
        """Every result dict has a category."""
        for result in all_coverage:
            assert "category" in result

    def test_all_results_have_tier(self, all_coverage: list[dict]) -> None:
        """Every result dict has a tier."""
        for result in all_coverage:
            assert "tier" in result
