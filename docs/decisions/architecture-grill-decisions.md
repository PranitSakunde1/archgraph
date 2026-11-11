# ArchGraph Architecture Grill: Final Decision Register

This document captures the final choices from the architecture grill. It is a structured decision record, not a verbatim chat transcript. Where a choice was corrected later, the latest choice is authoritative.

## Decisions 1-20: Core discovery and governance

1. **AI-free discovery:** no LLM calls during scanning, parsing, graph generation, enrichment or validation.
2. **Graph location:** local-first canonical ArchGraph graph; optional hybrid shared Neo4j later.
3. **Scan triggers:** manual scans, incremental pre-commit, full pre-push and CI.
4. **Git hooks:** incremental pre-commit and full pre-push consistency scan; policy-based blocking.
5. **Artifacts:** commit full lightweight graph artifacts; keep detailed caches and transient data local.
6. **Unsupported code:** fail/report unsupported services visibly, retain them in inventory, never silently omit.
7. **Unsupported enforcement:** baseline-aware; existing unsupported services report, newly introduced unsupported coverage regressions can block.
8. **Service identity:** codeway.yaml is authoritative for service/Lambda identity.
9. **Metadata layout:** each service/Lambda maintains its own codeway.yaml.
10. **Relationship evidence:** deterministic combination of code, contracts, infrastructure and configuration.
11. **Ambiguity:** separate confirmed and inferred relationships; prompt for important ambiguous relationships during manual scans.
12. **Journey authorship:** engineers author critical journeys; ArchGraph may suggest candidates, but execution requires human approval.
13. **Journey environment:** dependency-aware minimal service and infrastructure selection.
14. **Async testing:** event-driven synchronization with correlation IDs, acknowledgements, bounded condition waits and business-outcome assertions.
15. **Test data:** baseline fixtures plus unique dynamic values, per-run namespace and deterministic cleanup.
16. **Change impact:** graph diffs, blast-radius analysis and journey recommendations.
17. **Snapshot strategy:** commit-level snapshots plus immutable approved release snapshots.
18. **Graph quality:** structural validation, coverage validation and quality scoring; critical gaps can independently fail.
19. **Quality thresholds:** fixed initial thresholds, provisional until measured against real repositories.
20. **Quality gates:** baseline-aware; existing gaps report, new regressions block.

## Decisions 21-40: Review, AWS, data and graph model

21. **Baseline approval:** CODEOWNERS review for baseline changes.
22. **PR workflow:** CI generates graph diff/report and validates committed artifacts; developers commit updates; CI does not push or rewrite branches.
23. **Graph size:** full reproducible snapshots plus incremental diffs.
24. **Diagrams:** layered system, domain/group, service, event-flow and journey diagrams.
25. **Domain grouping:** explicit ownership metadata plus human-approved suggestions.
26. **Ownership:** both domain ownership and per-service team ownership; optional technical owner/escalation.
27. **API discovery:** combine OpenAPI/contracts, source handlers and client discovery.
28. **API enforcement:** severity/policy based; breaking or critical consumer incompatibility may block, non-breaking gaps warn.
29. **Event discovery:** code-first.
30. **Event infrastructure evidence:** Terraform, CloudFormation, SAM and AsyncAPI enrich/validate code evidence; missing definitions warn rather than automatically fail.
31. **AWS scope:** broad AWS-native resources and event services, not only SNS/SQS.
32. **SDK coverage:** all languages/SDKs in the monorepo are in scope; Go and Java are first-class, with other SDK adapters.
33. **AWS resource discovery:** derive from source code, resource files, dependencies and configuration; no live AWS account queries in v1.
34. **Wrapper discovery:** built-in wrapper adapters plus generic source inspection fallback.
35. **Wrapper mapping location:** built-ins in ArchGraph and repository-specific mappings under .archgraph/plugins/.
36. **Adapter style:** YAML declarative mappings for simple wrappers; trusted Go interfaces for complex cases.
37. **Cycles:** detect/classify synchronous cycles and asynchronous feedback loops; intentional documented loops can be allowed.
38. **Deployment topology:** out of scope; application relationships only.
39. **Database discovery:** table/collection-level relationships and access.
40. **Database schema:** full schema discovery where statically available, including fields, keys and indexes.

## Decisions 41-60: Data lineage, security and language support

41. **Shared database access:** observe reads/writes; do not warn merely because multiple services access the same table or collection; no mandatory single-writer policy.
42. **Data flow:** comprehensive deterministic data-flow discovery across APIs, events, services, databases and field mappings.
43. **Lineage representation:** both field-level lineage edges and explicit transformation nodes.
44. **Transformation depth:** structural and operational transformations, nested objects/collections, serialization/deserialization and cross-service propagation where statically identifiable.
45. **Sensitive data:** classify sensitive fields and track lineage without storing actual values.
46. **Classification rules:** built-in defaults plus repository-specific extensions/overrides.
47. **Classification precedence:** repository rules win for ordinary categories; protected credential/secret categories cannot be downgraded; retain original built-in classification for traceability.
48. **Secret detection:** deterministic pattern matching plus context analysis.
49. **Secret output handling:** redact secret values from terminal output, logs, graph artifacts and CI reports.
50. **Secret enforcement:** confirmed secrets block; potential secrets are reported for review; approved false-positive exceptions are auditable.
51. **Security scope:** integrate dedicated scanners and correlate findings with architecture; ArchGraph does not replace scanners.
52. **Security gates:** baseline-aware; newly introduced critical/high findings block, existing backlog is reported and tracked.
53. **Scanner extensibility:** hybrid built-in integrations, custom Go adapters and YAML configuration.
54. **Exclusions:** repository-level configuration and ignore rules with explicit reasons and auditability.
55. **Source access:** read-only scanning by default.
56. **Graph identity:** stable IDs based on logical identity, not file paths or line numbers alone.
57. **Schema evolution:** explicit graph schema versioning, validation and migration support.
58. **Extractor compatibility:** version native language/framework adapters and validate their output against the canonical graph schema; fail clearly on incompatibility.
59. **Parsing fallback:** unsupported constructs are reported with evidence and partial coverage, never silently discarded.
60. **Language targets (updated):** Go, Java, Python, React, TypeScript and JavaScript. React gets a framework-specific adapter for components, hooks, routes and API clients.

## Decisions 61-70: Runtime, developer experience and rollout

61. **Dynamic behavior:** mark unresolved static analysis paths explicitly; do not invent call or data-flow relationships.
62. **Journey execution:** execute only approved journeys and prioritize those affected by graph changes.
63. **Test isolation:** isolated Docker environment and per-run data namespace with deterministic cleanup.
64. **Parallel execution:** parallelize independent journeys; serialize journeys sharing mutable resources.
65. **Runtime intelligence:** optional future enhancement; static analysis remains the v1 foundation.
66. **Reports:** machine-readable JSON, human-readable Markdown and Mermaid diagrams.
67. **CLI:** commands for scan, diff, validate, query, impact, journey and report.
68. **MCP:** add ArchGraph-specific tools for normalized architecture queries, impact analysis, journey planning or reporting when those capabilities are ready.
69. **Performance:** incremental scans, caching and dependency-aware invalidation, with full reproducible CI scans.
70. **Rollout:** phased delivery from graph foundation to relationship intelligence, impact analysis, journey testing, security and optional hybrid scale.

## Corrections and explicit guardrails

- React is a framework adapter, not a programming language.
- Dependency presence is not proof that a service actually uses an AWS resource.
- Shared database access is a fact, not automatically a defect.
- Proposed quality thresholds are starting targets, not measured coverage claims.
- The actual codeway.yaml format and internal AWS wrapper patterns must be inspected before parser behavior is finalized.
- Docker one-shot mode is the selected operating model: mount workspace read-only, scan, export, exit. No daemon/API server in v1.
