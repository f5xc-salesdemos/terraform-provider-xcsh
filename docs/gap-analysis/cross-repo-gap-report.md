# F5 XC Terraform Provider - Gap Analysis Report

## 1. Executive Summary

### Extension Status Overview

| Status | Count |
|--------|-------|
| CONSUMED | 16 |
| DEFINED_UNUSED | 4 |
| DOMAIN_ONLY | 12 |
| NOT_EMITTED | 8 |
| PARSED_NOT_RENDERED | 16 |
| UNKNOWN | 1 |

### Top 10 Gaps

| Rank | Gap | Priority | Impact | Reach | Effort |
|------|-----|----------|--------|-------|--------|
| 1 | Consume x-f5xc-constraints to generate Terraform validators | 3.0 | 3 | 3 | 2 |
| 2 | Implement ConflictsWith validators | 2.5 | 3 | 2 | 2 |
| 3 | Emit and replace hardcoded maps | 2.5 | 2 | 3 | 2 |
| 4 | Consume for Required field accuracy | 2.0 | 2 | 2 | 2 |
| 5 | Generate stringvalidator.OneOf from enum fields | 2.0 | 3 | 1 | 2 |
| 6 | Render in generated code | 2.0 | 2 | 2 | 2 |
| 7 | Fix struct shape mismatch and embed in docs | 2.0 | 1 | 3 | 2 |
| 8 | Audit 49 orphan data source files | 2.0 | 1 | 1 | 1 |
| 9 | Consume operation-level extensions for destruction warnings | 1.3 | 2 | 2 | 3 |
| 10 | Wire up index-derived dependency map or remove | 1.0 | 1 | 1 | 2 |

## 2. Extension Consumption Matrix

| Extension | Status |
|-----------|--------|
| x-f5xc-acronyms | DOMAIN_ONLY |
| x-f5xc-api-url | NOT_EMITTED |
| x-f5xc-best-practices | CONSUMED |
| x-f5xc-category | CONSUMED |
| x-f5xc-cli-domain | DOMAIN_ONLY |
| x-f5xc-cli-metadata | DOMAIN_ONLY |
| x-f5xc-completion | PARSED_NOT_RENDERED |
| x-f5xc-complexity | CONSUMED |
| x-f5xc-conditions | PARSED_NOT_RENDERED |
| x-f5xc-confirmation-required | CONSUMED |
| x-f5xc-conflicts-with | CONSUMED |
| x-f5xc-constraints | CONSUMED |
| x-f5xc-critical-resources | DOMAIN_ONLY |
| x-f5xc-danger-level | CONSUMED |
| x-f5xc-defaults | PARSED_NOT_RENDERED |
| x-f5xc-deprecated | PARSED_NOT_RENDERED |
| x-f5xc-description | CONSUMED |
| x-f5xc-description-long | DEFINED_UNUSED |
| x-f5xc-description-medium | CONSUMED |
| x-f5xc-description-short | CONSUMED |
| x-f5xc-discovered-at | NOT_EMITTED |
| x-f5xc-discovered-error-catalog | NOT_EMITTED |
| x-f5xc-discovered-rate-limits | NOT_EMITTED |
| x-f5xc-discovered-response-time | UNKNOWN |
| x-f5xc-display-name | PARSED_NOT_RENDERED |
| x-f5xc-displayorder | PARSED_NOT_RENDERED |
| x-f5xc-doc-section | DOMAIN_ONLY |
| x-f5xc-enriched-version | NOT_EMITTED |
| x-f5xc-error-resolution | DEFINED_UNUSED |
| x-f5xc-example | CONSUMED |
| x-f5xc-examples | PARSED_NOT_RENDERED |
| x-f5xc-glossary | DOMAIN_ONLY |
| x-f5xc-guided-workflows | DOMAIN_ONLY |
| x-f5xc-icon | DOMAIN_ONLY |
| x-f5xc-is-preview | CONSUMED |
| x-f5xc-logo-svg | DOMAIN_ONLY |
| x-f5xc-minimum-configuration | PARSED_NOT_RENDERED |
| x-f5xc-namespace-scope | DEFINED_UNUSED |
| x-f5xc-operation-metadata | PARSED_NOT_RENDERED |
| x-f5xc-primary-resources | DOMAIN_ONLY |
| x-f5xc-recommended-oneof-variant | PARSED_NOT_RENDERED |
| x-f5xc-recommended-value | PARSED_NOT_RENDERED |
| x-f5xc-related-domains | DOMAIN_ONLY |
| x-f5xc-required-fields | PARSED_NOT_RENDERED |
| x-f5xc-required-for | CONSUMED |
| x-f5xc-required-for-operations | PARSED_NOT_RENDERED |
| x-f5xc-requires-tier | CONSUMED |
| x-f5xc-response-time-ms | NOT_EMITTED |
| x-f5xc-server-default | CONSUMED |
| x-f5xc-side-effects | CONSUMED |
| x-f5xc-summary | DEFINED_UNUSED |
| x-f5xc-terraform-resource | PARSED_NOT_RENDERED |
| x-f5xc-uniqueness | PARSED_NOT_RENDERED |
| x-f5xc-upstream-etag | NOT_EMITTED |
| x-f5xc-upstream-timestamp | NOT_EMITTED |
| x-f5xc-use-cases | DOMAIN_ONLY |
| x-f5xc-validation | PARSED_NOT_RENDERED |

## 3. Resource-Level Drill-Downs

### AI

| Resource | Total Fields | Constraints % | Descriptions % | Enums % |
|----------|-------------|--------------|---------------|---------|
| ai_gateway | 0 | 0% | 0% | 0% |
| ai_policy | 0 | 0% | 0% | 0% |

### Infrastructure

| Resource | Total Fields | Constraints % | Descriptions % | Enums % |
|----------|-------------|--------------|---------------|---------|
| aws_vpc_site | 0 | 0% | 0% | 0% |
| azure_vnet_site | 0 | 0% | 0% | 0% |
| cloud_credentials | 71 | 27% | 38% | 0% |
| container_registry | 49 | 35% | 41% | 0% |
| endpoint | 68 | 28% | 31% | 0% |
| fleet_config | 205 | 49% | 56% | 0% |
| gcp_vpc_site | 0 | 0% | 0% | 0% |
| k8s_cluster_role | 103 | 32% | 45% | 0% |
| mk8s_cluster | 0 | 0% | 0% | 0% |
| pod_security_policy | 0 | 0% | 0% | 0% |
| registration_token | 92 | 43% | 46% | 0% |
| service_discovery | 0 | 0% | 0% | 0% |
| site | 474 | 53% | 36% | 0% |
| site_config | 0 | 0% | 0% | 0% |
| site_mesh_group | 0 | 0% | 0% | 0% |
| virtual_k8s | 64 | 30% | 41% | 0% |
| virtual_site | 43 | 30% | 33% | 0% |
| workload | 20 | 20% | 10% | 0% |

### Networking

| Resource | Total Fields | Constraints % | Descriptions % | Enums % |
|----------|-------------|--------------|---------------|---------|
| app_firewall | 137 | 21% | 24% | 0% |
| cdn_loadbalancer | 223 | 20% | 23% | 0% |
| cdn_origin_pool | 0 | 0% | 0% | 0% |
| dns_domain | 59 | 31% | 32% | 0% |
| dns_load_balancer | 85 | 31% | 47% | 0% |
| dns_zone | 263 | 48% | 43% | 0% |
| healthcheck | 70 | 41% | 51% | 0% |
| http_loadbalancer | 258 | 11% | 12% | 0% |
| malicious_user_detection | 0 | 0% | 0% | 0% |
| network_connector | 56 | 16% | 27% | 0% |
| origin_pool | 53 | 30% | 30% | 0% |
| rate_limit_threshold | 48 | 25% | 42% | 0% |
| rate_limiter | 57 | 25% | 40% | 0% |
| rate_limiter_policy | 49 | 16% | 39% | 0% |
| service_policy | 169 | 28% | 37% | 0% |
| tcp_loadbalancer | 80 | 12% | 24% | 0% |
| virtual_network | 82 | 33% | 44% | 0% |

### Operations

| Resource | Total Fields | Constraints % | Descriptions % | Enums % |
|----------|-------------|--------------|---------------|---------|
| alert | 0 | 0% | 0% | 0% |
| alert_policy | 0 | 0% | 0% | 0% |
| analytics_query | 0 | 0% | 0% | 0% |
| audit_log | 0 | 0% | 0% | 0% |
| dashboard | 0 | 0% | 0% | 0% |
| data_export | 0 | 0% | 0% | 0% |
| insight_query | 0 | 0% | 0% | 0% |
| log_receiver | 0 | 0% | 0% | 0% |
| metrics_receiver | 0 | 0% | 0% | 0% |
| saved_query | 0 | 0% | 0% | 0% |
| support_case | 68 | 46% | 54% | 0% |
| telemetry_receiver | 0 | 0% | 0% | 0% |

### Platform

| Resource | Total Fields | Constraints % | Descriptions % | Enums % |
|----------|-------------|--------------|---------------|---------|
| api_credential | 98 | 58% | 44% | 0% |
| authentication_policy | 0 | 0% | 0% | 0% |
| bigip_device | 0 | 0% | 0% | 0% |
| bigip_pool | 0 | 0% | 0% | 0% |
| bucket | 0 | 0% | 0% | 0% |
| marketplace_item | 9 | 56% | 56% | 0% |
| namespace_role | 0 | 0% | 0% | 0% |
| nginx_config | 0 | 0% | 0% | 0% |
| nginx_upstream | 0 | 0% | 0% | 0% |
| node_config | 0 | 0% | 0% | 0% |
| object_store | 53 | 51% | 40% | 0% |
| otp_policy | 0 | 0% | 0% | 0% |
| quota | 54 | 17% | 46% | 0% |
| session | 0 | 0% | 0% | 0% |
| static_asset | 0 | 0% | 0% | 0% |
| subscription | 0 | 0% | 0% | 0% |
| token | 0 | 0% | 0% | 0% |
| ui_component | 0 | 0% | 0% | 0% |
| usage_report | 53 | 70% | 38% | 0% |
| user | 0 | 0% | 0% | 0% |
| user_profile | 0 | 0% | 0% | 0% |
| user_role | 0 | 0% | 0% | 0% |
| vpm_config | 0 | 0% | 0% | 0% |

### Security

| Resource | Total Fields | Constraints % | Descriptions % | Enums % |
|----------|-------------|--------------|---------------|---------|
| api_definition | 19 | 63% | 53% | 0% |
| api_endpoint | 27 | 30% | 56% | 0% |
| api_rate_limit | 0 | 0% | 0% | 0% |
| blindfold_secret | 0 | 0% | 0% | 0% |
| bot_defense_instance | 20 | 30% | 55% | 0% |
| ca_certificate | 41 | 27% | 44% | 0% |
| certificate | 98 | 21% | 37% | 0% |
| certificate_chain | 40 | 25% | 42% | 0% |
| data_classification | 61 | 28% | 46% | 0% |
| ddos_mitigation_rule | 201 | 20% | 28% | 0% |
| ddos_protection | 783 | 35% | 37% | 0% |
| forward_proxy_policy | 24 | 0% | 0% | 0% |
| malicious_user_rule | 49 | 18% | 35% | 0% |
| mitigation_policy | 49 | 18% | 35% | 0% |
| network_firewall | 64 | 19% | 30% | 0% |
| network_policy | 226 | 29% | 38% | 0% |
| policy_document | 59 | 32% | 36% | 0% |
| secret_policy | 97 | 31% | 40% | 0% |
| sensitive_data_policy | 35 | 23% | 40% | 0% |
| shape_app_firewall | 0 | 0% | 0% | 0% |
| shape_recognizer | 0 | 0% | 0% | 0% |
| threat_campaign_policy | 198 | 57% | 66% | 0% |
| threat_category | 43 | 28% | 44% | 0% |

## 4. Validator Opportunity Analysis

- **Constraint validators available**: 1800 fields across 95 resources
- **Enum validators available**: 0 fields with enum values that could generate stringvalidator.OneOf
- **ConflictsWith validators available**: 77 fields with conflict declarations

## 5. Schema Fidelity Findings

- 41 of 57 extensions are not fully consumed
- 3338 fields lack enriched descriptions (of 5347 total)
- 3707 of 3707 operations have danger level annotations (100%)

## 6. Spec Enrichment Priorities

| Config File | Purpose | Priority |
|------------|---------|----------|
| extension_constants.py | Extension registry | High |
| index.json | Resource-to-domain map | High |
| Domain spec files | Per-field extensions | Medium |
| Operation metadata | Operation-level annotations | Medium |

## 7. Downstream Impact Assessment

| Consumer | Impact Area | Dependency |
|----------|------------|------------|
| terraform-provider-f5xc | Schema generation | api-specs-enriched |
| terraform-provider-f5xc | Validator generation | x-f5xc-constraints, enum |
| terraform-provider-f5xc | Documentation | x-f5xc-description-*, x-f5xc-best-practices |
| api-specs-enriched | Extension emission | extension_constants.py |

## 8. Prioritized Action Items

| Rank | Action | Repo | Priority Score |
|------|--------|------|---------------|
| 1 | Consume x-f5xc-constraints to generate Terraform validators | terraform-provider-f5xc | 3.0 |
| 2 | Implement ConflictsWith validators | terraform-provider-f5xc | 2.5 |
| 3 | Emit and replace hardcoded maps | both | 2.5 |
| 4 | Consume for Required field accuracy | terraform-provider-f5xc | 2.0 |
| 5 | Generate stringvalidator.OneOf from enum fields | terraform-provider-f5xc | 2.0 |
| 6 | Render in generated code | terraform-provider-f5xc | 2.0 |
| 7 | Fix struct shape mismatch and embed in docs | terraform-provider-f5xc | 2.0 |
| 8 | Audit 49 orphan data source files | terraform-provider-f5xc | 2.0 |
| 9 | Consume operation-level extensions for destruction warnings | terraform-provider-f5xc | 1.3 |
| 10 | Wire up index-derived dependency map or remove | terraform-provider-f5xc | 1.0 |

## GitHub Issue Templates

### Consume x-f5xc-constraints to generate Terraform validators

**Repo**: terraform-provider-f5xc
**Priority Score**: 3.0
**User Impact**: 3 | **Downstream Reach**: 3 | **Effort**: 2

The x-f5xc-constraints extension is parsed into Go structs but never rendered into Terraform validation logic. Consuming this extension would auto-generate validators for min/max, pattern, and custom constraints.

### Implement ConflictsWith validators

**Repo**: terraform-provider-f5xc
**Priority Score**: 2.5
**User Impact**: 3 | **Downstream Reach**: 2 | **Effort**: 2

The x-f5xc-conflicts-with extension is parsed but not rendered. Generating ConflictsWith validators would prevent users from setting mutually exclusive fields.

### Emit and replace hardcoded maps

**Repo**: both
**Priority Score**: 2.5
**User Impact**: 2 | **Downstream Reach**: 3 | **Effort**: 2

The x-f5xc-namespace-scope extension is registered but not emitted in any spec file. Once emitted, it can replace the hardcoded namespace scope maps in the Terraform provider.

### Consume for Required field accuracy

**Repo**: terraform-provider-f5xc
**Priority Score**: 2.0
**User Impact**: 2 | **Downstream Reach**: 2 | **Effort**: 2

The x-f5xc-minimum-configuration extension defines the minimal set of fields needed for a valid resource. Consuming it would improve Required field accuracy in the Terraform schema.

### Generate stringvalidator.OneOf from enum fields

**Repo**: terraform-provider-f5xc
**Priority Score**: 2.0
**User Impact**: 3 | **Downstream Reach**: 1 | **Effort**: 2

Enum fields in the spec define valid string values. Generating stringvalidator.OneOf validators from these would catch invalid values at plan time instead of apply time.

### Render in generated code

**Repo**: terraform-provider-f5xc
**Priority Score**: 2.0
**User Impact**: 2 | **Downstream Reach**: 2 | **Effort**: 2

The x-f5xc-required-for extension is parsed but not rendered. Rendering it would document which fields are required for specific operations (create, update, etc.).

### Fix struct shape mismatch and embed in docs

**Repo**: terraform-provider-f5xc
**Priority Score**: 2.0
**User Impact**: 1 | **Downstream Reach**: 3 | **Effort**: 2

The x-f5xc-best-practices extension has a struct shape mismatch between the spec definition and the Go struct. Fixing this and embedding best practices in generated documentation would guide users toward optimal configurations.

### Audit 49 orphan data source files

**Repo**: terraform-provider-f5xc
**Priority Score**: 2.0
**User Impact**: 1 | **Downstream Reach**: 1 | **Effort**: 1

There are approximately 49 data source files that are not connected to any primary resource in index.json. These should be audited and either connected or removed.

### Consume operation-level extensions for destruction warnings

**Repo**: terraform-provider-f5xc
**Priority Score**: 1.3
**User Impact**: 2 | **Downstream Reach**: 2 | **Effort**: 3

The x-f5xc-danger-level extension marks operations that may cause service disruption. Consuming it would add destruction warnings and confirmation prompts to dangerous operations.

### Wire up index-derived dependency map or remove

**Repo**: terraform-provider-f5xc
**Priority Score**: 1.0
**User Impact**: 1 | **Downstream Reach**: 1 | **Effort**: 2

The dependency map derived from index.json is computed but never wired into the provider. It should either be connected to resource ordering logic or removed as dead code.
