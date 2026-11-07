# Docker Runtime Design

## Operating contract

ArchGraph v1 is a one-shot container:

1. Start the container.
2. Discover Git repositories under the mounted workspace.
3. Read source and configuration through read-only mounts.
4. Run deterministic MVP discovery for service identity, source inventory and selected code/resource patterns.
5. Export architecture.json and architecture.md to the mounted output directory.
6. Exit with a documented status code.
7. The caller removes the container with docker run --rm.

There is no background daemon, HTTP API or persistent server in v1.

## Windows PowerShell

    docker build --pull -t archgraph:dev .

    New-Item -ItemType Directory -Force .\archgraph-output | Out-Null

    docker run --rm -v "D:\work\repositories:/workspace:ro" -v "${PWD}\archgraph-output:/output" archgraph:dev scan --path /workspace --output /output

## Multiple repositories

Mount a parent workspace containing multiple Git repositories. ArchGraph recursively discovers nested Git repositories and stops descending after finding a repository boundary.

## Security and reproducibility

- Run as a non-root user.
- Mount source repositories read-only.
- Mount output separately with only the permissions required to write artifacts.
- Do not pass AWS credentials; v1 does not query live AWS accounts.
- Do not make LLM calls during discovery.
- The current MVP has no external runtime dependencies beyond the Go binary and CA certificates.
- Docker build uses versioned Go and Alpine image tags. Production releases should pin base images by OCI digest and retain an SBOM.
- The image can scan offline after it has been built.

## Exit status

- 0: scan completed without warnings.
- 1: execution, filesystem or configuration error.
- 2: command usage error.
- 3: output generated but discovery is incomplete or warnings exist.

## Deferred integrations

Graphify invocation and graph normalization are target architecture items, not part of the current MVP binary. They will be integrated only after pinning the Graphify output schema and adding acceptance fixtures.

A long-running server/API/MCP mode is intentionally deferred. It must not be added until graph query and journey execution interfaces are stable.
