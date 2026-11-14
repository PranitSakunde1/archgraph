# ArchGraph: Production Build Prompt

> **Purpose:** This is the approved master engineering specification for GitHub Copilot to implement ArchGraph. Treat it as the source of truth for the build. Do not begin implementation until the user approves the first implementation phase.

## 1. Role and mission

Act as the principal software architect and senior engineer responsible for evolving the existing ArchGraph repository into a production-quality, local-first architecture intelligence and knowledge-generation CLI.

ArchGraph analyzes software repositories and produces a persistent, evidence-backed architecture knowledge base in two forms:

1. A machine-readable canonical graph.
2. Human-readable architecture documentation.

The central product outcome is reusable engineering knowledge. Developers and AI agents should be able to answer new engineering questions from previously generated ArchGraph knowledge instead of repeatedly scanning every source repository.

ArchGraph itself must remain AI-free. It must not require an LLM, cloud service, SaaS product or external telemetry destination.

## 2. Non-negotiable principles

- Local-first execution and zero mandatory cloud spend.
- No LLM dependency and no SaaS dependency.
- No invented architecture relationships.
- Every source-derived relationship carries evidence and provenance.
- Static, declared, inferred, unresolved, conflicting and runtime evidence remain distinct.
- Individual repositories own their detailed architecture knowledge.
- A parent workspace aggregates repository knowledge and cross-repository relationships.
- Incremental regeneration is the normal workflow.
- AI agents retrieve existing knowledge through focused interfaces rather than rescanning the whole workspace.
- Human-readable documentation is a first-class product output, not a decorative export.
- Generated artifacts are deterministic wherever possible.
- Never overwrite user-maintained files or modify application source and infrastructure.
- No destructive Git operations, automatic commits or pushes.
- Work in approved phases with explicit acceptance gates.

## 3. Mandatory first step: inspect before editing

Before changing source code:

1. Inspect the complete repository tree and current branch.
2. Read README, AGENTS.md, PROJECT_CONTEXT.md if present, and all existing architecture blueprints, ADRs, decision logs and runtime documentation.
3. Inspect go.mod, CLI entry points, existing packages, tests, Makefile, Dockerfile and CI workflows.
4. Identify existing capabilities, technical debt, duplicated responsibilities, oversized files and mismatches between documentation and implementation.
5. Preserve useful existing behavior. Do not assume existing code or documentation is correct without inspection.

Create an assessment under docs/implementation/:

- repository-assessment.md
- existing-capabilities.md
- architecture-gap-analysis.md
- implementation-roadmap.md
- proposed-package-structure.md

Present the assessment, proposed architecture, phase boundaries, dependencies, acceptance criteria and risks. Do not modify application source code until the user approves Phase 1.

## 4. Product scope and supported technologies

First-class application languages:

- Go
- Java, including Java Spring Boot
- Python

Discover, where supported:

- Repositories, modules, packages, namespaces, classes, interfaces, functions and methods.
- Imports, dependencies, internal calls and cross-package relationships.
- HTTP/REST endpoints, gRPC endpoints and clients.
- Messaging producers, consumers, queues and topics.
- Database clients, data stores and discoverable read/write relationships.
- Framework integrations, generated clients and dependency injection.
- Configuration references, environment variables and build dependencies.
- External service dependencies and deployment relationships.

Infrastructure and specification formats:

- Dockerfiles and Docker Compose.
- Kubernetes manifests and Helm configuration.
- Terraform.
- AWS configuration and infrastructure declarations.
- CI/CD workflows.
- YAML, JSON, properties and relevant environment configuration.
- OpenAPI, Swagger and Protocol Buffers.

Use compiler-aware or semantic tooling where required to establish reliable relationships. Do not build deep architecture claims on regex-only scanning. Where resolution is not possible, classify the result as inferred, unresolved or conflicting and explain the limitation.

Report analyzer coverage by repository, language, framework, infrastructure type and evidence classification. Do not claim support without fixtures and tests.

## 5. AWS discovery

AWS discovery must be extensible and must not be restricted to a hardcoded shortlist.

Discover AWS integrations from SDK client construction, method invocation, dependency injection, configuration references, resource identifiers, environment variables, infrastructure declarations and IAM policies where supported.

Examples include SNS, SQS, S3, Lambda, DynamoDB, RDS, EventBridge, Secrets Manager, CloudWatch and Bedrock, plus other services supported by implemented analyzers.

Distinguish these findings:

- SDK dependency detected.
- Client constructed.
- Service operation invoked.
- Resource configured.
- Infrastructure resource declared.
- Runtime interaction observed.

An SDK import alone is not proof of an actual AWS operation or resource interaction.

## 6. Canonical graph model

ArchGraph owns a language-neutral, versioned canonical graph schema.

Use a hybrid model: typed core entities with extensible metadata, explicit relationship types, stable deterministic IDs, schema versioning and explicit backward-compatible migrations.

Entity categories should include:

- Workspace, repository and service.
- Module, package, class, interface, function and method.
- API and contract.
- Database, table or collection.
- Queue, topic and messaging resource.
- AWS resource and generic infrastructure resource.
- Configuration reference.
- Container and deployment.

Relationship types must preserve architectural meaning. Examples include CALLS, IMPORTS, DEPENDS_ON, EXPOSES, CONSUMES, PRODUCES, PUBLISHES_TO, SUBSCRIBES_TO, READS_FROM, WRITES_TO, DEPLOYED_ON, CONFIGURED_BY, IMPLEMENTS, CONTAINS and OBSERVED_TO_CALL.

Avoid representing every relationship with a generic edge.

Evidence classifications:

- STATIC_CONFIRMED
- DECLARED
- RUNTIME_OBSERVED
- INFERRED
- UNRESOLVED
- CONFLICTING

Never silently promote one evidence classification into another.

Every relationship must carry provenance sufficient to explain its origin, including repository identity, source location where available, analyzer/discovery method, evidence classification and relevant observation metadata. Confidence scores may supplement classifications but must not be presented as statistical certainty unless calibrated against validated data.

Stable identity must survive ordinary source edits and file moves wherever identity can be resolved reliably. Do not base identity solely on file paths or content hashes. Document identity rules and test collision, rename and reconciliation behavior.

## 7. Repository-level knowledge base

Every analyzed repository must own its own generated architecture knowledge.

Target output:

~~~text
repository/
  .archgraph/
    docs/
      README.md
      architecture-overview.md
      service-catalog.md
      code-structure.md
      dependency-map.md
      api-catalog.md
      data-access.md
      aws-integrations.md
      infrastructure.md
      configuration.md
      security-and-exposure.md
      execution-flows.md
      evidence-and-limitations.md
      diagrams/
        component-architecture.mmd
        dependency-map.mmd
        request-flows.mmd
        infrastructure.mmd
    graph/
      graph.json
      nodes.json
      relationships.json
      evidence.json
      manifest.json
      schema-version.json
    metadata/
      repository.json
      analysis-summary.json
      coverage.json
~~~

The canonical graph is the source for generated documentation. Do not maintain a second independently authored architecture model. Generated documents must explain the architecture in readable language, include navigation and cross-links, and expose evidence and limitations.

## 8. Workspace-level knowledge base and detection

When ArchGraph runs against a directory containing multiple repositories, detect the workspace and generate:

- An independent knowledge base inside each child repository.
- A consolidated workspace knowledge base.
- Cross-repository relationships in the workspace graph.

Target workspace output:

~~~text
workspace/
  .archgraph/
    docs/
      README.md
      workspace-overview.md
      repository-catalog.md
      system-architecture.md
      cross-service-dependencies.md
      end-to-end-request-flows.md
      shared-apis-and-contracts.md
      shared-data-and-messaging.md
      aws-resource-landscape.md
      infrastructure-topology.md
      ownership-map.md
      security-and-exposure.md
      unresolved-dependencies.md
      evidence-and-coverage.md
      diagrams/
        system-context.mmd
        service-dependency-map.mmd
        cross-service-flows.mmd
        aws-landscape.mmd
        deployment-topology.mmd
    graph/
      workspace-graph.json
      repositories.json
      cross-repository-relationships.json
      shared-resources.json
      evidence-index.json
      manifest.json
      schema-version.json
    metadata/
      workspace.json
      analysis-summary.json
      coverage.json
      repository-inventory.json
~~~

Rules:

- Child repositories own their detailed code graphs.
- The workspace graph references stable repository and entity IDs and owns cross-repository relationships.
- Do not duplicate complete child graphs unnecessarily.
- Repository graphs remain independently usable.
- Do not mistake arbitrary folders for repositories.
- Provide explicit discovery rules and configuration overrides.

If the parent directory is itself a Git repository, automatically separate its own repository knowledge from workspace knowledge:

~~~text
.archgraph/
  repository/
    docs/
    graph/
    metadata/
  workspace/
    docs/
    graph/
    metadata/
~~~

If the parent is not a Git repository, generate workspace output without inventing a repository identity.

## 9. Incremental regeneration and dependency invalidation

Incremental regeneration is the default.

Use both:

1. Git-based change detection for commit workflows.
2. Filesystem fingerprints for uncommitted and untracked local changes.

Detect additions, modifications, deletions, renames, configuration changes, API contract changes, build dependency changes and infrastructure declaration changes.

When one repository changes:

1. Identify changed files.
2. Identify affected analyzer outputs and graph entities.
3. Re-analyze affected components.
4. Update that repository's graph and documentation.
5. Reconcile affected cross-repository relationships.
6. Update affected workspace graph and documentation.
7. Leave unaffected repository knowledge untouched.
8. Update source manifests.

A public API or contract change may invalidate consumers in other repositories even when their source files have not changed. Implement dependency-aware invalidation; do not equate incremental with only scanning changed files.

Provide an explicit full rebuild for recovery and verification. Do not perform a full workspace rebuild during every normal incremental scan.

## 10. Freshness and stale-knowledge policy

Every knowledge base must have a source manifest recording:

- Repository and workspace identity.
- Git commit where available.
- Source and relevant configuration fingerprints.
- Analyzer versions.
- Graph schema version.
- Generation timestamp.
- Supported analysis capabilities.
- Coverage and limitations.

Use timestamps for quick checks and source manifests for correctness validation.

AI retrieval must expose freshness status.

Stale retrieval policy:

1. If fresh, return the knowledge.
2. If stale, attempt incremental regeneration.
3. If regeneration succeeds, return refreshed knowledge.
4. If regeneration fails, return existing knowledge only when explicitly marked stale, including last analyzed commit, affected files and limitations.
5. Never label stale knowledge as current.

Protect concurrent regeneration with locking or equivalent coordination, atomic output replacement and snapshot consistency. Concurrent readers must not observe partially written knowledge bases.

## 11. AI-agent knowledge retrieval: central product capability

Support agent-agnostic integration through:

- A local MCP server.
- CLI query commands.
- Machine-readable JSON interfaces.

AI agents must reuse the generated knowledge base rather than scan all source repositories for every question.

Implement layered retrieval:

1. Retrieve a compact workspace or repository overview and index.
2. Identify relevant services and entities.
3. Retrieve relevant dependencies, APIs, data flows, AWS resources or journeys.
4. Retrieve supporting evidence and source references.
5. Retrieve related documentation or diagrams only as needed.

Support both complete graph exports and focused context retrieval. Focused retrieval is the default for agent queries.

Expose stable entity IDs, relationship IDs, evidence references, freshness status and source locations.

Provide read-only MCP tools for workspace/repository overview, service lookup, dependency exploration, impact analysis, API discovery, AWS resource discovery, request journeys, architecture search, ownership lookup, security/exposure discovery, relationship explanation, evidence retrieval and freshness validation.

MCP must be read-only by default. Do not expose arbitrary filesystem access, shell execution or unrestricted file reads.

## 12. Human-readable documentation

Generate both:

1. An architecture handbook.
2. Detailed technical reference documents.

Organize content in layers:

- Workspace overview.
- Repository overview.
- Service architecture.
- Dependencies and interactions.
- APIs and contracts.
- Data access and messaging.
- AWS integrations.
- Infrastructure and deployment.
- Execution/request flows.
- Source-level evidence and limitations.

Each generated document should include title, purpose, scope, generation timestamp, source repository/workspace, related documents, relevant diagrams and known limitations.

Use navigable indexes and links. Documentation must be understandable without an AI agent and must be generated from the canonical graph.

## 13. Runtime observer

Implement an optional native local process on the Docker host, including WSL2/Linux. Do not require an observer container.

Docker metadata discovery may include containers, images and digests, labels, Compose metadata, networks, aliases, ports and lifecycle events.

Docker topology is not proof of application-level calls.

Support optional local OpenTelemetry ingestion for HTTP/REST, gRPC, database-client and messaging interactions. Do not automatically modify application containers or require instrumentation to be present.

The observer must never start, stop, restart or modify application containers.

Target commands:

~~~bash
archgraph runtime start --background
archgraph runtime status
archgraph runtime snapshot
archgraph runtime stop
~~~

Process events in memory. At snapshot time, reconcile and aggregate sanitized relationship metadata into graph JSON.

Do not persist raw traces, individual event histories, request/response bodies, credentials, authorization headers, cookies, raw SQL parameter values or message payloads.

Aggregated runtime metadata may include source/target IDs, evidence classification, observation count, protocol, first/last observed timestamps, capture window and instrumentation source.

An observation count is evidence of activity, not proof of success or reliability.

Report unsnapshotted observations. Implement PID management, health checks, graceful shutdown and safe process recovery. With no raw event persistence, do not promise recovery of in-memory observations after a crash.

## 14. Security and privacy

ArchGraph has read-only analysis permissions.

Permitted writes are limited to:

- ArchGraph-owned generated output.
- Explicitly approved Git hook installation.
- Explicitly requested export bundles.

Never modify application source or infrastructure.

Record secret references and metadata, never actual secret values. Redact credentials, tokens, passwords, private keys and sensitive values from graphs, documentation, logs, MCP responses and bundles.

Operate locally. No external telemetry, cloud APIs, SaaS or remote runtime dependencies.

Use supported local Docker interfaces. Never expose an unauthenticated Docker TCP API.

Security/exposure analysis identifies architecture exposure and suspicious configuration patterns. Do not claim to replace SAST, secrets scanning, dependency vulnerability scanning or cloud security posture tools.

## 15. Architecture queries and graph diff

Implement:

- Dependency explorer.
- Change impact analysis.
- AWS resource explorer.
- Request journey tracing.
- Architecture search.
- Ownership explorer.
- Security and exposure explorer.
- Git commit-to-commit graph diff.

Graph diff must compare arbitrary Git commits and report added, removed and changed entities and relationships, affected services and evidence differences.

All results must be evidence-backed.

## 16. CLI, TUI and configuration

Support a command-based CLI and an optional interactive terminal UI. Both must use the same application services and business logic.

Target commands include:

~~~text
archgraph init
archgraph scan
archgraph scan --incremental
archgraph scan --full
archgraph status
archgraph validate
archgraph report
archgraph graph
archgraph graph diff
archgraph explain
archgraph search
archgraph dependencies
archgraph impact
archgraph journey
archgraph aws
archgraph security
archgraph export
archgraph knowledge
archgraph runtime
archgraph mcp serve
archgraph git install-hooks
archgraph git uninstall-hooks
~~~

Refine names and flags consistently. Every command must provide useful help, structured errors, meaningful exit codes, human-readable output and JSON output where appropriate.

Use YAML configuration, including archgraph.yaml and service-level codeway.yaml.

codeway.yaml is authoritative for service identity when present. ArchGraph may discover and suggest missing metadata but must not silently invent authoritative identity.

## 17. Git integration

Support manual regeneration, opt-in pre-commit hooks and CI stale-artifact validation.

Hooks must be installed explicitly, preserve existing hooks and Git configuration, stage only ArchGraph-owned files and never stage unrelated developer changes.

CI must validate deterministic generation and fail when committed artifacts are stale. CI must not push, rewrite branches or automatically commit regenerated artifacts.

Do not perform automatic commits or pushes.

## 18. Exports and sharing

Support:

- Canonical JSON.
- Markdown.
- Mermaid.
- Graphviz DOT.
- Portable ZIP knowledge bundle.

Example:

~~~bash
archgraph export --format bundle --output architecture-bundle.zip
~~~

The bundle must contain graph data, human-readable documentation, diagrams and a manifest containing schema version, source commit, repository identities, analyzer versions, generation timestamp, freshness status, coverage and limitations.

The bundle must be portable and useful without the original source repositories. Never include secrets or raw runtime traces.

## 19. Quality, tests and performance

Mandatory test categories:

- Unit tests.
- Integration tests.
- End-to-end tests.

Target minimum Go test coverage: 90%. Coverage must not be gamed with meaningless tests. Prioritize assertions about relationship correctness, evidence, identity, freshness and safe behavior.

Create realistic fixtures for Go, Java/Spring Boot, Python, multi-repository workspaces, Docker Compose, Kubernetes, Terraform, AWS integrations, API contracts, runtime observations and incremental regeneration.

Target workspace scale: 25–100 microservices. Provide reproducible performance benchmarks.

Test deterministic generation, stable IDs, schema migrations, cross-repository relationships, incremental invalidation, stale knowledge handling, atomic writes, concurrent reads/regeneration, MCP response correctness, secret redaction, runtime shutdown and bundle integrity.

Do not claim support for a framework or format without corresponding fixtures and tests.

## 20. Distribution

V1 platforms:

- Ubuntu/Linux, including WSL2.
- Native Windows.

Distribution:

- APT packages for Ubuntu/Debian.
- WinGet packages for Windows.
- Stable release channel only.

Version checks are allowed. Automatic upgrades are not.

Provide reproducible builds, version metadata and release checksums. Docker packaging is optional. Static analysis must work without Docker. Do not claim macOS as a supported v1 platform.

## 21. Implementation governance

Work in explicit phases. Before each phase, present:

1. Purpose and scope.
2. Architectural changes.
3. Packages/files to create or modify.
4. Dependencies.
5. Acceptance criteria.
6. Tests.
7. Risks and unresolved decisions.

Wait for user approval before implementing that phase.

After approval, implement only that phase. At its completion:

1. Run tests and formatting.
2. Run static analysis.
3. Check coverage.
4. Verify acceptance criteria.
5. Update documentation.
6. Report changed files, test results and limitations.
7. Stop and wait for approval before proceeding.

Do not silently proceed to the next phase or make unrelated refactors. Do not install software, change global machine settings or execute destructive commands without explicit approval. Do not modify GitHub or push commits unless explicitly instructed.

## 22. Definition of done

The product is complete when:

1. A single repository can be scanned and receives an independently usable knowledge base.
2. A multi-repository workspace produces child knowledge bases and a consolidated workspace knowledge base.
3. Parent-repository and workspace outputs are separated when the parent is also a Git repository.
4. Incremental changes update affected knowledge and cross-repository relationships correctly.
5. Freshness validation detects stale outputs.
6. AI agents retrieve knowledge through MCP, CLI and JSON without rescanning the entire workspace.
7. Layered retrieval returns focused context with evidence and freshness.
8. Developers can read and share the architecture handbook and technical references.
9. AWS and infrastructure relationships are evidence-backed.
10. Runtime observation enriches the graph without raw trace persistence.
11. Portable knowledge bundles are generated and validated.
12. Git integration protects developer workflows.
13. Tests validate correctness on realistic multi-service fixtures.
14. Supported Linux and Windows distribution workflows are verified.

The ultimate success criterion:

A developer submits a new engineering requirement to an AI coding agent. The agent retrieves existing ArchGraph knowledge, understands the relevant architecture and dependencies, and inspects source files only when targeted evidence is needed.

ArchGraph must become reusable architecture memory for the engineering organization.

**Start with repository assessment and the proposed implementation roadmap. Do not begin implementation until the user explicitly approves Phase 1.**
