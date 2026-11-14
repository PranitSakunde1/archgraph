# ADR-001: Hybrid, evidence-first architecture analysis

- **Status:** Accepted
- **Date:** 2026-11-15
- **Decision owners:** ArchGraph maintainers
- **Scope:** ArchGraph core, language analyzers, runtime evidence, CLI and generated artifacts

## Context

ArchGraph currently has a Go CLI with workspace/repository discovery, lightweight source inventory and deterministic pattern-based candidates. That implementation is not yet a semantic Java or Go architecture analyzer. The product goal is a single local-first tool that generates reusable architecture context for Java Spring Boot and Go microservices without requiring LLMs or cloud services.

A single generic parser cannot provide compiler-level semantic resolution for both languages. Runtime observations are useful but cannot enumerate all possible execution paths. External graph generators can accelerate extraction, but must not dictate the product's graph semantics or evidence policy.

## Decision

### 1. One product

Continue development in this ArchGraph repository. The Go core owns CLI orchestration, repository discovery, canonical graph schema, stable identity, evidence, relationship reconciliation, validation, graph diff and export.

### 2. Hybrid parser architecture

- **Go:** native Go analysis libraries for package loading, AST, type information, SSA and statically resolvable call graphs.
- **Java:** Java/compiler-aware analysis with Spring-specific extraction for controllers, routes, beans, injection, repositories, entities, messaging and build/module metadata.
- **Tree-sitter:** shared syntax extraction and fallback for supported languages. It does not substitute for semantic/type resolution.
- **External tools:** Graphify may be evaluated as an optional adapter or reference. It is not the canonical graph engine and is not a mandatory dependency.

### 3. Optional runtime evidence

Static scanning and exports must work without runtime instrumentation. OpenTelemetry traces may be ingested as optional runtime evidence. Runtime-observed relationships remain separately classified and include environment, time/window and trace provenance. Absence of a runtime observation is not proof that a relationship does not exist.

### 4. Deployment

The primary static-analysis interface is one native Go CLI. Docker is optional for reproducible packaging and isolated execution. Static scanning remains one-shot. Runtime observation is a separate native local process on the Docker host, including WSL/Linux, not an observer container. The observer must not require Docker socket access from a container.

### 5. Initial language scope

Java Spring Boot and Go are first-class targets. Other languages may be added through adapters. React is treated as a framework adapter over JavaScript/TypeScript.

### 6. Canonical graph and evidence contract

ArchGraph owns a versioned language-neutral graph schema. Adapters emit facts into that schema. Every published source-derived edge carries provenance, including repository revision, file and source range where available, extractor/version, method and evidence classification.

Evidence classifications must distinguish at least:
- `STATIC_CONFIRMED`
- `DECLARED`
- `RUNTIME_OBSERVED`
- `INFERRED`
- `UNRESOLVED`
- `CONFLICTING`

No relationship may be silently promoted from inferred, unresolved or runtime-observed to statically confirmed. Unknown behavior remains visible. The product must not claim 100% completeness; it must report supported coverage, unresolved constructs and analysis limitations.

### 7. Local Docker runtime observation

The runtime observer runs as a native host process and discovers local containers through the Docker host API/CLI integration. Docker events and inspection provide container lifecycle and declared topology metadata only; they do not prove application-level communication.

Actual HTTP/RPC, database and messaging relationships are captured through optional OpenTelemetry instrumentation and a local OTLP receiver/export path. Runtime collection, correlation, storage and graph updates stay on the local machine. No SaaS backend, cloud service, LLM or external telemetry destination is required.

Static evidence, Docker topology evidence and runtime-observed application relationships remain separate. Runtime observations carry capture time/window, environment, service/container identity and trace correlation metadata where available. Missing instrumentation is an observability gap, not evidence that a relationship does not exist.

Do not persist request/response bodies, credentials, authorization headers or raw SQL parameter values by default. Redact sensitive attributes and provide configurable local retention/deletion. Static scanning must work independently when the observer is stopped.

### 8. Generated context and Git workflow

The target service repository receives generated Markdown, Mermaid and JSON architecture context. JSON is the canonical machine-readable representation; diagrams and prose are derived from it.

Regeneration occurs in the Git commit workflow for relevant source changes. Hooks must stage only generated files they own and must not stage unrelated user changes. CI performs a full deterministic verification and fails on stale or invalid committed artifacts. CI does not push commits or rewrite branches.

## Consequences

### Positive
- One tool and one graph contract for multiple languages.
- Better semantic precision than regex or syntax-only analysis alone.
- AI-independent, locally reproducible architecture context.
- Runtime evidence enriches the graph without becoming a hard dependency.
- Portable JSON/Markdown/Mermaid artifacts can live beside service code.

### Costs and risks
- Java semantic analysis may require a bundled local Java adapter process and a versioned JSON interface to the Go core.
- Analyzer versions and language toolchains must be pinned and tested.
- Static analysis cannot resolve every reflection, dynamic dispatch, runtime configuration or environment-dependent path.
- Runtime observations are bounded by instrumentation and exercised workloads.
- Graph schema and evidence contracts require compatibility/versioning discipline.

## Implementation sequence

1. Refactor the existing Go CLI and graph model into internal packages.
2. Define versioned node, edge and evidence schemas with stable logical IDs.
3. Implement Go semantic extraction and golden fixtures.
4. Implement Java/Spring extraction and golden fixtures.
5. Add strict validation, deterministic export and CI drift checks.
6. Add optional OpenTelemetry ingestion and evidence reconciliation.
7. Add impact analysis and advanced graph queries after the graph foundation is reliable.
8. Implement the native host runtime observer, Docker lifecycle/topology discovery, local OTLP ingestion, privacy controls and runtime evidence reconciliation.

## Acceptance criteria

- Go and Spring Boot fixture repositories produce expected graph snapshots.
- Every published source-derived relationship has evidence and provenance.
- Unresolved and unsupported constructs are reported, not silently omitted.
- Repeated scans of the same source revision and pinned analyzer versions produce equivalent canonical graph content.
- Runtime evidence is optional and separately classified.
- Native CLI works without Docker; Docker packaging remains supported.
- CI detects stale generated documentation without modifying repository branches.
- Runtime observer runs natively on the host, uses local-only telemetry handling and does not require a privileged observer container.
- Docker topology facts are not misrepresented as observed application calls.
