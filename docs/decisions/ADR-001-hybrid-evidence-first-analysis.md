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

The primary interface is one native Go CLI. Docker is optional for reproducible packaging and isolated execution. No always-on daemon or server is required for core scanning.

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

### 7. Generated context and Git workflow

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

## Acceptance criteria

- Go and Spring Boot fixture repositories produce expected graph snapshots.
- Every published source-derived relationship has evidence and provenance.
- Unresolved and unsupported constructs are reported, not silently omitted.
- Repeated scans of the same source revision and pinned analyzer versions produce equivalent canonical graph content.
- Runtime evidence is optional and separately classified.
- Native CLI works without Docker; Docker packaging remains supported.
- CI detects stale generated documentation without modifying repository branches.
