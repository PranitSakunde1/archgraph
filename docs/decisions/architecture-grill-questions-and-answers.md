# ArchGraph Architecture Grill: Questions and Answers

> **Purpose:** This companion record preserves the architecture-grill decisions as question-and-answer pairs for future implementation, review and onboarding.
>
> **Interpretation:** Questions below are normalized summaries of the decisions, not verbatim chat transcripts. Answers reflect the final recorded choice. Where a decision was later corrected, the latest/superseding decision is authoritative. The numbered decision register remains the canonical detailed decision log.

## 1. Core discovery and governance

### Q1. Should ArchGraph use AI/LLMs during repository discovery?
**A:** No. Discovery, parsing, graph generation, enrichment and validation must be deterministic and AI-free.

### Q2. Where should the canonical architecture graph live?
**A:** Locally in the ArchGraph knowledge base. A shared Neo4j layer may be considered later as an optional hybrid extension.

### Q3. What should trigger architecture scans?
**A:** Manual scans, incremental pre-commit scans, full pre-push scans and CI validation.

### Q4. What Git hooks should be supported?
**A:** Opt-in incremental pre-commit and full pre-push consistency scans, with policy-based blocking.

### Q5. Which generated artifacts should be committed?
**A:** Commit full, lightweight graph artifacts. Keep detailed caches and transient data local.

### Q6. How should unsupported code or services be handled?
**A:** Report them visibly and retain them in the service inventory. Never silently omit them.

### Q7. Should unsupported coverage block a scan or CI?
**A:** Use baseline-aware enforcement. Existing unsupported services are reported; newly introduced unsupported-coverage regressions may block.

### Q8. What is the source of truth for service and Lambda identity?
**A:** `codeway.yaml` is authoritative when present.

### Q9. Where should service metadata live?
**A:** Each service or Lambda maintains its own `codeway.yaml`.

### Q10. What evidence should establish architecture relationships?
**A:** Deterministic evidence combined from source code, contracts, infrastructure and configuration.

### Q11. How should ambiguous relationships be represented?
**A:** Keep confirmed and inferred relationships separate. During manual scans, prompt for important ambiguous relationships.

### Q12. Who authors critical user/business journeys?
**A:** Engineers author critical journeys. ArchGraph may suggest candidates, but execution requires human approval.

### Q13. How should journey environments select dependencies?
**A:** Select the minimal service and infrastructure set required by the journey, based on dependency relationships.

### Q14. How should asynchronous journeys be tested?
**A:** Use event-driven synchronization with correlation IDs, acknowledgements, bounded condition waits and business-outcome assertions.

### Q15. How should journey test data be isolated?
**A:** Use baseline fixtures plus unique dynamic values, a per-run namespace and deterministic cleanup.

### Q16. What should change-impact analysis provide?
**A:** Graph diffs, blast-radius analysis and recommendations for affected journeys.

### Q17. What snapshot model should be used?
**A:** Commit-level snapshots plus immutable snapshots for approved releases.

### Q18. How should graph quality be evaluated?
**A:** Use structural validation, coverage validation and quality scoring. Critical gaps may independently fail validation.

### Q19. Should quality thresholds be treated as proven?
**A:** No. Start with fixed provisional thresholds and calibrate them against real repositories.

### Q20. How should quality gates treat existing gaps?
**A:** Be baseline-aware. Report existing gaps and block newly introduced regressions.

## 2. Review, AWS, data and graph model

### Q21. Who approves changes to the architecture baseline?
**A:** CODEOWNERS review is required for baseline changes.

### Q22. What is the pull-request and CI workflow?
**A:** CI generates graph diffs/reports and validates committed artifacts. Developers commit updates. CI must not push or rewrite branches.

### Q23. How should graph history and size be managed?
**A:** Keep full reproducible snapshots and incremental diffs.

### Q24. Which architecture diagrams should be generated?
**A:** Layered system, domain/group, service, event-flow and journey diagrams.

### Q25. How should services be grouped into domains?
**A:** Use explicit ownership metadata, supplemented by human-approved suggestions.

### Q26. How should ownership be represented?
**A:** Capture both domain ownership and per-service team ownership, with optional technical owner/escalation metadata.

### Q27. How should API discovery work?
**A:** Combine OpenAPI/contracts, source handler discovery and client discovery.

### Q28. Should API compatibility gaps fail validation?
**A:** Apply severity and policy. Breaking or critical consumer incompatibility may block; non-breaking gaps should warn.

### Q29. How should event discovery work?
**A:** Use code-first discovery.

### Q30. What role should infrastructure definitions play in event discovery?
**A:** Terraform, CloudFormation, SAM and AsyncAPI enrich or validate code evidence. Missing infrastructure definitions warn rather than automatically fail.

### Q31. How broad should AWS discovery be?
**A:** Support broad AWS-native resources and event services, not only a fixed shortlist.

### Q32. Which AWS SDK languages should be supported?
**A:** All languages and SDKs present in the monorepo are in scope. Go and Java are first-class, with other SDK adapters.

### Q33. Should ArchGraph query live AWS accounts in v1?
**A:** No. Discover AWS resources from source code, resource files, dependencies and configuration. Do not query live AWS accounts in v1.

### Q34. How should custom AWS wrappers be discovered?
**A:** Provide built-in wrapper adapters and a generic source-inspection fallback.

### Q35. Where should wrapper mappings be maintained?
**A:** Built-in mappings belong in ArchGraph; repository-specific mappings belong under `.archgraph/plugins/`.

### Q36. How should wrapper adapters be implemented?
**A:** Use declarative YAML mappings for simple wrappers and trusted Go interfaces for complex cases.

### Q37. How should dependency cycles be handled?
**A:** Detect and classify synchronous cycles and asynchronous feedback loops. Intentional documented loops may be allowed.

### Q38. Is deployment topology part of the initial relationship model?
**A:** No. The selected scope is application relationships, not deployment topology.

### Q39. What database relationships should be discovered?
**A:** Table- and collection-level access and relationships.

### Q40. How deeply should database schemas be discovered?
**A:** Discover full schemas where statically available, including fields, keys and indexes.

## 3. Data lineage, security and language support

### Q41. Should shared database access automatically be flagged as a defect?
**A:** No. Observe reads and writes, but do not warn merely because multiple services access the same table or collection. Do not impose a mandatory single-writer policy.

### Q42. How comprehensive should data-flow discovery be?
**A:** Discover deterministic data flow across APIs, events, services, databases and field mappings.

### Q43. How should data lineage be represented?
**A:** Use both field-level lineage edges and explicit transformation nodes.

### Q44. How deep should transformation discovery go?
**A:** Capture structural and operational transformations, nested objects/collections, serialization/deserialization and cross-service propagation where statically identifiable.

### Q45. Should actual sensitive data values be stored?
**A:** No. Classify sensitive fields and track their lineage without storing actual values.

### Q46. Where should data-classification rules come from?
**A:** Provide built-in defaults plus repository-specific extensions and overrides.

### Q47. Which classification rule wins when rules conflict?
**A:** Repository rules win for ordinary categories. Protected credential/secret categories cannot be downgraded, and the original built-in classification remains traceable.

### Q48. How should secrets be detected?
**A:** Use deterministic pattern matching combined with context analysis.

### Q49. How should detected secret values be handled in outputs?
**A:** Redact secret values from terminal output, logs, graph artifacts and CI reports.

### Q50. What should happen when a secret is detected?
**A:** Confirmed secrets block. Potential secrets are reported for review. Approved false-positive exceptions must be auditable.

### Q51. Is ArchGraph intended to replace dedicated security scanners?
**A:** No. Integrate dedicated scanners and correlate findings with architecture. ArchGraph does not replace them.

### Q52. How should security gates handle existing findings?
**A:** Use baseline-aware gates. Newly introduced critical/high findings block; existing backlog is reported and tracked.

### Q53. How extensible should security scanner integration be?
**A:** Use built-in integrations, custom Go adapters and YAML configuration.

### Q54. How should excluded files and resources be handled?
**A:** Support repository-level configuration and ignore rules with explicit reasons and auditability.

### Q55. Should ArchGraph modify scanned source repositories?
**A:** No. Source access is read-only by default.

### Q56. How should graph entity IDs be generated?
**A:** Use stable IDs based on logical identity, not file paths or line numbers alone.

### Q57. How should graph schema changes be managed?
**A:** Use explicit schema versioning, validation and migration support.

### Q58. How should analyzer compatibility be enforced?
**A:** Version native language/framework adapters and validate their output against the canonical graph schema. Fail clearly on incompatibility.

### Q59. What happens when a parser cannot understand a construct?
**A:** Report the unsupported construct with evidence and partial coverage. Never silently discard it.

### Q60. Which languages and frameworks are in scope?
**A:** Go, Java, Python, React, TypeScript and JavaScript. React is a framework-specific adapter for components, hooks, routes and API clients, not a programming language.

## 4. Runtime, developer experience and rollout

### Q61. How should unresolved dynamic behavior be represented?
**A:** Mark unresolved static-analysis paths explicitly. Do not invent call or data-flow relationships.

### Q62. Which journeys may be executed?
**A:** Execute only approved journeys, prioritizing journeys affected by graph changes.

### Q63. How should journey test environments be isolated?
**A:** Use an isolated Docker environment and per-run data namespace with deterministic cleanup.

### Q64. How should journeys run in parallel?
**A:** Parallelize independent journeys. Serialize journeys that share mutable resources.

### Q65. Is runtime intelligence mandatory for v1?
**A:** No. Runtime intelligence is an optional future enhancement; static analysis is the v1 foundation.

### Q66. Which report formats are required?
**A:** Machine-readable JSON, human-readable Markdown and Mermaid diagrams.

### Q67. Which core CLI capabilities are required?
**A:** Commands for scan, diff, validate, query, impact, journey and report.

### Q68. What agent integration should ArchGraph expose?
**A:** ArchGraph-specific MCP tools for normalized architecture queries, impact analysis, journey planning and reporting as those capabilities become ready.

### Q69. What performance strategy should be used?
**A:** Incremental scans, caching and dependency-aware invalidation, alongside full reproducible CI scans.

### Q70. How should delivery be phased?
**A:** Deliver in phases, progressing from graph foundation to relationship intelligence, impact analysis, journey testing, security and optional hybrid scale.

## 5. Hybrid analyzer architecture and product boundaries

### Q71. Should ArchGraph become a separate Java/Go graph product or remain one product?
**A:** Continue in the existing ArchGraph repository. ArchGraph remains the single user-facing product and owns orchestration, canonical graph semantics, evidence validation and exports.

### Q72. Which parser/analyzer strategy should be used?
**A:** Use hybrid semantic analysis. Go gets native semantic analysis through the Go toolchain ecosystem, including packages, AST, types, SSA and callgraph. Java gets compiler-aware analysis plus Spring-specific discovery. Tree-sitter is a shared syntax extraction and fallback layer, not a substitute for compiler/type resolution.

### Q73. Should Graphify or another external analyzer own the architecture model?
**A:** No. External analyzers may be evaluated as optional adapters or implementation references. They must not own ArchGraph's canonical schema, evidence policy, identity rules or validation. No mandatory Graphify dependency is selected.

### Q74. Is runtime analysis required for static graph generation?
**A:** No. Static analysis and export must work without runtime instrumentation or trace data. Optional OpenTelemetry observations may enrich the graph but must remain distinguishable from static evidence and include provenance.

### Q75. What is the primary deployment/interface model?
**A:** One native Go CLI is the primary product interface. Docker is optional packaging/execution. The core workflow does not require an always-on daemon or local server.

### Q76. Which languages receive first-class implementation priority?
**A:** Java Spring Boot and Go are first-class. Python, TypeScript, JavaScript and React are secondary/future adapters. React is a framework adapter, not a language.

### Q77. Who owns the canonical graph schema?
**A:** ArchGraph owns a versioned, language-neutral schema. Analyzer adapters emit normalized facts; the core validates and reconciles them. JSON is canonical; Markdown and Mermaid are derived views.

### Q78. What evidence standard applies to published relationships?
**A:** Every source-derived relationship must carry provenance and evidence. Unresolved, inferred, declared, statically confirmed, runtime-observed and conflicting evidence remain distinguishable. Never promote unsupported relationships to confirmed or claim 100% completeness.

### Q79. Where should generated architecture documentation live?
**A:** In the target service repository as Markdown, Mermaid and structured JSON. Generated artifacts must be reproducible, versioned and separated from manually maintained documentation.

### Q80. What should Git refresh and CI validation do?
**A:** Regenerate affected artifacts during the Git commit workflow and verify them in CI. Hooks stage only ArchGraph-owned files. CI fails on stale or invalid committed artifacts and never pushes or rewrites branches.

## 6. Local Docker runtime observation

### Q81. Where should the runtime observer run?
**A:** As a native local process on the Docker host, including WSL/Linux, not as a container with Docker socket access.

### Q82. What can Docker metadata prove?
**A:** It can establish container lifecycle, image, labels, networks and port configuration. It does not prove application-level communication.

### Q83. How should actual application interactions be observed?
**A:** Use optional OpenTelemetry instrumentation and a local OTLP receiver/export path for service calls, request paths and supported database/messaging client operations. Missing instrumentation is an observability gap, not proof that no relationship exists.

### Q84. Where may runtime data be sent?
**A:** Nowhere outside the local machine. Collection, correlation, persistence and graph updates remain local. No SaaS, cloud service, LLM or external telemetry destination is required.

### Q85. How should static, Docker and runtime evidence be represented?
**A:** Keep Docker topology, statically confirmed relationships and runtime-observed relationships as distinct evidence classes. Runtime observations include capture window/time, environment, source identity, instrumentation/source and trace-correlation metadata where available.

### Q86. What runtime data must not be persisted?
**A:** Do not persist request/response bodies, credentials, authorization headers, raw SQL parameter values or sensitive payloads by default. Redact sensitive attributes and support configurable retention and local deletion. Do not give an observer container Docker socket access.

### Q87. How should the runtime observer operate?
**A:** Static scanning remains a one-shot CLI workflow. Runtime observation has a separate start/stop lifecycle and may run while application containers are active. It can enrich/update the local graph continuously or through an explicit snapshot, but is not required for static graph generation.

## 7. Corrections and interpretation guardrails

- React is a framework adapter, not a programming language.
- Dependency presence alone is not proof that a service uses an AWS resource.
- Shared database access is an observed fact, not automatically a defect.
- Initial quality thresholds are provisional until measured against real repositories.
- Inspect the actual `codeway.yaml` format and internal AWS wrapper patterns before finalizing parser behavior.
- The selected Docker scan mode is one-shot: mount the workspace read-only, scan, export and exit. The native CLI remains primary.
- Decisions 71-80 supersede earlier interpretations that ArchGraph is Docker-only. Native CLI first, Docker optional.
- Runtime observation is optional enrichment, not a prerequisite for graph generation.
- Hybrid semantic analysis supersedes Tree-sitter-only or regex-only extraction as the selected strategy.
- Evidence traceability is mandatory; complete discovery or 100% accuracy is not promised.
- The existing Go core remains the orchestration and graph authority. Java analysis may use a bundled, versioned local adapter process if Java tooling requires it.
- The runtime observer is a native host process. Docker metadata is topology evidence only, never proof of HTTP, RPC, database or messaging interactions.
- Runtime data is local-only by default. Any future remote forwarding/export requires a separate explicit decision and opt-in.

## Related decision records

- [Architecture Grill: Final Decision Register](./architecture-grill-decisions.md)
- [ADR-001: Hybrid Evidence-First Analysis](./ADR-001-hybrid-evidence-first-analysis.md)
- [Architecture Blueprint](../architecture/architecture-blueprint.md)
- [Docker Runtime](../architecture/docker-runtime.md)
- [Production Build Prompt](../implementation/ARCHGRAPH-MASTER-BUILD-PROMPT.md)
