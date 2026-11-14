# ArchGraph Architecture Blueprint

**Status:** Consolidated target architecture  
**Implementation:** Go  
**Static execution:** Native Go CLI (optional one-shot Docker packaging)  
**Runtime observation:** Native local host process  
**Discovery:** Deterministic and AI-free  
**Language targets:** Go, Java, Python, React, TypeScript, JavaScript

## 1. Purpose

ArchGraph discovers application architecture across a workspace containing multiple Git repositories. It identifies services and Lambda functions from recursively discovered codeway.yaml files, then correlates source code, contracts, AWS resource definitions, dependencies and configuration into an evidence-backed graph.

The system must expose uncertainty. A graph edge is not automatically true because a naming convention or dependency suggests it.

## 2. Operating model

Static analysis is a one-shot native Go CLI workflow. Docker remains an optional packaging/execution mode for reproducible isolated scans; it is not the required deployment model. Source scanning is read-only by default and writes generated artifacts only to the configured output directory. No always-on server is required for static analysis.

Runtime observation is a separate native local process running on the Docker host, including WSL/Linux. It watches local container lifecycle/topology and receives optional application telemetry through a local OTLP endpoint. It can run while service containers are active and updates local runtime evidence in the canonical graph. It is not deployed as a container and does not require exposing the Docker socket to an observer container.

All runtime data remains local. No SaaS observability backend, cloud service, LLM or external telemetry destination is required. Runtime observation can be stopped without affecting static scans or exports.

Optional Docker scan workflow (PowerShell):

    docker build --pull -t archgraph:dev .
    New-Item -ItemType Directory -Force .\archgraph-output | Out-Null
    docker run --rm -v "D:\work\repositories:/workspace:ro" -v "${PWD}\archgraph-output:/output" archgraph:dev scan --path /workspace --output /output

## 3. Architecture

    Workspace with many Git repositories
       -> codeway.yaml discovery and service identity
       -> native language and framework extractors
       -> API, event, AWS and database discoverers
       -> Go normalization and stable identity
       -> evidence and confidence classification
       -> graph validation and coverage
       -> architecture.json, Markdown and Mermaid
       -> diff, impact analysis and journey selection

ArchGraph owns its extraction and graph-construction pipeline. Native Go-based language and framework adapters produce structured architecture facts, which the Go core normalizes, correlates and validates before adding them to the canonical graph. The core handles repository discovery, identity, resource correlation, evidence, validation, diff, impact and journey orchestration. Neo4j is optional future infrastructure, not a v1 requirement.

## 4. Repository discovery and identity

The root command accepts a workspace path and recursively finds Git repositories. It also supports explicit repository paths.

Every codeway.yaml / codeway.yml is discovered recursively. The file is authoritative for service/Lambda identity. Actual field names and nesting must be confirmed from representative files before finalizing the parser. Duplicate identities, conflicting declarations and source-to-service association ambiguity are reported.

Domain and team ownership may be explicit metadata. Folder names and naming patterns can generate suggestions but cannot silently establish ownership.

## 5. Language and framework discovery

First-class targets:

- Go: packages, modules, handlers, routes, clients, dependency injection, database access, AWS SDK operations and event handlers.
- Java: Maven/Gradle modules, Spring Boot controllers, clients, beans, repositories, entities/DTOs, messaging and AWS SDK calls.
- Python: packages, FastAPI/Flask/Django routes, clients, ORM/document DB operations, boto3 and event handlers.
- TypeScript and JavaScript: modules, routes, clients, event handlers, database access and AWS SDK v3.
- React framework adapter: components, hooks, routes, state/data-fetching layers and frontend-to-backend API dependencies.

React is a framework adapter, not a separate language. Native language and framework adapters are responsible for extracting code structure and domain-specific architecture facts. Shared extraction interfaces and the canonical schema should prevent duplicated work across adapters. Unsupported syntax and languages remain in inventory with unsupported or partial analysis status.

## 6. AWS and event intelligence

Event discovery is code-first. Infrastructure and contract files enrich or validate code evidence.

Discover AWS resources and operations from application source and SDK call sites, internal wrapper implementations, Terraform, CloudFormation, SAM, dependency manifests, lockfiles, application configuration and resource references.

Dependency presence alone is not proof of usage. Trace the chain:

    application call -> internal wrapper -> SDK operation -> resource reference

Coverage includes SNS, SQS, EventBridge, Lambda, Kinesis, DynamoDB Streams, Amazon MQ, Step Functions, API Gateway and other identified AWS services.

Go and Java AWS SDKs, Python boto3 and Node.js AWS SDK v3 are initial first-class SDK targets. Repository-specific wrappers may use declarative YAML mappings for simple patterns or trusted Go adapters for complex patterns under .archgraph/plugins/. Do not execute arbitrary untrusted plugin code.

No live AWS account calls or credentials are required in v1. Deployment topology (accounts, regions, clusters and environment orchestration) is out of scope.

## 7. API, database and data lineage

API discovery combines OpenAPI/contracts, source handlers, route registrations, clients and configuration. Compare declared contracts with implementation and consumers. Breaking changes or critical consumer incompatibility may block; missing/stale specs and non-breaking mismatches warn.

Database discovery aims at schema and access level:
- MongoDB: databases, collections, fields, inferred types, indexes, reads/writes.
- SQL: schemas, tables, columns, types, keys, indexes and queries.
- DynamoDB: tables, keys, indexes, attributes and access patterns.
- Redis: key patterns, recognizable structures and operations.
- Elasticsearch/OpenSearch: indices, mappings, fields and queries.

Shared database access is an observed fact, not automatically a policy violation. Do not impose a single-writer rule.

Represent both field-level lineage edges and explicit transformation nodes. Trace nested structures, DTO/entity mapping, serialization/deserialization and cross-service propagation when statically identifiable. Dynamic unresolved paths remain unresolved; do not invent field mappings.

Sensitive-data classifications are built-in plus repository extensions. Preserve declared/confirmed and heuristic classifications separately. Track lineage without storing actual values.

## 8. Evidence and canonical graph

Core nodes include Repository, Domain, Team, Service, LambdaFunction, Module, API, APIEndpoint, APIClient, Event, EventProducer, EventConsumer, AWSResource, Database, Table, Collection, Field, Transformation, SecurityFinding, Dependency, Journey and Evidence.

Core edges include CONTAINS, OWNED_BY, CALLS, EXPOSES, CONSUMES_API, PRODUCES_EVENT, CONSUMES_EVENT, PUBLISHES_TO, DELIVERS_TO, READS_FROM, WRITES_TO, MAPS_TO, TRANSFORMS, SERIALIZES, DESERIALIZES, DEPENDS_ON, USES_RESOURCE, AFFECTS, PARTICIPATES_IN and HAS_FINDING.

Each relationship preserves:
- status: confirmed, inferred or unresolved
- deterministic confidence semantics
- evidence type and source location
- discoverer and adapter version
- resolution reference, if manually approved

Stable IDs are based on logical repository/service identity, entity kind and logical name/namespace, not line number or path alone. Source locations are evidence attributes.

Manual scans may prompt for important ambiguous relationships. Hooks and CI are non-interactive and use approved version-controlled overrides. Stale overrides are reopened when evidence changes.

## 9. Security intelligence

Integrate dedicated scanners such as Trivy, OSV and Checkov through built-in adapters, custom Go adapters and YAML settings. ArchGraph correlates findings to dependencies, services, APIs, events, resources, data flows and affected journeys; it does not replace scanners or claim static correlation proves exploitability.

Secret detection uses deterministic patterns plus context analysis:
- CONFIRMED_SECRET: block.
- POTENTIAL_SECRET: report for review.
- NON_SECRET: no finding.

Redact secret values from terminal output, logs, graph artifacts and CI reports. Baseline-aware security gates block newly introduced critical/high findings while tracking existing backlog. Baseline changes require reviewed approval and an audit trail.

## 10. Graph artifacts, snapshots and quality

Commit lightweight reproducible artifacts: architecture.json, architecture.md, architecture.mmd, graph-diff.json and graph-diff.md.

Detailed extraction caches, temporary files, debug logs and transient test artifacts remain local and ignored by default.

Snapshots associate with Git SHAs. Keep immutable approved release snapshots and support comparisons across revisions. CI generates reports and validates committed artifacts; developers commit artifact changes. CI must not push commits or rewrite branches.

Version the canonical graph schema and validate every extractor's output against explicit adapter contracts. Unsupported or incompatible data must fail clearly rather than being silently reinterpreted.

Quality dimensions include structural validity, service/API/event/resource coverage, relationship evidence, unresolved critical edges, unsupported/partial analysis, staleness and security coverage. Initial thresholds are provisional and must be calibrated on real repository data. Existing gaps are baseline-aware; newly introduced regressions block. CODEOWNERS reviews baseline changes.

## 11. Git and CI workflow

- Manual scan during development.
- Incremental pre-commit scan.
- Full consistency pre-push scan.
- CI validation on pull requests.

Hooks are bypassable; CI is authoritative. CI compares target and proposed snapshots, computes blast radius, recommends the smallest relevant approved journey set, publishes JSON/Markdown reports and verifies committed artifacts match generated output.

## 12. Business journey testing

Engineers author critical journeys. ArchGraph may suggest candidates, but a human must approve them before execution.

Plan the smallest dependency-aware Docker environment for the journey. Do not include every consumer of a shared topic automatically. Use Docker Compose, Testcontainers and LocalStack where appropriate.

Async tests use correlation IDs, event acknowledgements, condition-based waits and bounded timeouts. No arbitrary sleeps. Assert final business outcomes, not just message delivery.

Use versioned baseline fixtures, unique per-run values/namespaces and deterministic cleanup. Retain failed-run data optionally. Parallelize independent journeys; serialize journeys that share mutable resources.

## 13. CLI and MCP

Initial CLI families:

    archgraph scan
    archgraph diff
    archgraph validate
    archgraph query
    archgraph impact
    archgraph journey validate
    archgraph journey plan
    archgraph journey run
    archgraph report

Add ArchGraph-specific MCP tools when the normalized architecture graph supports reliable queries, impact analysis, journey planning or reporting.

## 14. Implementation phases

1. Reconnaissance: inspect real codeway files, internal AWS wrappers, SDK versions, IaC conventions and a representative business journey.
2. Local graph foundation: CLI, workspace discovery, identity, native extraction adapters, canonical schema, artifacts and validation.
3. Relationship intelligence: six language/framework targets, API/event/AWS/database discovery, wrapper mapping, lineage.
4. Impact and CI: snapshots, diffs, coverage baselines, hooks, PR reports and CODEOWNERS review.
5. Journey engine: journey schema, approval, minimal environment planning, event-driven runner and cleanup.
6. Security: scanner integrations, finding correlation, secret redaction and baseline gates.
7. Optional hybrid scale: assess ArchGraph query and MCP requirements, plus an optional Neo4j backend.

## 15. Validation requirements before production claims

- Obtain representative real codeway.yaml files; do not invent schema.
- Collect internal AWS wrapper examples.
- Validate resource logical-name to code-reference mapping.
- Document static analysis limits for reflection, dynamic dispatch, generated queries and serializers.
- Measure baseline coverage before enforcing provisional thresholds.
- Verify scanner versions, output formats, licenses and supported ecosystems.
- Benchmark representative repository subsets before promising scan times.

The largest engineering risk is incorrect inference from undocumented conventions, wrapper abstractions and dynamic behavior. Make uncertainty visible, measure coverage and validate against repository fixtures before enforcing aggressive gates.

 
## 16. Local Docker runtime observer

### 16.1 Deployment boundary

The observer is a native local process on the Docker host, not a Docker container. It supports the developer's local Docker environment, including Docker Desktop with WSL2/Linux workflows. It is separate from the one-shot static scan command.

The observer uses the local Docker API/CLI integration to discover containers, images, labels, networks, port mappings and lifecycle events. Access to the Docker daemon is powerful and must be treated as privileged host access. Keep the observer process under the user's local account, do not expose the Docker API over an unauthenticated TCP socket, and do not mount the Docker socket into an observer container.

### 16.2 Two evidence channels

**Docker topology channel**
- Container create/start/stop/die/restart events.
- Image identity and digest where available.
- Compose project/service labels and user-defined labels.
- Network membership, aliases, exposed/published ports and declared environment metadata, with secret values redacted.
- Container-to-container network reachability candidates.

This channel describes deployment/container topology. It does not prove that one application called another.

**Application telemetry channel**
- Local OTLP receiver/export path for OpenTelemetry spans.
- HTTP and RPC client/server spans for observed request paths.
- Instrumented database client operations, represented without raw parameter values.
- Instrumented messaging producer/consumer spans and correlation context.
- Trace/span relationships and service identity mapping.

Actual application relationships are runtime-observed only when supported by telemetry evidence. Services without instrumentation remain marked as runtime coverage gaps. Network connections alone must not be upgraded to semantic API or database relationships.

### 16.3 Local-only processing and privacy

- No telemetry is sent to external services by default.
- No cloud account, hosted observability service or LLM is required.
- Do not capture or persist request/response bodies, authorization headers, credentials, cookies, raw SQL parameter values or message payloads by default.
- Redact sensitive span attributes before persistence and export.
- Make retention windows, storage paths and explicit local deletion configurable.
- Keep raw traces separate from committed architecture artifacts; commit only sanitized, aggregated graph evidence when repository policy allows it.
- Runtime snapshots record observation time/window and environment so short-lived observations are not mistaken for permanent architecture declarations.

### 16.4 Graph reconciliation

Keep these facts distinguishable:
- static source evidence;
- declared configuration and Compose topology;
- Docker container/network topology;
- runtime-observed application interactions;
- inferred or unresolved relationships.

A runtime observation enriches the canonical graph; it does not overwrite static evidence. An absent span is not proof of an absent edge. Runtime evidence should be linked to service/container identity using explicit OpenTelemetry resource attributes and Compose labels where available. Ambiguous identity mappings remain unresolved for review.

### 16.5 Lifecycle and CLI shape

The intended command family is:

    archgraph runtime start
    archgraph runtime status
    archgraph runtime snapshot
    archgraph runtime stop

The exact process supervision and IPC contract will be specified during implementation. Static commands remain independently usable:

    archgraph scan
    archgraph validate
    archgraph export

The observer can run while application containers are active. It can produce periodic local snapshots or update the local graph incrementally. Stopping the observer must not stop, restart or modify application containers.

### 16.6 Acceptance criteria

- Observer runs as a native host process and does not run as a privileged observer container.
- Docker start/stop and topology facts are discovered locally.
- Instrumented service calls appear as runtime-observed edges with trace provenance.
- Uninstrumented services are reported as runtime coverage gaps.
- No telemetry leaves the machine by default.
- Sensitive payloads and credentials are not persisted.
- Runtime observations never silently become static-confirmed edges.
- Static scanning and export work with the observer stopped.
