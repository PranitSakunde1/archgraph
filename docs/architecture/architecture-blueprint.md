# ArchGraph Architecture Blueprint

**Status:** Consolidated target architecture  
**Implementation:** Go  
**Execution:** One-shot Docker container  
**Discovery:** Deterministic and AI-free  
**Language targets:** Go, Java, Python, React, TypeScript, JavaScript

## 1. Purpose

ArchGraph discovers application architecture across a workspace containing multiple Git repositories. It identifies services and Lambda functions from recursively discovered codeway.yaml files, then correlates source code, contracts, AWS resource definitions, dependencies and configuration into an evidence-backed graph.

The system must expose uncertainty. A graph edge is not automatically true because a naming convention or dependency suggests it.

## 2. Operating model

A single Docker invocation scans one workspace or an explicit set of repositories, writes artifacts to a mounted output directory and exits. Source mounts are read-only. Temporary extraction work happens inside the container. Generated artifacts go to /output. No daemon, API server or persistent container is part of v1.

Example PowerShell workflow:

    docker build --pull -t archgraph:dev .
    New-Item -ItemType Directory -Force .\archgraph-output | Out-Null
    docker run --rm -v "D:\work\repositories:/workspace:ro" -v "${PWD}\archgraph-output:/output" archgraph:dev scan --path /workspace --output /output

## 3. Architecture

    Workspace with many Git repositories
       -> codeway.yaml discovery and service identity
       -> Graphify code-only extraction
       -> language and framework adapters
       -> API, event, AWS and database discoverers
       -> Go normalization and stable identity
       -> evidence and confidence classification
       -> graph validation and coverage
       -> architecture.json, Markdown and Mermaid
       -> diff, impact analysis and journey selection

Graphify supplies local code-level graph extraction. Its raw artifacts remain separate until a validated normalization adapter maps them into the canonical graph. The Go core handles repository discovery, identity, enrichers, resource correlation, normalization, evidence, validation, diff, impact and journey orchestration. Neo4j is optional future infrastructure, not a v1 requirement.

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

React is a framework adapter, not a separate language. Graphify remains the default code graph; native adapters should add domain facts Graphify does not provide, not duplicate its work. Unsupported syntax and languages remain in inventory with unsupported or partial analysis status.

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

Version the canonical graph schema. Pin Graphify and validate its output schema. Unknown/incompatible output fails clearly rather than being silently reinterpreted.

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

Reuse Graphify MCP initially. Add ArchGraph-specific MCP tools only for normalized architecture queries, impact analysis, journey planning or reporting that Graphify does not provide.

## 14. Implementation phases

1. Reconnaissance: inspect real codeway files, Graphify output, internal AWS wrappers, SDK versions, IaC conventions and a representative business journey.
2. Local graph foundation: CLI, workspace discovery, identity, Graphify adapter, canonical schema, artifacts and validation.
3. Relationship intelligence: six language/framework targets, API/event/AWS/database discovery, wrapper mapping, lineage.
4. Impact and CI: snapshots, diffs, coverage baselines, hooks, PR reports and CODEOWNERS review.
5. Journey engine: journey schema, approval, minimal environment planning, event-driven runner and cleanup.
6. Security: scanner integrations, finding correlation, secret redaction and baseline gates.
7. Optional hybrid scale: assess Graphify MCP gaps and optional Neo4j backend.

## 15. Validation requirements before production claims

- Obtain representative real codeway.yaml files; do not invent schema.
- Pin and inspect actual Graphify output format/version.
- Collect internal AWS wrapper examples.
- Validate resource logical-name to code-reference mapping.
- Document static analysis limits for reflection, dynamic dispatch, generated queries and serializers.
- Measure baseline coverage before enforcing provisional thresholds.
- Verify scanner versions, output formats, licenses and supported ecosystems.
- Benchmark representative repository subsets before promising scan times.

The largest engineering risk is incorrect inference from undocumented conventions, wrapper abstractions and dynamic behavior. Make uncertainty visible, measure coverage and validate against repository fixtures before enforcing aggressive gates.
