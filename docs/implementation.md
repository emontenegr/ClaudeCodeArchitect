# Implementation Guidance

## What to Specify

### Architecture

* System decomposition (*components*, *layers*, *services*)
* Communication patterns (`sync/async`, *protocols*, *message formats*)
* Data flow (`input → processing → output`, *pipelines*)
* State management (*where state lives*, *how persisted*, *lifecycle*)
* Concurrency model (`threads/goroutines/async`, *synchronization*)

### Data Structures

* Complete type definitions (all fields)
* Memory layout and size (if performance-critical)
* Invariants and constraints
* Lifetime and ownership (who creates, who destroys)

### Algorithms

* Step-by-step procedure
* Mathematical formulation (if applicable)
* Complexity analysis (Big-O time and space)
* Edge cases and error conditions

### Interfaces

* Complete API specifications (all routes, all parameters)
* Data format examples (actual data, not just schema)
* Error responses (all error cases)
* Rate limits, quotas, constraints

### Operations

* Deployment configuration (specific platform, specific approach)
* Monitoring and observability (what to log, what to metric)
* Backup and recovery (strategy, frequency, restoration procedure)
* Scaling strategy (how system scales, at what thresholds)

## What Not to Specify

### Omit

* Obvious implementation details AI knows
* Standard error handling patterns (unless domain-specific)
* Common idioms in target language
* Generic best practices without specifics
* Actual implementation code (AI generates this)

### Boundary

Specify WHAT and WHY.

Specify HOW only when:

* Non-obvious approach
* Performance-critical
* Architectural constraint
* Multiple valid approaches exist and choice matters

Otherwise let AI choose HOW.

# Meta-Validation

## Self-Exemplification

This specification follows its own requirements:

**Zero unresolved options** (format specified, no alternatives)

**Zero conditionals** (no "if spec is complex, use...")

**Complete structure** (all sections defined)

**Quantified targets** (`90-95%` completion rate)

**No weak language** (uses "must", not "should")

**Decision frameworks** (not prescriptive templates)

**Anti-patterns with reasoning** (demonstrates principles)

**Domain-agnostic** (applicable to any project type)

***This document IS the pattern it teaches.***

## Cross-Domain Validation

### Test

Can this guide be used to write specs for:

* Fintech trading platform?
* Real-time multiplayer game?
* ML training pipeline?
* IoT device firmware?
* Blockchain validator?
* E-commerce backend?

If NO to any: Guide is too domain-specific, needs generalization.

If YES to all: Guide successfully domain-agnostic.

### Note

This guide teaches HOW TO REASON about specifications.

Not WHAT ARTIFACTS to create.

Reasoning transfers across domains.

Artifacts do not.

## Usage

To use this specification system:

1. Understand your domain and requirements (planning phase)
2. Make all architectural decisions for YOUR system
3. Resolve all options and conditionals
4. Create modular structure using AsciiDoc
5. Apply decision frameworks from this guide
6. Use `include::` for composition, `tag::` for cross-cutting concerns
7. Validate against completeness checklist
8. Test: Can unfamiliar person implement without questions?
9. Compile: `asciidoctor -b markdown MANIFEST.adoc > compiled-spec.md`
10. Submit to Claude Code for implementation
11. Expect 90-95% completion without iteration

This process works for ANY domain.

Frameworks are universal.

Your application determines specifics.
