# ArchGraph

Deterministic, local-first architecture discovery for multi-repository workspaces.

ArchGraph scans a workspace, discovers nested Git repositories and service identities from recursively found codeway.yaml/codeway.yml files, and emits a JSON graph plus a Markdown inventory report. The current scanner uses deterministic source-pattern candidates. Those matches are evidence for further analysis, not proof of a complete runtime relationship.

## Current MVP

- Recursive discovery of Git repositories under a workspace.
- Service/Lambda identity candidates from common codeway YAML keys.
- Inventory of Go, Java, Python, TypeScript and JavaScript source files.
- Initial deterministic candidates for AWS operations, API routes, database access and Terraform resources.
- Evidence attached to discovered nodes and inferred source-to-service containment edges.
- JSON and Markdown output.
- One-shot Docker execution with read-only source mounts.

Not yet implemented: Graphify invocation/normalization, deep AST analysis, cross-repository edge resolution, full resource mapping, schema-level database discovery, field lineage, security scanner integrations, graph diff/impact analysis and journey execution.

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
