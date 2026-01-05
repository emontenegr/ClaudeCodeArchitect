# Decision Frameworks

## Explicitness Framework

### Principle

**Balance:** Be *explicit* when it matters, let AI decide when it doesn't.

### Be Explicit When

Implementation choice affects:

* Performance (latency, throughput, memory)
* Correctness (algorithm determines results)
* Architecture (concurrency, state management, data flow)
* Interoperability (protocols, formats, interfaces)

### Let AI Decide When

Choice is:

* Stylistic (naming, formatting, organization)
* Standard practice (error handling, logging)
* Obvious from context (helper functions, utilities)

### Decision Heuristic

Ask: "If AI chooses differently, will it break something?"

**YES (breaks) → Be explicit:**

* Performance targets missed
* Algorithm produces wrong results
* Architectural constraint violated
* Integration fails

**NO (just different) → Let AI decide:**

* Code still correct
* Meets all requirements
* Stylistic preference only

### Example Reasoning

**Lookup operation:**

* If called millions of times → **Specify:** `const table, O(1), no allocation`
* If called occasionally → **Let AI decide:** `map, array, whatever works`

**Error handling:**

* If error recovery critical → **Specify:** `retry logic, fallback, timeout`
* If standard error propagation → **Let AI decide:** knows language idioms

**Data structure:**

* If cache locality critical → **Specify:** `contiguous array, not map`
* If just storing data → **Let AI decide:** appropriate for access pattern

## Project Organization Framework

### Principle

Organize repositories by *project*, where each project contains its specification.

### Structure Pattern

```
repository/
├── project_1/
│   ├── MANIFEST.adoc
│   ├── core/
│   ├── concerns/
│   ├── references/
│   └── (implementation files)
├── project_2/
│   ├── MANIFEST.adoc
│   ├── core/
│   ├── concerns/
│   ├── references/
│   └── (implementation files)
└── (non-code directories)
```

Each project contains its specification directly.

### Project Identification

What constitutes a "project" varies by domain:

| Architecture | Projects |
|--------------|----------|
| Microservices | Each service is a project |
| Monolith | Separate by architectural layer or bounded context |
| Game | Engine, gameplay, ui, tools as projects |
| ML | Training, inference, data pipeline as projects |

Decompose YOUR system into architectural units.

Each unit gets `plan/` directory.
