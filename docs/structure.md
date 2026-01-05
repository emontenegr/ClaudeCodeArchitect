# Structural Requirements

## Plan File Naming

### Specification

Plan file ***MUST*** be named `MANIFEST.adoc` (uppercase).

| Status | ❌ NOT | ✅ IS |
|--------|--------|--------|
| Naming | • architecture.adoc<br>• manifest.adoc (lowercase)<br>• spec.adoc<br>• Architecture.adoc (mixed case) | • MANIFEST.adoc (uppercase) |

**Location:** `project/MANIFEST.adoc`

### Rationale

* **Consistent naming** across all components
* Immediately recognizable as *executable plan*
* Distinguished from documentation (`architecture.md`, `design.md`)
* Convention signals importance (like `README`, `LICENSE`)
* Entry point for modular specification system

## Modular Organization

### The Problem

Monolithic specification files hit fundamental limits:

* **Tool limits:** Read tool cannot process `1000+` line files efficiently
* **Maintainability:** Finding specific sections is difficult
* **Editability:** Small changes require reading entire file
* **Cross-cutting concerns:** Performance/security constraints scattered throughout

### Solution: File-Per-Concern

```
project/
├── MANIFEST.adoc              # Entry point, imports everything
├── core/
│   ├── metadata.adoc          # System context, paradigm, foundation
│   ├── types.adoc             # Complete type definitions
│   └── principles.adoc        # Architectural constraints, invariants
├── architecture/
│   ├── structure.adoc         # File tree, package organization
│   ├── data-flow.adoc         # How data moves through system
│   └── concurrency.adoc       # Threading/async model
├── interfaces/
│   ├── api.adoc               # API routes, params, responses
│   └── formats.adoc           # Data format examples
├── algorithms/
│   ├── search.adoc            # Search algorithm specification
│   ├── indexing.adoc          # Indexing algorithm specification
│   └── performance.adoc       # Performance requirements & derivations
├── concerns/                  # Cross-cutting concerns
│   ├── performance.adoc       # System-wide performance constraints
│   ├── security.adoc          # System-wide security requirements
│   └── observability.adoc     # Logging, metrics, monitoring
├── operations/
│   ├── deployment.adoc        # Platform, approach, configuration
│   ├── testing.adoc           # Test requirements, coverage
│   └── monitoring.adoc        # Observability specifications
├── references/                # Domain-specific artifacts
│   ├── sample-data.json       # Example data
│   └── schemas/               # External schemas if needed
└── (implementation files)     # AI implements here
```

### File Size Guidelines

Files should be readable by your AI tooling and navigable by humans.

**The constraints:**

* AI tools have read limits (varies by tool - know yours)
* Humans navigate better with focused files
* Large files are hard to search, edit, and reason about

**Decision criteria:**

* Can your AI read the entire file in one operation?
* Can you quickly find what you need?
* Does it address one cohesive concern?

**Split when:** File exceeds your tool's read limit OR becomes hard to navigate OR addresses multiple unrelated concerns

**Example:** If `algorithms/` becomes too large:

```
algorithms/
├── sorting.adoc
├── indexing.adoc
├── caching.adoc
└── performance.adoc    # Cross-cutting performance for algorithms
```

### MANIFEST.adoc Structure

The entry point that composes all parts:

```
# Project Name Specification

// Define reusable attributes (like constants)

// Import core definitions
include::core/metadata.adoc[]

include::core/types.adoc[]

include::core/principles.adoc[]

// Import architecture
include::architecture/structure.adoc[]

include::architecture/data-flow.adoc[]

include::architecture/concurrency.adoc[]

// Import interfaces
include::interfaces/api.adoc[]

include::interfaces/formats.adoc[]

// Import algorithms
include::algorithms/performance.adoc[]

// Import operations
include::operations/deployment.adoc[]

include::operations/testing.adoc[]
```

### Atomic Editability

Each file addresses ONE concern:

* Want to change API? → Edit `interfaces/api.adoc` only
* Want to update types? → Edit `core/types.adoc` only
* Want to fix deployment? → Edit `operations/deployment.adoc` only

No hunting through monolithic files.

Git diffs show exactly what changed:

```bash
$ git diff
diff --git a/concerns/performance.adoc b/concerns/performance.adoc
-:api-p99-latency: 100ms
+:api-p99-latency: 50ms
```

Single line change, system-wide effect.

## Cross-Cutting Concerns

### The Problem

Some requirements apply across many components:

* Performance constraints (all APIs must be <100ms)
* Security requirements (all endpoints need auth)
* Observability standards (all operations must log)

Inline duplication leads to:

* Inconsistency (requirements differ across files)
* Maintenance burden (update 10 places for one change)
* Missing coverage (easy to forget one place)

### Solution: Tagged Sections + Attributes

AsciiDoc provides native solutions:

**1. Attributes for Constants**

Define once in `MANIFEST.adoc` or `concerns/performance.adoc`:

```asciidoc

:api-p99-latency: 100ms
:api-p50-latency: 20ms
:circuit-breaker-threshold: 3
:circuit-breaker-window: 10s
:db-connection-pool: 25
```

Use everywhere:

```asciidoc
All endpoints must meet: P99 < {api-p99-latency}, P50 < {api-p50-latency}
```

**2. Tagged Sections for Reusable Blocks**

`concerns/performance.adoc`:

```asciidoc
# Performance Requirements

// tag::api-latency[]
**API Latency Requirements:**

* P99 latency: <{api-p99-latency}
* P50 latency: <{api-p50-latency}
* Circuit breaker: {circuit-breaker-threshold} failures in {circuit-breaker-window}
* Timeout: 30s
// end::api-latency[]

// tag::database-performance[]
**Database Query Performance:**

* P95 query time: <10ms
* Connection pool: {db-connection-pool} connections
* Query timeout: 5s
// end::database-performance[]
```

`interfaces/api.adoc` includes tagged section:

```asciidoc
# API Specification

## POST /users

include::../concerns/performance.adoc[tag=api-latency]

include::../concerns/security.adoc[tag=oauth-required]

**Endpoint Details:**

Route: POST /users
Request Body: { username, email, password }
...
```

**Result when compiled:**

```asciidoc
## POST /users

**API Latency Requirements:**

* P99 latency: <100ms
* P50 latency: <20ms
* Circuit breaker: 3 failures in 10s
* Timeout: 30s

**OAuth Required:**
...

**Endpoint Details:**
Route: POST /users
...
```

### Benefits

**Single source of truth:** Change `concerns/performance.adoc`, propagates everywhere

**Consistent:** Impossible to have conflicting requirements

**Discoverable:** All performance constraints in one file

**Auditable:** Can extract all components that include a concern

**Tool-friendly:** Standard AsciiDoc, no custom preprocessing

### When to Use Cross-Cutting Concerns

**Use `tag::` includes when:**

* Requirement applies to >3 components
* Requirement must stay consistent across system
* Requirement changes together (tighten all API latency)

**Use inline when:**

* Requirement specific to one component
* Requirement varies by component
* Requirement is part of component's unique contract

## Metadata Section

### Specification

Every `MANIFEST.adoc` begins with metadata describing system context.

**Required fields:**

* `system_name`: Specific name of this component
* `paradigm`: Fundamental approach/architecture pattern
* `foundation`: Core technologies/algorithms
* `language`: Exact language + version
* `dependencies`: Critical external services (if applicable)

Metadata should be SPECIFIC to your system. NOT generic. NOT vague.

**Example:**

```asciidoc
# User Authentication Service

[metadata]

*System Name:* User Authentication Service
*Paradigm:* OAuth 2.0 + JWT token-based authentication
*Foundation:* PostgreSQL 16 for user storage, Redis 7 for session cache
*Language:* Go 1.21
*Dependencies:*
- PostgreSQL 16
- Redis 7
- SMTP service (SendGrid API v3)

```

### Reasoning

Metadata provides AI with immediate context about:

* What kind of system is being built
* What constraints apply
* What paradigms to follow
* What dependencies exist

Enables AI to make appropriate inferences for the rest of spec.

### Adaptability Note

Metadata fields vary by domain.

| Domain | Typical Metadata |
|--------|------------------|
| Web App | framework, database, auth_provider |
| ML System | model_type, training_framework, inference_engine |
| Game Engine | render_pipeline, physics_engine, asset_format |
| IoT Firmware | microcontroller, RTOS, communication_protocol |
| Fintech | regulatory_framework, ledger_system, settlement_protocol |

Include fields relevant to YOUR system's architectural decisions.

## Type Definitions

### Specification

Include complete type definitions for core data structures.

**Language-specific guidance:**

**Go:**

* All struct fields with exact types (`uint32`, not "integer")
* Memory layout annotations
* Size calculations if relevant

**TypeScript:**

* All interface properties
* Exact types (`number`, `string`, `Date`, custom types)
* Optional vs required (`?` notation)

**Python:**

* Dataclass or TypedDict definitions
* Type hints for all fields

Include brief comments for non-obvious fields.

Specify invariants and constraints.

**Example (TypeScript):**

```asciidoc
# Type Definitions

## User

```typescript
interface User {
  id: string;              // UUID v4
  username: string;        // 3-32 alphanumeric characters
  email: string;           // RFC 5322 format, unique constraint
  passwordHash: string;    // bcrypt hash, cost=12
  createdAt: Date;         // UTC timestamp
  lastLoginAt: Date | null; // null if never logged in
  roles: Role[];           // At least one role required
}
```

**Invariants:**

* `id` must be globally unique UUID v4
* `username` must match regex: `^[a-zA-Z0-9]{3,32}$`
* `email` must be unique across all users
* `roles` array cannot be empty

### Reasoning

Complete type definitions enable AI to:

* Generate correct code immediately
* Understand memory layout and performance implications
* Make appropriate trade-offs
* Avoid ambiguous representations

Without complete types, AI must guess or ask questions.

### Domain Examples

| Domain | Core Types |
|--------|------------|
| Fintech | Transaction, Account, Ledger |
| Game | Entity, Transform, PhysicsBody |
| ML | Tensor, Layer, Activation |
| Web | User, Session, Request |

Define types central to YOUR domain's core abstractions.

## Algorithm Specifications

### Specification

Specify algorithms with precise steps, not pseudocode.

**Include:**

* Mathematical formulation (if applicable)
* Step-by-step procedure
* Edge case handling
* Performance characteristics (Big-O, actual timing estimates)
* Library to use (if standard implementation exists)

Use actual language syntax where helpful for clarity.

**Example:**

```asciidoc
# Search Algorithm

## User Search by Email Prefix

**Algorithm:** Binary search on sorted email index

**Input:** email_prefix (string, 1-64 chars)

**Output:** List<User>, max 20 results

**Procedure:**

1. Normalize input: `email_prefix = toLowerCase(trim(email_prefix))`
2. Validate: If length < 1 or > 64, return empty list
3. Query index: `SELECT * FROM user_email_index WHERE email >= $1 AND email < $2 LIMIT 20`
   - $1 = email_prefix
   - $2 = email_prefix + '~' (next lexicographic string)
4. Fetch user records for matching IDs
5. Return list, empty if no matches

**Performance:**

* Index: B-tree on `users.email` column
* Complexity: O(log n) for index lookup + O(20) for fetch = O(log n)
* Expected latency: P95 < 5ms for 1M user database

**Edge Cases:**

* Empty prefix → return empty list
* Prefix too long → return empty list
* No matches → return empty list
* Special characters in prefix → exact match only (no normalization)

**Library:** Use database native B-tree index, no custom implementation needed
```

### Reasoning

Precise algorithm specs prevent:

* AI choosing inefficient approach
* Correctness bugs from wrong algorithm
* Performance regressions

When algorithm choice is critical, be explicit.

When standard approach exists, reference library.

### Decision Framework

Ask: Is there a standard algorithm/library for this?

| Answer | Action |
|--------|--------|
| YES, standard exists | Reference library, specify version |
| NO, custom needed | Specify algorithm precisely |
| MAYBE, domain-specific | Specify algorithm if non-obvious |

## Data Format Specifications

### Specification

All data formats include complete examples with real values.

**For JSON/JSONL:**

* Full example objects (not schema, actual data)
* All fields with representative values
* Type indicators if helpful
* Token/size estimates if relevant for performance

**For APIs:**

* Complete request/response pairs
* All headers
* Status codes (success and error cases)
* Error response formats

**For binary/wire protocols:**

* Byte layout diagrams
* Endianness specifications
* Field offsets and sizes

**Example:**

```asciidoc
# API Data Formats

## POST /users - Create User

**Request:**

```http
POST /users HTTP/1.1
Host: api.example.com
Content-Type: application/json
Authorization: Bearer eyJhbGc...

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response 201 Created:**

```http
HTTP/1.1 201 Created
Content-Type: application/json
Location: /users/550e8400-e29b-41d4-a716-446655440000

{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "john_doe",
  "email": "john@example.com",
  "createdAt": "2024-01-15T10:30:00Z",
  "roles": ["user"]
}
```

**Response 400 Bad Request:**

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "error": "validation_failed",
  "details": [
    {
      "field": "email",
      "message": "Invalid email format"
    }
  ]
}
```

**Response 409 Conflict:**

```http
HTTP/1.1 409 Conflict
Content-Type: application/json

{
  "error": "email_already_exists",
  "message": "User with this email already exists"
}
```

### Reasoning

Complete format examples prevent:

* Misinterpretation of structure
* Missing fields
* Wrong types
* Incompatible formats between components

Show, don't describe.

## Performance Specifications

### Specification

Performance requirements MUST include concrete numbers.

| Type | ❌ NOT | ✅ IS |
|--------|--------|--------|
| Latency | "Should be fast" | "P99 latency <100ms" |
| Throughput | "Low latency required" | "Support 1000 requests/second" |
| Processing | "Efficiently process" | "Process 10K items in <1ms" |
| Scale | "Handle large scale" | "Support 1M concurrent users" |

**Include:**

* Latency targets (P50, P95, P99 percentiles)
* Throughput requirements (requests/sec, items/sec)
* Memory constraints (max RAM, per-item size)
* Computational complexity (Big-O with practical timing)

### Reasoning

Quantified performance enables AI to:

* Choose appropriate algorithms
* Make data structure trade-offs
* Add caching/optimization where needed
* Validate implementation meets requirements

Without numbers, "fast" is subjective and unmeasurable.

## File Structure Specifications

### Specification

Use local reference frame principle.

Each project's spec specifies ONLY its project's local file structure.

NOT the entire repository, NOT sibling projects, NOT parent structure.

**Include:**

* Project's root directory
* All subdirectories and packages
* Key files with purpose annotations
* Internal organization

**Omit:**

* Repository root structure
* Other projects
* Parent/sibling directories

Like general relativity: specify local coordinate system, not global universe.

Each project's spec specifies ONLY its project's local file structure.

### Generic Structure Template

```
myproject/
├── MANIFEST.adoc              # Main specification
├── core/
├── concerns/
├── interfaces/
├── references/
└── (implementation files)
    ├── main.{ext}             # Entry point
    ├── package_a/
    ├── package_b/
    └── package_c/
```

Adapt structure to YOUR project's needs.

Names, packages, organization reflect YOUR domain.

### Domain Variations

**Web API service:**

```
myproject/
├── MANIFEST.adoc
├── core/
├── concerns/
├── handlers/
├── models/
├── db/
└── middleware/
```

**Game engine:**

```
myproject/
├── MANIFEST.adoc
├── core/
├── concerns/
├── renderer/
├── physics/
├── entities/
└── assets/
```

**ML training pipeline:**

```
myproject/
├── MANIFEST.adoc
├── core/
├── concerns/
├── data_loaders/
├── models/
├── training/
└── evaluation/
```

Structure emerges from YOUR system's architecture. Not prescribed by template.

## Dependency Specifications

### Specification

All external dependencies MUST specify exact versions.

Format varies by ecosystem:

| Ecosystem | Format |
|-----------|--------|
| Go | `github.com/user/repo@v1.2.3` |
| npm | `package@18.2.0` |
| Python | `package==1.24.0` |
| Docker | `postgres:16-alpine` |
| Cloud Services | Provider + tier/version |

Version pinning prevents drift and ensures reproducibility.

### Reasoning

Unversioned dependencies cause:

* Build failures (breaking changes)
* Non-reproducible builds
* Deployment inconsistencies
* Debug nightmares

Pin everything. Use lock files.
