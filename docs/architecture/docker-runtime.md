# Docker Runtime Design

## Operating contract

ArchGraph v1 is a one-shot container:

1. Start the container.
2. Discover repositories under the mounted workspace or use explicit repository paths.
3. Read source and configuration through read-only mounts.
4. Stage temporary extraction copies inside the container when a tool requires writable inputs.
5. Run deterministic discovery and Graphify in code-only mode.
6. Validate and export artifacts to the mounted output directory.
7. Exit with a documented status code.
8. Container is removed by docker run --rm.

There is no background daemon, HTTP API or persistent server in v1.

## Windows PowerShell

    docker build --pull -t archgraph:dev .

    New-Item -ItemType Directory -Force .\archgraph-output | Out-Null

    docker run --rm `
      -v "D:\work\repositories:/workspace:ro" `
      -v "${PWD}\archgraph-output:/output" `
      archgraph:dev scan --path /workspace --output /output

## Multiple explicit repositories

    docker run --rm `
      -v "D:\work\service-a:/repos/service-a:ro" `
      -v "D:\work\service-b:/repos/service-b:ro" `
      -v "${PWD}\archgraph-output:/output" `
      archgraph:dev scan --repo /repos/service-a --repo /repos/service-b --output /output

## Security and reproducibility

- Run as a non-root user.
- Mount source repositories read-only.
- Mount output separately with only the permissions required to write artifacts.
- Do not pass AWS credentials; v1 does not query live AWS accounts.
- Do not make LLM calls during discovery.
- Pin Graphify and runtime versions. Production releases should additionally pin base images by OCI digest and retain an SBOM.
- Network is not needed during scanning once the image is built; image construction may need network unless dependencies are pre-cached.
- A fully hash-locked Python dependency wheelhouse and immutable base-image digests remain prerequisites for bit-for-bit reproducible builds.

## Exit status

- 0: scan completed without warnings.
- 1: execution, filesystem or configuration error.
- 2: reserved for critical architecture validation failure; not yet implemented by the MVP.
- 3: output generated but discovery is incomplete or warnings exist.

## Future, not v1

A long-running server/API/MCP mode is intentionally deferred. It must not be added until graph query and journey execution interfaces are stable.
