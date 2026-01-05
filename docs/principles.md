# Foundational Principles

## Decided, Not Deciding

**Priority:** *Absolute*

Specifications are the ***output*** of planning, not the ***input*** to planning.

| Status | ❌ NOT | ✅ IS |
|--------|--------|--------|
| Undecided | • "Consider using X or Y"<br>• "This could be implemented with multiple approaches"<br>• "Optional: Add feature Z" | • "Use X v2.1.0"<br>• "Implement with approach Y"<br>• "Include feature Z" (or omit entirely) |
| Conditional | "If performance is insufficient, use caching" | "Use caching with 5-minute TTL" |

Every decision made BEFORE writing spec.

AI receives decisions, not choices.

This principle applies to ALL domains: fintech, games, ML systems, web apps.

## Zero Conditionals

**Priority:** *Absolute*

Eliminate all *conditional logic* from specifications.

| Type | ❌ NOT | ✅ IS |
|--------|--------|--------|
| Performance | "If query is slow, add index" | "Create index on (user_id, created_at)" |
| Error Handling | "Should this fail, fall back to..." | Single implementation path specified |
| Scaling | "When users exceed 10K, consider sharding" | "Shard by user_id hash when count > 1000" |

Conditionals indicate incomplete planning.

Resolve all conditions before spec creation.

Applies universally across all problem domains.

## Complete Specifications

**Priority:** *Absolute*

Specifications must be ***complete*** enough for one-shot implementation.

**Include:**

* Exact library versions with pinned numbers
* Complete type definitions (all fields, all types)
* Full file tree (every directory, every file in project)
* Database schemas (tables, columns, indexes, constraints)
* API specifications (routes, parameters, responses, status codes)
* Performance requirements (quantified with numbers)

**Omit:**

* Implementation code (AI generates this)
* Obvious standard practices (AI knows these)
* Generic advice ("follow best practices")

Specification completeness determines implementation success rate.

## Mathematical Derivation

**Priority:** *Critical*

All numeric constants must have clear *derivation* or *rationale*.

| Constant | ❌ NOT | ✅ IS |
|----------|--------|--------|
| Buffer Size | "Buffer size: 256" | "Buffer size: 256 = 2^8 (power of 2 for efficient allocation)" |
| Timeout | "Timeout: 30 seconds" | "Timeout: 30s = 3× p99 latency (10s observed)" |
| Max Items | "Max items: 1024" | "Max items: 1024 = 2^10 (cache line efficiency on target hardware)" |

Derived constants prevent arbitrary choices.

Enables AI to adjust coherently when needed.

Applies to any system with numeric parameters.

## Disambiguation Pattern

**Priority:** *Critical*

Use **NOT/IS pattern** to eliminate ambiguity.

**Format:**

```
NOT:
- [What this is NOT]
- [Common misconceptions]
- [Rejected approaches]

IS:
- [What this IS]
- [Actual implementation]
- [Precise definition]
```

Disambiguation prevents AI from making wrong inferences.

Universal pattern applicable to any specification.

## Local Reference Frame

**Priority:** *Critical*

Each project's `MANIFEST.adoc` specifies only its own *local structure*.

Like local reference frames in general relativity:

* Specify what's local, not the entire universe
* Project doesn't need to know about siblings
* Parent repository structure irrelevant

**Include in project spec:**

* This project's directory structure
* Internal packages/modules
* Project-specific dependencies

**Omit from project spec:**

* Repository-level structure
* Sibling projects
* Global concerns

Enables project independence and focused specifications.
