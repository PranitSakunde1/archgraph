# ArchGraph by Persona

ArchGraph is not only a developer dependency map. It is a shared, evidence-backed view of how products are implemented across repositories, services, APIs, events, AWS resources and data stores.

Different roles need different answers from the same architecture graph. The graph is the shared source of evidence; reports, queries and workflows should be tailored to each persona.

> **Product boundary:** These are target capabilities and intended outcomes. They must not be presented as fully implemented until the corresponding discovery, graph, reporting and workflow capabilities exist and are validated against real repositories.

## Product Manager

### Why ArchGraph matters
Product Managers need to understand the technical surface area behind a product capability without having to reconstruct it from engineering conversations and scattered documentation.

### Questions ArchGraph should help answer
- Which services, APIs, events and data stores support a product capability?
- Which customer journeys could be affected by a proposed feature or change?
- What upstream or downstream dependencies could expand scope?
- Which parts of the capability have weak documentation or unresolved ownership?
- What technical dependencies or integration risks could affect a release?

### Useful outputs
- Business-capability-to-service map.
- Journey and dependency view.
- Change impact summary in product language.
- Known risks, ownership gaps and unresolved relationships.

### Value
Better scope conversations, earlier visibility of cross-team dependencies and fewer surprises during release planning. ArchGraph does not replace product discovery, customer research or delivery estimation.

## Business Analyst

### Why ArchGraph matters
Business Analysts translate business processes into system behavior. ArchGraph should connect process steps to the APIs, events, services and data transformations that implement them.

### Questions ArchGraph should help answer
- Which systems participate in a business process?
- Where does a business event originate and which services consume it?
- Which fields are created, transformed, persisted or propagated?
- Where do API contracts, event payloads and implementation behavior disagree?
- Which process steps have unresolved technical ownership or missing evidence?

### Useful outputs
- Business journey-to-system traceability.
- API and event flow diagrams.
- Field-level data lineage and transformation evidence.
- Contract-to-implementation mismatch reports.
- Explicit unresolved links requiring analyst or engineering confirmation.

### Value
More precise functional specifications, traceable requirements and faster impact assessment. ArchGraph cannot infer business meaning reliably from code alone; analysts and domain owners must validate business semantics.

## Tech Lead

### Why ArchGraph matters
Tech Leads need a trustworthy view of dependencies and technical boundaries across repositories, especially when a change crosses service, API, event, AWS or database boundaries.

### Questions ArchGraph should help answer
- What will this change affect directly and transitively?
- Which services publish or consume a particular event or API?
- Which internal wrappers hide AWS SDK operations?
- Which services access the same database tables, collections or fields?
- Where are dependency cycles, unresolved relationships or architecture drift?
- Which tests or business journeys should be run for a proposed change?

### Useful outputs
- Cross-repository dependency graph.
- Evidence-backed blast-radius analysis.
- API/event/resource relationship views.
- Database access and data-lineage maps.
- Architecture quality and coverage reports.
- Impact-based journey recommendations.

### Value
Faster design reviews, safer refactoring and better-informed technical decisions. Inferred relationships must remain visibly distinct from confirmed relationships.

## QA Engineer

### Why ArchGraph matters
QA Engineers need to know which components participate in a behavior and what can break when one component changes. A repository list or manually maintained test checklist is not enough for distributed workflows.

### Questions ArchGraph should help answer
- Which services and integrations participate in this journey?
- Which APIs, events and data stores are involved?
- What downstream consumers may be affected by a changed contract or event?
- Which fixtures, data fields and transformations need assertions?
- Which dependencies are unresolved and need test-environment investigation?

### Useful outputs
- Service and integration map for a test scenario.
- Journey dependency plan.
- API/event contract and data-lineage references.
- Suggested affected test areas, with evidence.
- Environment dependency inventory for local integration testing.

### Value
More focused test design and fewer missed integration paths. ArchGraph can recommend test scope, but QA engineers still define assertions, expected outcomes and test adequacy.

## QA Lead

### Why ArchGraph matters
QA Leads need a portfolio-level view of test coverage and risk across teams, not just individual test cases or pipeline status.

### Questions ArchGraph should help answer
- Which critical business journeys have approved automated coverage?
- Which services or integrations are not covered by those journeys?
- What changed in the dependency graph between releases?
- Which journeys should be prioritized for regression after a change?
- Where do unsupported languages, unresolved edges or stale contracts weaken confidence?
- Are test environments exercising the relevant dependency set?

### Useful outputs
- Journey coverage mapped to services, APIs, events and data stores.
- Change-based regression recommendations.
- Coverage gaps and unresolved dependency reports.
- Cross-team quality and architecture risk summaries.
- Trendable graph-quality and test-coverage measures.

### Value
More defensible regression selection and visibility into systemic coverage gaps. A graph relationship is not proof that a test covers the behavior; coverage must be linked to executed tests and verified outcomes.

## Developer

### Why ArchGraph matters
Developers need fast, concrete answers while changing code: what calls this, what does this call, which resource is used, and what else might break?

### Questions ArchGraph should help answer
- Which services, functions or consumers depend on this component?
- Where is an API or event produced and consumed?
- Which AWS resource does an internal wrapper resolve to?
- Which database collections, tables or fields are read or written?
- What is the likely blast radius of this change?
- Which relevant tests or approved journeys should run?

### Useful outputs
- Local service and code dependency graph.
- Source-linked relationship evidence.
- API/event/AWS/database navigation.
- Change impact and graph diff.
- Relevant test and journey recommendations.
- Clear unsupported or unresolved analysis warnings.

### Value
Less time reverse-engineering unfamiliar repositories, more targeted changes and fewer avoidable integration regressions. ArchGraph should explain why it created each relationship and point back to source evidence.

## Shared value and role-specific views

| Shared graph capability | Product Manager | Business Analyst | Tech Lead | QA Engineer | QA Lead | Developer |
|---|---|---|---|---|---|---|
| Service and domain inventory | Product scope | Process systems | Boundaries | Test components | Portfolio coverage | Find code owners |
| API and event relationships | Feature dependencies | Process transitions | Contract design | Integration paths | Regression scope | Trace producers/consumers |
| Data lineage | Data impact | Field meaning and mapping | Data architecture | Data assertions | Data risk | Trace reads/writes |
| Change impact | Release scope | Requirement impact | Blast radius | Test selection | Regression prioritization | Safe implementation |
| Journey mapping | Customer/business flow | End-to-end process | System design | Scenario design | Journey coverage | Local reproduction |
| Evidence and coverage | Delivery risk | Traceability gaps | Confidence in design | Unknown test paths | Coverage confidence | Source-level proof |

## Product principles for persona experiences

1. **One graph, multiple views.** Do not create six disconnected architecture models.
2. **Translate, do not distort.** Role-specific language may simplify presentation, but it must preserve technical meaning and uncertainty.
3. **Evidence is inspectable.** Every important relationship should link to source files, contracts, configuration or infrastructure declarations.
4. **Separate fact from inference.** Confirmed, inferred and unresolved relationships must be visible in every persona view.
5. **Show gaps honestly.** Unsupported languages, dynamic behavior and missing contracts must appear as coverage gaps, not disappear from reports.
6. **Do not overclaim test coverage.** Suggested journeys and impacted tests are recommendations until execution evidence confirms coverage.
7. **Keep business semantics human-approved.** Product capabilities, business journeys and domain ownership require validation by people who understand the business.

## Suggested delivery order

- **Foundation:** service inventory, ownership metadata, evidence-backed API/event/resource relationships.
- **Engineering utility:** developer and Tech Lead dependency navigation, graph diff and impact analysis.
- **Quality workflows:** QA Engineer test-scope recommendations, then QA Lead journey coverage and regression reporting.
- **Business visibility:** Business Analyst journey and field-lineage views, then Product Manager capability and release-impact summaries.

This order builds on the same graph and avoids shipping persona dashboards before the underlying relationships are trustworthy.
