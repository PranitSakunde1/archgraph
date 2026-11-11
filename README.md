# ArchGraph

Deterministic, local-first architecture discovery for multi-repository workspaces.

ArchGraph is being built to give AI agents a reusable, evidence-backed view of a software system. Instead of asking an LLM to repeatedly inspect repositories and reconstruct service relationships for every task, ArchGraph aims to discover and maintain those relationships once, then expose focused architectural context for queries, impact analysis and engineering workflows.

> **Core idea:** Let ArchGraph handle architecture discovery and relationship mapping; let AI focus on reasoning, trade-offs and decisions.

## Why ArchGraph?

When an AI agent works across a large codebase, it may need to rediscover the same information repeatedly: which service owns an endpoint, who publishes or consumes an event, which resources are used, and what downstream components may be affected by a change. That discovery consumes context-window capacity, tokens, model calls and wall-clock time before the actual engineering question is addressed.

ArchGraph is designed to move repeatable discovery out of the LLM loop:

1. Scan repositories and extract architecture signals using deterministic analysis.
2. Normalize those signals into a canonical graph of repositories, services, APIs, events, resources and data dependencies.
3. Preserve evidence and distinguish confirmed, inferred and unresolved relationships.
4. Let tools and AI agents query only the relevant part of the graph for a task.

This can reduce repeated context construction and unnecessary LLM work. **Actual token, cost and latency savings are not yet measured and will depend on graph coverage, freshness, query selectivity and the cost of maintaining the graph.** ArchGraph does not eliminate the need for AI reasoning or direct source inspection.

## What makes this useful?

- **Reusable architectural context:** discover relationships once and reuse them across multiple questions and sessions.
- **Cross-repository visibility:** connect services, APIs, event producers/consumers and resource usage across repository boundaries.
- **Focused retrieval:** provide a task-specific subgraph instead of sending large amounts of unrelated source code to an LLM.
- **Impact analysis foundation:** traverse dependency relationships to identify potentially affected services and journeys.
- **Evidence and explainability:** attach source references and discovery methods to relationships; keep confirmed, inferred and unresolved facts separate.
- **Controlled graph traversal:** bound traversal by depth, scope and visited-node tracking. Detect and classify cycles so traversal terminates safely, while preserving intentional architecture loops as relationships rather than treating every cycle as an error.
- **AI-independent discovery:** scanning, parsing, graph construction and validation are designed to run without LLM calls.

## How this can save AI resources

| AI-only workflow | ArchGraph-assisted workflow |
|---|---|
| Repeatedly inspect source files to rediscover dependencies | Reuse previously discovered graph relationships |
| Reconstruct cross-service flows in the prompt | Traverse the relevant producer/consumer or dependency path |
| Send broad repository context to the model | Retrieve a smaller, task-specific subgraph with evidence |
| Spend model reasoning on basic architecture discovery | Reserve model reasoning for intent, ambiguity and decisions |

The expected savings are less repeated source ingestion, fewer redundant discovery steps and potentially fewer tokens or model calls. This is an architectural hypothesis, not a benchmark result. Graph generation, refreshes, storage and query execution also consume compute; the goal is to make that work reusable and cheaper than repeatedly rebuilding the same context.

## Trade-offs and limitations

ArchGraph adds an architecture index that must be built, validated and kept current. It is not a free replacement for reading code.

| Pros | Cons / risks |
|---|---|
| Reusable context across tasks and AI sessions | Requires initial scanning and ongoing refreshes |
| Better visibility across repositories and services | Incomplete parsers can miss relationships |
| Evidence-backed, inspectable relationships | Incorrect inference can mislead downstream analysis |
| Focused graph retrieval can reduce prompt size | Large or noisy graphs need filtering and traversal limits |
| Repeatable dependency and cycle analysis | Dynamic dispatch, reflection and runtime behavior can be hard to resolve statically |
| No LLM required for core discovery | Graph maintenance, storage and query costs still exist |
| Can support impact analysis and journey selection | Impact results are only as complete as the graph and its evidence |

**Important:** a graph edge is not automatically a runtime truth. Naming conventions, dependency declarations or matching resource names are clues, not proof. ArchGraph must preserve evidence, confidence and uncertainty rather than silently turning guesses into confirmed relationships.

## Current MVP

ArchGraph scans a workspace, discovers nested Git repositories and service identities from recursively found codeway.yaml/codeway.yml files, and emits a JSON graph plus a Markdown inventory report. The current scanner uses deterministic source-pattern candidates. Those matches are evidence for further analysis, not proof of a complete runtime relationship.

- Recursive discovery of Git repositories under a workspace.
- Service/Lambda identity candidates from common codeway YAML keys.
- Inventory of Go, Java, Python, TypeScript and JavaScript source files.
- Initial deterministic candidates for AWS operations, API routes, database access and Terraform resources.
- Evidence attached to discovered nodes and inferred source-to-service containment edges.
- JSON and Markdown output.
- One-shot Docker execution with read-only source mounts.

**Not yet implemented:** Graphify invocation/normalization, deep AST analysis, cross-repository edge resolution, full resource mapping, schema-level database discovery, field lineage, security scanner integrations, graph diff/impact analysis and journey execution. Controlled multi-hop traversal, cycle classification and AI/MCP-facing architecture queries are target capabilities, not claims about the current MVP.

**Graphify status:** Graphify is part of the target architecture as a potential code-level extraction engine, but it is not currently integrated into the MVP. ArchGraph's canonical graph, normalization and correlation layer remain the core product responsibility; Graphify should be an optional/pluggable extraction source rather than a hard runtime dependency.

## Docker workflow

    docker build --pull -t archgraph:dev .

    New-Item -ItemType Directory -Force .\archgraph-output | Out-Null

    docker run --rm -v "D:\work\repositories:/workspace:ro" -v "${PWD}\archgraph-output:/output" archgraph:dev scan --path /workspace --output /output

The workspace can contain multiple nested Git repositories. ArchGraph discovers them automatically. The output directory is mounted separately so source repositories remain read-only.

## Local development

    go test ./...
    go vet ./...
    go run ./cmd/archgraph scan --path ./fixtures --output ./out

## Documentation

- [ArchGraph by persona: Product Manager, Business Analyst, Tech Lead, QA Engineer, QA Lead and Developer](docs/product/persona-use-cases.md)
- [Architecture blueprint](docs/architecture/architecture-blueprint.md)
- [Architecture grill decision register](docs/decisions/architecture-grill-decisions.md)
- [Docker runtime design](docs/architecture/docker-runtime.md)

## Target language/framework coverage

Go, Java, Python, React, TypeScript and JavaScript. React requires framework-specific component, hook, route and API-client discovery in addition to TypeScript/JavaScript analysis. The full target coverage is planned; current MVP matching is intentionally shallow.

## Design constraints

- No LLM calls during discovery, parsing, graph generation or validation.
- One-shot Docker execution: scan, export, exit.
- Read-only source mounts by default.
- No live AWS account queries in v1.
- Application architecture is in scope; deployment topology is not.
- Confirmed, inferred and unresolved relationships must remain distinguishable as discovery matures.

## First engineering milestone

Validate real codeway.yaml examples, internal AWS wrapper implementations, resource definitions and one end-to-end business journey. These become acceptance fixtures before broadening discovery or introducing CI gates.
