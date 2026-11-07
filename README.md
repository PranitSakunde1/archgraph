# ArchGraph

Deterministic, local-first architecture intelligence for multi-repository workspaces.

ArchGraph scans a workspace, discovers Git repositories and service identities from `codeway.yaml`, correlates application code with AWS resources, APIs, events, databases and data flows, then emits evidence-backed architecture artifacts.

## Design principles

- No LLM calls in discovery, parsing, graph generation or validation.
- One-shot Docker execution: scan, export artifacts, exit.
- Read-only source mounts by default.
- Local-first graph artifacts; Neo4j is optional future infrastructure.
- Graphify-first code graph extraction, with deterministic ArchGraph enrichers.
- Confirmed, inferred and unresolved relationships remain distinguishable.
- No live AWS account queries in v1.
- Application architecture is in scope; deployment topology is not.

## Docker workflow

```powershell
docker build --pull -t archgraph:dev .

New-Item -ItemType Directory -Force .\archgraph-output | Out-Null

docker run --rm `
  -v "D:\work\repositories:/workspace:ro" `
  -v "${PWD}\archgraph-output:/output" `
  archgraph:dev scan --path /workspace --output /output
```

The workspace can contain multiple nested Git repositories. ArchGraph discovers them automatically. The output directory is mounted separately so source repositories remain read-only.

## Documentation

- [Architecture blueprint](docs/architecture/architecture-blueprint.md)
- [Architecture grill decision register](docs/decisions/architecture-grill-decisions.md)
- [Docker runtime design](docs/architecture/docker-runtime.md)

## Language targets

Go, Java, Python, React, TypeScript and JavaScript. React has framework-specific discovery in addition to TypeScript/JavaScript parsing.

## Implementation status

This repository is the implementation home and design record. The blueprint describes the target system; it is not a claim that every capability is already implemented. Discovery coverage, adapters and gates must be validated against representative repositories before being treated as production-ready.

## First engineering milestone

Validate real `codeway.yaml` examples, Graphify output, internal AWS wrapper implementations, infrastructure resource definitions and one end-to-end business journey. These become acceptance fixtures for deterministic discovery.
