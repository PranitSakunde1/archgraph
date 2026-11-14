# Docker Runtime Design

## Static scan operating contract

ArchGraph static scanning is a one-shot CLI workflow. The CLI can run natively; an optional Docker image packages the same scan workflow. The scan container is not the runtime observer.

1. Start the container.
2. Discover Git repositories under the mounted workspace.
3. Read source and configuration through read-only mounts.
4. Run deterministic MVP discovery for service identity, source inventory and selected code/resource patterns.
5. Export architecture.json and architecture.md to the mounted output directory.
6. Exit with a documented status code.
7. The caller removes the container with docker run --rm.

There is no background daemon or API server for static scanning. Runtime observation is handled separately by a native host process described below.

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

 
## Native local runtime observer

The runtime observer is a separate native process on the Docker host, including WSL/Linux. It is not run as a container and does not require mounting the Docker socket into any container.

It connects to the local Docker daemon to observe container lifecycle and declared topology metadata. Docker metadata identifies running containers, images, Compose labels, networks and ports. These facts do not prove application-level requests or data access.

Actual HTTP/RPC, database and messaging relationships require optional OpenTelemetry instrumentation and a local OTLP receiver/export path. The observer correlates telemetry with local container/service identity and writes runtime evidence to local storage. No telemetry is sent to an external destination by default.

Do not persist request/response bodies, credentials, authorization headers, cookies, raw SQL parameter values or message payloads by default. Redact sensitive attributes, configure retention, and keep raw traces separate from committed architecture artifacts.

Suggested lifecycle commands (target interface, not implemented yet):

    archgraph runtime start
    archgraph runtime status
    archgraph runtime snapshot
    archgraph runtime stop

The observer must not start, stop, restart or otherwise modify application containers. Static scanning and export remain usable while the observer is stopped.
