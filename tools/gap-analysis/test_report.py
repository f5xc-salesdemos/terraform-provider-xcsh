# ruff: noqa: INP001, S101, D102, PLR2004, D415, RUF015
# Copyright (c) 2026 Robin Mordasiewicz. MIT License.

"""Tests for generate_report.py - Gap Analysis Report Generator.

Tests use synthetic data to verify report generation logic without
requiring real repository access.
"""

from __future__ import annotations

import pytest
from generate_report import (
    GAP_DEFINITIONS,
    compute_priority_score,
    create_gap_items,
    generate_markdown_report,
)

# =========================================================================
# Test 1: compute_priority_score
# =========================================================================


class TestComputePriorityScore:
    """Verify priority score calculation."""

    def test_high_impact_low_effort(self) -> None:
        """(3,3,1) = 6.0"""
        assert compute_priority_score(3, 3, 1) == 6.0

    def test_high_impact_medium_effort(self) -> None:
        """(3,3,2) = 3.0"""
        assert compute_priority_score(3, 3, 2) == 3.0

    def test_low_impact_high_effort(self) -> None:
        """(1,1,3) ~= 0.667"""
        assert compute_priority_score(1, 1, 3) == pytest.approx(0.667, abs=0.001)

    def test_returns_float(self) -> None:
        """Return type is always float."""
        result = compute_priority_score(2, 2, 1)
        assert isinstance(result, float)

    def test_equal_impact_and_effort(self) -> None:
        """(2,2,2) = 2.0"""
        assert compute_priority_score(2, 2, 2) == 2.0


# =========================================================================
# Test 2: create_gap_items
# =========================================================================


# Synthetic extension map for testing
SAMPLE_EXT_MAP: dict[str, str] = {
    "x-f5xc-description-medium": "CONSUMED",
    "x-f5xc-server-default": "CONSUMED",
    "x-f5xc-constraints": "DEFINED_UNUSED",
    "x-f5xc-conflicts-with": "PARSED_NOT_RENDERED",
    "x-f5xc-danger-level": "UNKNOWN",
    "x-f5xc-namespace-scope": "NOT_EMITTED",
}

# Synthetic coverage data for testing
SAMPLE_COVERAGE_DATA: list[dict] = [
    {
        "resource_name": "http_loadbalancer",
        "domain_file": "virtual.json",
        "domain_name": "virtual",
        "category": "Load Balancers",
        "tier": "basic",
        "total_fields": 50,
        "fields_with_constraints": 10,
        "fields_with_description_medium": 30,
        "fields_with_required_for": 5,
        "fields_with_server_default": 8,
        "fields_with_conflicts_with": 3,
        "fields_with_enum": 12,
    },
    {
        "resource_name": "tcp_loadbalancer",
        "domain_file": "virtual.json",
        "domain_name": "virtual",
        "category": "Load Balancers",
        "tier": "basic",
        "total_fields": 30,
        "fields_with_constraints": 5,
        "fields_with_description_medium": 20,
        "fields_with_required_for": 2,
        "fields_with_server_default": 4,
        "fields_with_conflicts_with": 1,
        "fields_with_enum": 6,
    },
    {
        "resource_name": "origin_pool",
        "domain_file": "views.json",
        "domain_name": "views",
        "category": "Networking",
        "tier": "basic",
        "total_fields": 40,
        "fields_with_constraints": 8,
        "fields_with_description_medium": 25,
        "fields_with_required_for": 3,
        "fields_with_server_default": 6,
        "fields_with_conflicts_with": 2,
        "fields_with_enum": 10,
    },
]


class TestCreateGapItems:
    """Verify gap item creation from hardcoded definitions."""

    def test_returns_list(self) -> None:
        """Return type is a list."""
        items = create_gap_items(SAMPLE_EXT_MAP, SAMPLE_COVERAGE_DATA)
        assert isinstance(items, list)

    def test_returns_ten_items(self) -> None:
        """GAP_DEFINITIONS has 10 entries, so 10 items returned."""
        items = create_gap_items(SAMPLE_EXT_MAP, SAMPLE_COVERAGE_DATA)
        assert len(items) == 10

    def test_items_sorted_by_priority_descending(self) -> None:
        """Items are sorted by priority_score from high to low."""
        items = create_gap_items(SAMPLE_EXT_MAP, SAMPLE_COVERAGE_DATA)
        scores = [item["priority_score"] for item in items]
        assert scores == sorted(scores, reverse=True)

    def test_each_item_has_required_fields(self) -> None:
        """Each gap item has title, description, repo, scores, key."""
        items = create_gap_items(SAMPLE_EXT_MAP, SAMPLE_COVERAGE_DATA)
        required = {
            "title",
            "description",
            "repo",
            "user_impact",
            "downstream_reach",
            "effort",
            "priority_score",
            "key",
        }
        for item in items:
            for field in required:
                assert field in item, (
                    f"Missing field {field} in item {item.get('key', '?')}"
                )

    def test_constraints_gap_present(self) -> None:
        """x-f5xc-constraints gap item exists."""
        items = create_gap_items(SAMPLE_EXT_MAP, SAMPLE_COVERAGE_DATA)
        keys = {item["key"] for item in items}
        assert "x-f5xc-constraints" in keys

    def test_orphan_data_sources_present(self) -> None:
        """orphan-data-sources gap item exists."""
        items = create_gap_items(SAMPLE_EXT_MAP, SAMPLE_COVERAGE_DATA)
        keys = {item["key"] for item in items}
        assert "orphan-data-sources" in keys

    def test_constraints_priority_score(self) -> None:
        """x-f5xc-constraints has score (3+3)/2 = 3.0."""
        items = create_gap_items(SAMPLE_EXT_MAP, SAMPLE_COVERAGE_DATA)
        constraints = [i for i in items if i["key"] == "x-f5xc-constraints"][0]
        assert constraints["priority_score"] == 3.0

    def test_orphan_data_sources_priority_score(self) -> None:
        """orphan-data-sources has score (1+1)/1 = 2.0."""
        items = create_gap_items(SAMPLE_EXT_MAP, SAMPLE_COVERAGE_DATA)
        orphan = [i for i in items if i["key"] == "orphan-data-sources"][0]
        assert orphan["priority_score"] == 2.0


# =========================================================================
# Test 3: GAP_DEFINITIONS constant
# =========================================================================


class TestGapDefinitions:
    """Verify the hardcoded GAP_DEFINITIONS dict."""

    def test_has_ten_entries(self) -> None:
        """GAP_DEFINITIONS should have exactly 10 entries."""
        assert len(GAP_DEFINITIONS) == 10

    def test_all_keys_present(self) -> None:
        """All 10 expected gap keys are present."""
        expected_keys = {
            "x-f5xc-constraints",
            "x-f5xc-conflicts-with",
            "x-f5xc-minimum-configuration",
            "x-f5xc-namespace-scope",
            "enum-validators",
            "x-f5xc-danger-level",
            "x-f5xc-required-for",
            "x-f5xc-best-practices",
            "orphan-data-sources",
            "dependency-dead-code",
        }
        assert set(GAP_DEFINITIONS.keys()) == expected_keys

    def test_each_definition_has_required_fields(self) -> None:
        """Each definition has title, description, repo, user_impact, downstream_reach, effort."""
        required = {
            "title",
            "description",
            "repo",
            "user_impact",
            "downstream_reach",
            "effort",
        }
        for key, defn in GAP_DEFINITIONS.items():
            for field in required:
                assert field in defn, f"Missing {field} in GAP_DEFINITIONS[{key}]"


# =========================================================================
# Test 4: generate_markdown_report
# =========================================================================


# Synthetic operation data for testing
SAMPLE_OP_DATA: dict[str, int] = {
    "total_operations": 100,
    "ops_with_operation_metadata": 60,
    "ops_with_danger_level": 20,
    "ops_with_confirmation_required": 15,
    "ops_with_side_effects": 10,
    "ops_with_required_fields": 25,
}


class TestGenerateMarkdownReport:
    """Verify markdown report generation."""

    @pytest.fixture
    def gap_items(self) -> list[dict]:
        return create_gap_items(SAMPLE_EXT_MAP, SAMPLE_COVERAGE_DATA)

    @pytest.fixture
    def report(self, gap_items: list[dict]) -> str:
        return generate_markdown_report(
            SAMPLE_EXT_MAP,
            SAMPLE_COVERAGE_DATA,
            gap_items,
            SAMPLE_OP_DATA,
        )

    def test_report_is_string(self, report: str) -> None:
        """Report is a string."""
        assert isinstance(report, str)

    def test_report_not_empty(self, report: str) -> None:
        """Report is not empty."""
        assert len(report) > 0

    # --- Section 1: Executive Summary ---
    def test_has_executive_summary(self, report: str) -> None:
        """Report contains Executive Summary header."""
        assert "## 1. Executive Summary" in report

    # --- Section 2: Extension Consumption Matrix ---
    def test_has_extension_matrix(self, report: str) -> None:
        """Report contains Extension Consumption Matrix header."""
        assert "## 2. Extension Consumption Matrix" in report

    def test_extension_names_in_matrix(self, report: str) -> None:
        """Extension names appear in the consumption matrix."""
        assert "x-f5xc-description-medium" in report
        assert "x-f5xc-constraints" in report

    # --- Section 3: Resource-Level Drill-Downs ---
    def test_has_resource_drilldowns(self, report: str) -> None:
        """Report contains Resource-Level Drill-Downs header."""
        assert "## 3. Resource-Level Drill-Downs" in report

    def test_resource_names_in_drilldowns(self, report: str) -> None:
        """Resource names appear in the drill-downs section."""
        assert "http_loadbalancer" in report
        assert "tcp_loadbalancer" in report
        assert "origin_pool" in report

    # --- Section 4: Validator Opportunity Analysis ---
    def test_has_validator_analysis(self, report: str) -> None:
        """Report contains Validator Opportunity Analysis header."""
        assert "## 4. Validator Opportunity Analysis" in report

    # --- Section 5: Schema Fidelity Findings ---
    def test_has_schema_fidelity(self, report: str) -> None:
        """Report contains Schema Fidelity Findings header."""
        assert "## 5. Schema Fidelity Findings" in report

    # --- Section 6: Spec Enrichment Priorities ---
    def test_has_enrichment_priorities(self, report: str) -> None:
        """Report contains Spec Enrichment Priorities header."""
        assert "## 6. Spec Enrichment Priorities" in report

    # --- Section 7: Downstream Impact Assessment ---
    def test_has_downstream_impact(self, report: str) -> None:
        """Report contains Downstream Impact Assessment header."""
        assert "## 7. Downstream Impact Assessment" in report

    # --- Section 8: Prioritized Action Items ---
    def test_has_action_items(self, report: str) -> None:
        """Report contains Prioritized Action Items header."""
        assert "## 8. Prioritized Action Items" in report

    # --- GitHub Issue Templates ---
    def test_has_issue_templates(self, report: str) -> None:
        """Report contains GitHub Issue Templates section."""
        assert "## GitHub Issue Templates" in report

    def test_issue_templates_contain_gap_titles(
        self, report: str, gap_items: list[dict]
    ) -> None:
        """Issue templates reference gap item titles."""
        for item in gap_items:
            assert item["title"] in report

    def test_issue_templates_contain_repos(self, report: str) -> None:
        """Issue templates contain repo references."""
        assert "terraform-provider-f5xc" in report

    def test_categories_grouped(self, report: str) -> None:
        """Resource drill-downs group by category."""
        assert "Load Balancers" in report
        assert "Networking" in report

    def test_status_counts_in_summary(self, report: str) -> None:
        """Executive summary includes extension status counts."""
        assert "CONSUMED" in report
        assert "DEFINED_UNUSED" in report
