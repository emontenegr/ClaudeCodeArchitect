# Anti-Patterns

## Unresolved Options

### Description

Presenting *multiple implementation choices* without deciding.

### Examples

| Category | ❌ Wrong | ✅ Correct |
|----------|---------|-----------|
| Database | "Storage: Database A or Database B" | "Storage: PostgreSQL 16" |
| Library | "Use Library X or Library Y for visualization" | "Use Library X v3.4.0 for visualization" |
| Deployment | "Deploy on Cloud A or Cloud B" | "Deploy on AWS (region us-east-1)" |
| Data Store | "Database: SQL or NoSQL" | "Database: PostgreSQL 16" |

### Why Bad

* AI must choose *arbitrarily*
* May pick *wrong option* for context
* Inconsistent choices across related decisions
* Planning responsibility shifted to AI

### Fix

Make decision during planning.

Include **only** chosen option.

Add rationale if non-obvious.

## Conditional Implementation

### Description

Implementation depends on runtime conditions not specified.

### Examples

| Category | ❌ Wrong | ✅ Correct |
|----------|---------|-----------|
| Caching | "If query is slow, add caching" | "Add caching with 5-minute TTL for queries" |
| Scaling | "When load exceeds threshold, implement sharding" | "Implement sharding on user_id (8 shards)" |
| Async | "Should this fail, fall back to synchronous" | "Use asynchronous processing with 10s timeout" |
| Rate Limiting | "Optional: Add rate limiting if needed" | "Implement rate limiting: 100 requests/minute" |

### Why Bad

* Conditional logic deferred to implementation
* AI cannot determine "if needed" threshold
* Incomplete feature set at launch

### Fix

Decide now: include or exclude.

Specify exact trigger thresholds if conditional.

Single implementation path only.

## Vague Requirements

### Description

*Underspecified* functionality or constraints.

### Examples

| Category | ❌ Wrong | ✅ Correct |
|----------|---------|-----------|
| Performance | "Should be fast" | "P99 latency <100ms" |
| Algorithm | "Use efficient algorithm" | "Use binary search O(log n)" |
| Error Handling | "Handle errors appropriately" | "Return 400 + error JSON on validation failure" |
| Scale | "Support large datasets" | "Support up to 1M records in memory (~64MB)" |
| Best Practices | "Follow best practices" | "Use language-standard error wrapping" |

### Why Bad

* AI interprets subjectively
* No measurable success criteria
* Impossible to validate correctness

### Fix

**Quantify** all requirements.

Specify **exact** behavior.

Include numbers, thresholds, concrete criteria.

## Missing Versions

### Description

Dependencies without version specifications.

### Examples

| Category | ❌ Wrong | ✅ Correct |
|----------|---------|-----------|
| Database | "Use database" | "Use PostgreSQL 16 (Docker: postgres:16-alpine)" |
| Framework | "Install framework and UI library" | "Install React 18.2.0 and Next.js 14.1.0" |
| Logging | "Add logging" | "Add winston@3.11.0 for structured logging" |
| Containers | "Deploy with containers" | "Deploy with Docker 24.0" |

### Why Bad

* Version drift causes build failures
* Breaking changes in dependencies
* Non-reproducible builds

### Fix

Pin all versions.

Use ecosystem's lock files.

Specify exact image tags for containers.

## Optional Sections

### Description

Features marked as optional or future work.

### Examples

| Category | ❌ Wrong | ✅ Correct |
|----------|---------|-----------|
| Features | "Optional: Add real-time updates" | Include feature in spec OR omit entirely |
| Future | "Future: Implement caching" | "Caching: Not included in current scope" |
| Nice-to-Have | "Nice to have: Admin dashboard" | Single feature set, fully specified |
| Consideration | "Consider adding: Rate limiting" | Include OR explicitly exclude |

### Why Bad

* Ambiguous scope
* AI may include or omit unpredictably
* Incomplete product definition

### Fix

Define exact feature set for this version.

Mark explicitly excluded features if helpful.

No "optional" - binary include/exclude.

## Weak Obligation Language

### Description

Using "should", "could", "might" instead of imperatives.

### Examples

| Category | ❌ Wrong | ✅ Correct |
|----------|---------|-----------|
| Error Handling | "Should implement error handling" | "Implement error handling (wrap with context)" |
| Connection Pool | "Could use connection pooling" | "Use connection pooling (max 25 connections)" |
| Rate Limiting | "Might need rate limiting" | "Implement rate limiting (100/min)" |
| Caching | "Probably cache results" | "Cache results (5min TTL)" |

### Why Bad

* AI interprets as optional
* Inconsistent implementation
* Unclear requirements

### Fix

Use imperative: "Do X"

Eliminate hedging.

Be decisive.

## Domain-Specific Examples

### Description

Examples that bias AI toward specific domain/technology.

### Examples

| Type | ❌ Wrong | ✅ Correct |
|--------|---------|-----------|
| Structure | Showing only web app structures | Abstract categorical descriptions |
| Patterns | Showing only ML pipeline patterns | Multiple diverse mini-examples |
| Products | Using specific product names as universal | Decision frameworks that generalize |
| Technical Details | Deep technical examples from one domain | "If your system has X, consider Y" |

### Why Bad

* LLMs fixate on concrete examples
* Cargo-cult implementation in wrong domain
* Inappropriate pattern transfer
* Loss of reasoning ability

### Fix

Teach reasoning, not templates.

Use categories, not specific artifacts.

Multiple diverse examples if examples needed.

Explicit "adapt to your context" guidance.
