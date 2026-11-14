# ArchGraph

**Local-first, deterministic architecture intelligence for software repositories and multi-repository workspaces.**

ArchGraph builds a reusable, evidence-backed knowledge base from software repositories. It discovers architecture facts, reconciles them into a versioned canonical graph, and generates human-readable documentation. Developers and AI coding agents can then retrieve focused architectural context instead of repeatedly rediscovering the same dependencies from source code.

> **Core idea:** Build architecture knowledge once, keep it fresh, and make it reusable. ArchGraph discovers and explains the system; developers and AI agents use that knowledge to reason about changes.

## Why ArchGraph?

When an engineer or AI coding agent works across a large system, it often has to rediscover basic facts: which service owns an endpoint, which component calls another, who publishes or consumes an event, which database resources are accessed, and what may be affected by a change.

That repeated discovery consumes engineering time, context-window capacity and potentially model tokens and calls. ArchGraph aims to move repeatable architecture discovery into a deterministic local tool:

1. Analyze repositories using language- and framework-aware analyzers.
2. Normalize discovered facts into a language-neutral canonical graph.
3. Attach evidence and provenance to relationships, preserving uncertainty.
4. Generate navigable architecture documentation and diagrams.
5. Expose full-graph exports and focused retrieval interfaces for developers and AI agents.
6. Refresh affected knowledge as source repositories change.

This may reduce repeated source inspection and prompt construction. **Token, cost and latency savings have not yet been benchmarked.** Results will depend on analyzer coverage, graph freshness, retrieval selectivity and the cost of maintaining the knowledge base. ArchGraph does not replace source inspection or AI reasoning.

## Product principles

- **Local-first:** core analysis and retrieval run locally. No mandatory cloud service, SaaS backend or LLM.
- **AI-free discovery:** scanning, parsing, graph construction, reconciliation and validation do not call an LLM.
- **Evidence-first:** every source-derived relationship carries provenance. Confirmed, declared, inferred, unresolved, conflicting and runtime-observed evidence remain distinct.
- **Canonical graph ownership:** ArchGraph owns a versioned, language-neutral graph schema. Analyzer adapters emit normalized facts; the core validates and reconciles them.
- **Reusable knowledge:** repository-level knowledge remains independently useful, while workspace knowledge connects repositories and cross-service relationships.
- **Human-readable outputs:** architecture handbooks, technical references and diagrams are first-class outputs, not just graph exports.
- **Incremental by default:** changed files and dependency invalidation drive refreshes; full scans remain available for verification and recovery.
- **Safe by default:** source repositories are read-only. Writes are limited to ArchGraph-owned outputs and explicitly approved integrations.
- **Honest coverage:** unresolved behavior and unsupported constructs are reported. ArchGraph does not claim complete discovery or 100% accuracy.

## What ArchGraph is designed to provide

- **Repository and workspace inventory:** discover repositories, services and their identities, including `codeway.yaml` metadata when present.
- **Code relationships:** packages, modules, types, functions, imports, calls and framework-specific relationships where analyzers can establish them.
- **API and event intelligence:** endpoints, contracts, clients, producers, consumers, queues and topics.
- **Data access and lineage:** database/table/collection access and field-level lineage where statically identifiable.
- **AWS and infrastructure discovery:** resource usage and declarations from source, dependencies, configuration and supported infrastructure definitions.
- **Impact analysis:** graph diffs, dependency traversal and affected-service/journey exploration.
- **Evidence and explainability:** source locations, analyzer identity, evidence classification and limitations.
- **Agent integration:** focused CLI/JSON retrieval and read-only MCP tools as capabilities are implemented.
- **Optional runtime enrichment:** local Docker metadata and OpenTelemetry observations, kept separate from static evidence.

A dependency declaration or SDK import is not proof that an application actually uses a resource. Docker network membership is not proof of an application call. ArchGraph must not silently promote clues into confirmed relationships.

## Architecture at a glance

ArchGraph is a **native Go CLI first**. Docker is an optional packaging and execution method, not a requirement for static analysis.

The architecture uses a hybrid analyzer strategy:

- **Go:** native semantic analysis using the Go toolchain ecosystem, including packages, AST, types, SSA and call-graph capabilities where appropriate.
- **Java / Spring Boot:** compiler-aware analysis plus Spring-specific discovery. A bundled, versioned local adapter process may be used where Java tooling requires it.
- **Tree-sitter:** shared syntax extraction and fallback layer. It does not replace compiler/type resolution where semantic analysis is required.
- **Other adapters:** Python, TypeScript, JavaScript and React are secondary/future coverage targets. React is a framework adapter, not a programming language.

External analyzers may be evaluated as optional adapters or implementation references. They do not own ArchGraph's canonical schema, evidence policy, identity rules or validation.

### Canonical graph and evidence

JSON is the canonical machine-readable graph. Markdown documentation and Mermaid diagrams are derived from that graph.

The graph uses typed entities and relationships, stable logical identities, schema versioning and explicit migrations. Relationship evidence must distinguish:

- `STATIC_CONFIRMED`
- `DECLARED`
- `RUNTIME_OBSERVED`
- `INFERRED`
- `UNRESOLVED`
- `CONFLICTING`

Every relationship should retain enough provenance to explain where it came from, how it was discovered and what limitations apply.

### Repository and workspace knowledge

Each analyzed repository owns its detailed generated knowledge under `.archgraph/`. A parent workspace maintains a consolidated graph and documentation for cross-repository relationships. If the parent directory is itself a Git repository, its repository knowledge and workspace knowledge are stored separately.

The graph and documentation are versioned generated artifacts. Manually maintained documentation must not be overwritten.

## Current implementation status

ArchGraph is an evolving project. The existing MVP provides an initial deterministic inventory and candidate-discovery workflow; it is **not yet the complete semantic architecture intelligence system described by the target architecture**.

The current MVP includes:

- Recursive discovery of Git repositories under a workspace.
- Service/Lambda identity candidates from common `codeway.yaml` / `codeway.yml` keys.
- Inventory of Go, Java, Python, TypeScript and JavaScript source files.
- Initial deterministic candidates for AWS operations, API routes, database access and Terraform resources.
- Evidence attached to discovered nodes and inferred source-to-service containment edges.
- JSON graph and Markdown inventory output.
- One-shot Docker execution with read-only source mounts.

The following are target capabilities and must not be treated as complete unless verified in the implementation and tests:

- Deep Go and Java/Spring semantic analysis.
- Reliable cross-repository relationship resolution.
- Comprehensive AWS wrapper and resource mapping.
- Full database schema and field-level lineage discovery.
- Security scanner integrations and correlated security analysis.
- Graph diff, dependency-aware impact analysis and journey execution.
- Incremental invalidation, freshness validation and atomic regeneration.
- Read-only MCP tools and layered agent-context retrieval.
- Native local runtime observation and OpenTelemetry enrichment.
- Production packaging, Windows support and release automation.

Coverage must be reported by repository, language/framework, infrastructure type and evidence class. Unsupported constructs and unresolved behavior must remain visible.

## Getting started

### Requirements

- Go toolchain compatible with the version declared in `go.mod`.
- Git.
- Optional: Docker, if you want to use the container workflow.

### Build and test

```bash
git clone https://github.com/PranitSakunde1/archgraph.git
cd archgraph

go test ./...
go vet ./...
go build ./cmd/archgraph
```

### Run the current MVP

```bash
go run ./cmd/archgraph scan --path ./fixtures --output ./out
```

Check `archgraph --help` for the commands currently implemented. The command list and flags will evolve as the CLI is developed.

## Optional Docker workflow

Docker is an optional one-shot execution path. Mount source repositories read-only and mount a separate output directory for generated artifacts.

Example for PowerShell:

```powershell
docker build --pull -t archgraph:dev .

New-Item -ItemType Directory -Force .\\archgraph-output | Out-Null

docker run --rm `
  -v "D:\\work\\repositories:/workspace:ro" `
  -v "${PWD}\\archgraph-output:/output" `
  archgraph:dev scan --path /workspace --output /output
```

The scan container should scan, export and exit. Static analysis must also work natively without Docker.

## Target language and framework coverage

| Priority | Technology | Intended analysis |
|---|---|---|
| First-class | Go | Compiler/toolchain-aware semantic analysis |
| First-class | Java / Spring Boot | Compiler-aware and Spring-specific analysis |
| Secondary / future | Python | Language adapter |
| Secondary / future | TypeScript and JavaScript | Syntax and framework-aware adapters |
| Framework adapter | React | Components, hooks, routes and API clients |

A technology is not considered supported merely because files are inventoried. Meaningful support requires analyzer behavior, fixtures, tests and honest coverage reporting.

## Planned developer and agent workflows

The target CLI surface includes commands in these areas:

- **Discovery:** init, scan, status and report.
- **Validation:** validate freshness, graph structure, coverage and generated artifacts.
- **Exploration:** graph, search, dependencies, AWS, security and explain.
- **Change analysis:** graph diff and impact.
- **Journeys:** inspect and run explicitly approved journeys.
- **Knowledge:** export full graph or retrieve focused context.
- **Integration:** read-only MCP server and opt-in Git hooks.
- **Runtime:** separate start/status/snapshot/stop workflow for optional local observation.

These are target capabilities. Use CLI help and the implementation to determine which commands are currently available.

## Runtime observation (optional)

Static analysis is the foundation and does not require runtime instrumentation.

The optional runtime observer is designed to run as a native local process on the Docker host, including WSL/Linux. Docker metadata provides container and network topology evidence only; it does not prove HTTP, RPC, database or messaging interactions. Optional OpenTelemetry instrumentation can provide application-interaction evidence.

Runtime collection and processing remain local. Sensitive payloads, credentials, request/response bodies and raw SQL parameter values must not be persisted by default. Runtime observations remain distinct from static evidence. The observer must not start, stop, restart or modify application containers.

## Git and safety model

- Source scanning is read-only by default.
- Generated artifacts are written only to ArchGraph-owned paths.
- Git hooks are opt-in and must preserve existing hooks and configuration.
- Hooks may stage only ArchGraph-owned files, never unrelated developer changes.
- CI validates committed generated artifacts and fails when they are stale or invalid.
- CI does not push commits, rewrite branches or automatically commit regenerated output.
- No automatic Git commit or push.
- No live AWS account queries are part of v1.

## Documentation and decisions

- [Architecture Grill: Questions and Answers](docs/decisions/architecture-grill-questions-and-answers.md)
- [Architecture Grill: Final Decision Register](docs/decisions/architecture-grill-decisions.md)
- [ADR-001: Hybrid Evidence-First Analysis](docs/decisions/ADR-001-hybrid-evidence-first-analysis.md)
- [Architecture blueprint](docs/architecture/architecture-blueprint.md)
- [Docker runtime design](docs/architecture/docker-runtime.md)
- [Production build prompt](docs/implementation/ARCHGRAPH-MASTER-BUILD-PROMPT.md)
- [Product use cases by persona](docs/product/persona-use-cases.md)

## First engineering milestone

Before broadening implementation, validate representative real-world inputs:

1. Actual `codeway.yaml` examples and service identity conventions.
2. Internal AWS wrapper implementations and SDK usage patterns.
3. Infrastructure/resource definitions used by target repositories.
4. Representative Go and Java/Spring Boot services.
5. At least one end-to-end business journey and its cross-service dependencies.

Turn these into versioned acceptance fixtures. Use them to measure analyzer coverage and relationship correctness before introducing strict quality gates.

## Contributing and implementation governance

Implementation is phased. Each phase must begin with a proposed scope, architecture changes, affected packages/files, dependencies, acceptance criteria, tests and risks. Wait for explicit approval before implementing a phase. At phase completion, run formatting, tests, static analysis and coverage checks, update documentation, report limitations and stop for approval before proceeding.

See the [production build prompt](docs/implementation/ARCHGRAPH-MASTER-BUILD-PROMPT.md) and the decision records for the full engineering constraints.
