# ruff: noqa: INP001
# Copyright (c) 2026 Robin Mordasiewicz. MIT License.

"""Gap Analysis Report Generator.

Combines extension consumption map and spec coverage data into a
prioritized 8-section markdown gap report with GitHub Issue templates.
"""

from __future__ import annotations

from collections import defaultdict

# =============================================================================
# Hardcoded gap definitions (10 items)
# =============================================================================

GAP_DEFINITIONS: dict[str, dict] = {
    "x-f5xc-constraints": {
        "title": "Consume x-f5xc-constraints to generate Terraform validators",
        "description": (
            "The x-f5xc-constraints extension is parsed into Go structs but "
            "never rendered into Terraform validation logic. Consuming this "
            "extension would auto-generate validators for min/max, pattern, "
            "and custom constraints."
        ),
        "repo": "terraform-provider-f5xc",
        "user_impact": 3,
        "downstream_reach": 3,
        "effort": 2,
    },
    "x-f5xc-conflicts-with": {
        "title": "Implement ConflictsWith validators",
        "description": (
            "The x-f5xc-conflicts-with extension is parsed but not rendered. "
            "Generating ConflictsWith validators would prevent users from "
            "setting mutually exclusive fields."
        ),
        "repo": "terraform-provider-f5xc",
        "user_impact": 3,
        "downstream_reach": 2,
        "effort": 2,
    },
    "x-f5xc-minimum-configuration": {
        "title": "Consume for Required field accuracy",
        "description": (
            "The x-f5xc-minimum-configuration extension defines the minimal "
            "set of fields needed for a valid resource. Consuming it would "
            "improve Required field accuracy in the Terraform schema."
        ),
        "repo": "terraform-provider-f5xc",
        "user_impact": 2,
        "downstream_reach": 2,
        "effort": 2,
    },
    "x-f5xc-namespace-scope": {
        "title": "Emit and replace hardcoded maps",
        "description": (
            "The x-f5xc-namespace-scope extension is registered but not "
            "emitted in any spec file. Once emitted, it can replace the "
            "hardcoded namespace scope maps in the Terraform provider."
        ),
        "repo": "both",
        "user_impact": 2,
        "downstream_reach": 3,
        "effort": 2,
    },
    "enum-validators": {
        "title": "Generate stringvalidator.OneOf from enum fields",
        "description": (
            "Enum fields in the spec define valid string values. Generating "
            "stringvalidator.OneOf validators from these would catch invalid "
            "values at plan time instead of apply time."
        ),
        "repo": "terraform-provider-f5xc",
        "user_impact": 3,
        "downstream_reach": 1,
        "effort": 2,
    },
    "x-f5xc-danger-level": {
        "title": "Consume operation-level extensions for destruction warnings",
        "description": (
            "The x-f5xc-danger-level extension marks operations that may "
            "cause service disruption. Consuming it would add destruction "
            "warnings and confirmation prompts to dangerous operations."
        ),
        "repo": "terraform-provider-f5xc",
        "user_impact": 2,
        "downstream_reach": 2,
        "effort": 3,
    },
    "x-f5xc-required-for": {
        "title": "Render in generated code",
        "description": (
            "The x-f5xc-required-for extension is parsed but not rendered. "
            "Rendering it would document which fields are required for "
            "specific operations (create, update, etc.)."
        ),
        "repo": "terraform-provider-f5xc",
        "user_impact": 2,
        "downstream_reach": 2,
        "effort": 2,
    },
    "x-f5xc-best-practices": {
        "title": "Fix struct shape mismatch and embed in docs",
        "description": (
            "The x-f5xc-best-practices extension has a struct shape mismatch "
            "between the spec definition and the Go struct. Fixing this and "
            "embedding best practices in generated documentation would guide "
            "users toward optimal configurations."
        ),
        "repo": "terraform-provider-f5xc",
        "user_impact": 1,
        "downstream_reach": 3,
        "effort": 2,
    },
    "orphan-data-sources": {
        "title": "Audit 49 orphan data source files",
        "description": (
            "There are approximately 49 data source files that are not "
            "connected to any primary resource in index.json. These should "
            "be audited and either connected or removed."
        ),
        "repo": "terraform-provider-f5xc",
        "user_impact": 1,
        "downstream_reach": 1,
        "effort": 1,
    },
    "dependency-dead-code": {
        "title": "Wire up index-derived dependency map or remove",
        "description": (
            "The dependency map derived from index.json is computed but "
            "never wired into the provider. It should either be connected "
            "to resource ordering logic or removed as dead code."
        ),
        "repo": "terraform-provider-f5xc",
        "user_impact": 1,
        "downstream_reach": 1,
        "effort": 2,
    },
}


# =============================================================================
# Function 1: compute_priority_score
# =============================================================================


def compute_priority_score(
    user_impact: int,
    downstream_reach: int,
    effort: int,
) -> float:
    """Compute a priority score for a gap item.

    Formula: (user_impact + downstream_reach) / effort

    Args:
        user_impact: User impact rating (1-3).
        downstream_reach: Downstream reach rating (1-3).
        effort: Effort estimate (1-3).

    Returns:
        Priority score as a float.
    """
    return (user_impact + downstream_reach) / effort


# =============================================================================
# Function 2: create_gap_items
# =============================================================================


def create_gap_items(
    ext_map: dict[str, str],  # noqa: ARG001
    coverage_data: list[dict],  # noqa: ARG001
) -> list[dict]:
    """Create prioritized gap items from hardcoded definitions.

    Uses GAP_DEFINITIONS to create gap items with computed priority
    scores, sorted by priority_score descending.

    Args:
        ext_map: Extension consumption map (extension -> status).
        coverage_data: List of resource coverage dicts.

    Returns:
        List of gap item dicts sorted by priority_score descending.
    """
    items: list[dict] = []

    for key, defn in GAP_DEFINITIONS.items():
        score = compute_priority_score(
            defn["user_impact"],
            defn["downstream_reach"],
            defn["effort"],
        )
        items.append(
            {
                "key": key,
                "title": defn["title"],
                "description": defn["description"],
                "repo": defn["repo"],
                "user_impact": defn["user_impact"],
                "downstream_reach": defn["downstream_reach"],
                "effort": defn["effort"],
                "priority_score": score,
            }
        )

    items.sort(key=lambda x: x["priority_score"], reverse=True)
    return items


# =============================================================================
# Function 3: generate_markdown_report
# =============================================================================


def generate_markdown_report(
    ext_map: dict[str, str],
    coverage_data: list[dict],
    gap_items: list[dict],
    op_data: dict[str, int],
) -> str:
    """Generate the complete 8-section markdown gap analysis report.

    Args:
        ext_map: Extension consumption map (extension -> status).
        coverage_data: List of resource coverage dicts from spec_coverage.
        gap_items: Prioritized gap items from create_gap_items.
        op_data: Operation-level extension coverage from
                 audit_operation_extensions.

    Returns:
        Complete markdown report as a string.
    """
    sections: list[str] = []

    # Title
    sections.append("# F5 XC Terraform Provider - Gap Analysis Report\n")

    # --- Section 1: Executive Summary ---
    sections.append(_section_executive_summary(ext_map, gap_items))

    # --- Section 2: Extension Consumption Matrix ---
    sections.append(_section_extension_matrix(ext_map))

    # --- Section 3: Resource-Level Drill-Downs ---
    sections.append(_section_resource_drilldowns(coverage_data))

    # --- Section 4: Validator Opportunity Analysis ---
    sections.append(_section_validator_analysis(coverage_data))

    # --- Section 5: Schema Fidelity Findings ---
    sections.append(_section_schema_fidelity(ext_map, coverage_data, op_data))

    # --- Section 6: Spec Enrichment Priorities ---
    sections.append(_section_enrichment_priorities())

    # --- Section 7: Downstream Impact Assessment ---
    sections.append(_section_downstream_impact())

    # --- Section 8: Prioritized Action Items ---
    sections.append(_section_action_items(gap_items))

    # --- GitHub Issue Templates ---
    sections.append(_section_issue_templates(gap_items))

    return "\n".join(sections)


# =============================================================================
# Private section generators
# =============================================================================


def _section_executive_summary(
    ext_map: dict[str, str],
    gap_items: list[dict],
) -> str:
    """Generate Section 1: Executive Summary."""
    lines: list[str] = []
    lines.append("## 1. Executive Summary\n")

    # Extension status counts
    status_counts: dict[str, int] = defaultdict(int)
    for status in ext_map.values():
        status_counts[status] += 1

    lines.append("### Extension Status Overview\n")
    lines.append("| Status | Count |")
    lines.append("|--------|-------|")
    lines.extend(
        f"| {status} | {status_counts[status]} |"
        for status in sorted(status_counts.keys())
    )
    lines.append("")

    # Top 10 gaps table
    lines.append("### Top 10 Gaps\n")
    lines.append("| Rank | Gap | Priority | Impact | Reach | Effort |")
    lines.append("|------|-----|----------|--------|-------|--------|")
    for i, item in enumerate(gap_items[:10], 1):
        lines.append(
            f"| {i} | {item['title']} | {item['priority_score']:.1f} "
            f"| {item['user_impact']} | {item['downstream_reach']} "
            f"| {item['effort']} |"
        )
    lines.append("")

    return "\n".join(lines)


def _section_extension_matrix(ext_map: dict[str, str]) -> str:
    """Generate Section 2: Extension Consumption Matrix."""
    lines: list[str] = []
    lines.append("## 2. Extension Consumption Matrix\n")
    lines.append("| Extension | Status |")
    lines.append("|-----------|--------|")
    lines.extend(f"| {ext} | {ext_map[ext]} |" for ext in sorted(ext_map.keys()))
    lines.append("")
    return "\n".join(lines)


def _section_resource_drilldowns(coverage_data: list[dict]) -> str:
    """Generate Section 3: Resource-Level Drill-Downs."""
    lines: list[str] = []
    lines.append("## 3. Resource-Level Drill-Downs\n")

    # Group resources by category
    by_category: dict[str, list[dict]] = defaultdict(list)
    for res in coverage_data:
        category = res.get("category", "Uncategorized")
        by_category[category].append(res)

    for category in sorted(by_category.keys()):
        lines.append(f"### {category}\n")
        lines.append(
            "| Resource | Total Fields | Constraints % | Descriptions % | Enums % |"
        )
        lines.append(
            "|----------|-------------|--------------|---------------|---------|"
        )

        for res in sorted(by_category[category], key=lambda r: r["resource_name"]):
            total = res["total_fields"]
            if total > 0:
                constraints_pct = (res["fields_with_constraints"] / total) * 100
                desc_pct = (res["fields_with_description_medium"] / total) * 100
                enum_pct = (res["fields_with_enum"] / total) * 100
            else:
                constraints_pct = desc_pct = enum_pct = 0.0

            lines.append(
                f"| {res['resource_name']} | {total} | "
                f"{constraints_pct:.0f}% | {desc_pct:.0f}% | "
                f"{enum_pct:.0f}% |"
            )
        lines.append("")

    return "\n".join(lines)


def _section_validator_analysis(coverage_data: list[dict]) -> str:
    """Generate Section 4: Validator Opportunity Analysis."""
    lines: list[str] = []
    lines.append("## 4. Validator Opportunity Analysis\n")

    total_constraints = sum(r["fields_with_constraints"] for r in coverage_data)
    total_enums = sum(r["fields_with_enum"] for r in coverage_data)
    total_conflicts = sum(r["fields_with_conflicts_with"] for r in coverage_data)

    lines.append(
        f"- **Constraint validators available**: {total_constraints} fields "
        f"across {len(coverage_data)} resources"
    )
    lines.append(
        f"- **Enum validators available**: {total_enums} fields with enum "
        f"values that could generate stringvalidator.OneOf"
    )
    lines.append(
        f"- **ConflictsWith validators available**: {total_conflicts} fields "
        f"with conflict declarations"
    )
    lines.append("")

    return "\n".join(lines)


def _section_schema_fidelity(
    ext_map: dict[str, str],
    coverage_data: list[dict],
    op_data: dict[str, int],
) -> str:
    """Generate Section 5: Schema Fidelity Findings."""
    lines: list[str] = []
    lines.append("## 5. Schema Fidelity Findings\n")

    # Count unconsumed extensions
    unconsumed = sum(1 for s in ext_map.values() if s != "CONSUMED")
    total = len(ext_map)

    lines.append(f"- {unconsumed} of {total} extensions are not fully consumed")

    # Count total fields without descriptions
    total_fields = sum(r["total_fields"] for r in coverage_data)
    fields_with_desc = sum(r["fields_with_description_medium"] for r in coverage_data)
    fields_without_desc = total_fields - fields_with_desc
    lines.append(
        f"- {fields_without_desc} fields lack enriched descriptions "
        f"(of {total_fields} total)"
    )

    # Operation-level coverage
    total_ops = op_data.get("total_operations", 0)
    ops_with_danger = op_data.get("ops_with_danger_level", 0)
    if total_ops > 0:
        lines.append(
            f"- {ops_with_danger} of {total_ops} operations have danger "
            f"level annotations ({ops_with_danger * 100 // total_ops}%)"
        )

    lines.append("")
    return "\n".join(lines)


def _section_enrichment_priorities() -> str:
    """Generate Section 6: Spec Enrichment Priorities."""
    lines: list[str] = []
    lines.append("## 6. Spec Enrichment Priorities\n")

    lines.append("| Config File | Purpose | Priority |")
    lines.append("|------------|---------|----------|")
    lines.append("| extension_constants.py | Extension registry | High |")
    lines.append("| index.json | Resource-to-domain map | High |")
    lines.append("| Domain spec files | Per-field extensions | Medium |")
    lines.append("| Operation metadata | Operation-level annotations | Medium |")
    lines.append("")

    return "\n".join(lines)


def _section_downstream_impact() -> str:
    """Generate Section 7: Downstream Impact Assessment."""
    lines: list[str] = []
    lines.append("## 7. Downstream Impact Assessment\n")

    lines.append("| Consumer | Impact Area | Dependency |")
    lines.append("|----------|------------|------------|")
    lines.append("| terraform-provider-f5xc | Schema generation | api-specs-enriched |")
    lines.append(
        "| terraform-provider-f5xc | Validator generation | x-f5xc-constraints, enum |"
    )
    lines.append(
        "| terraform-provider-f5xc | Documentation | "
        "x-f5xc-description-*, x-f5xc-best-practices |"
    )
    lines.append("| api-specs-enriched | Extension emission | extension_constants.py |")
    lines.append("")

    return "\n".join(lines)


def _section_action_items(gap_items: list[dict]) -> str:
    """Generate Section 8: Prioritized Action Items."""
    lines: list[str] = []
    lines.append("## 8. Prioritized Action Items\n")

    lines.append("| Rank | Action | Repo | Priority Score |")
    lines.append("|------|--------|------|---------------|")
    for i, item in enumerate(gap_items, 1):
        lines.append(
            f"| {i} | {item['title']} | {item['repo']} | {item['priority_score']:.1f} |"
        )
    lines.append("")

    return "\n".join(lines)


def _section_issue_templates(gap_items: list[dict]) -> str:
    """Generate GitHub Issue Templates section."""
    lines: list[str] = []
    lines.append("## GitHub Issue Templates\n")

    for item in gap_items:
        lines.append(f"### {item['title']}\n")
        lines.append(f"**Repo**: {item['repo']}")
        lines.append(f"**Priority Score**: {item['priority_score']:.1f}")
        lines.append(
            f"**User Impact**: {item['user_impact']} | "
            f"**Downstream Reach**: {item['downstream_reach']} | "
            f"**Effort**: {item['effort']}\n"
        )
        lines.append(f"{item['description']}\n")

    return "\n".join(lines)


# =============================================================================
# CLI entry point
# =============================================================================

if __name__ == "__main__":
    # When run standalone, generate a report with sample data
    sample_ext_map = {
        "x-f5xc-description-medium": "CONSUMED",
        "x-f5xc-server-default": "CONSUMED",
        "x-f5xc-constraints": "DEFINED_UNUSED",
        "x-f5xc-conflicts-with": "PARSED_NOT_RENDERED",
    }

    sample_coverage: list[dict] = []
    sample_op: dict[str, int] = {
        "total_operations": 0,
        "ops_with_danger_level": 0,
    }

    items = create_gap_items(sample_ext_map, sample_coverage)
    report = generate_markdown_report(sample_ext_map, sample_coverage, items, sample_op)
    print(report)  # noqa: T201
