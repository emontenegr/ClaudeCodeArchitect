---
name: adoc
description: Handle AsciiDoc (.adoc) files. Use when you see .adoc files, MANIFEST.adoc, or spec/ folders. Never read .adoc files directly - use cca to compile them first.
---

# Architecture Specification Workflow

This project uses `cca` to manage architecture specifications. The spec is the source of truth - implement from it, don't guess.

## Understanding Specs: Context + Implementation

CCA specs have two layers:

**Context** - Intent and foundation (what/why/scope):
- Identity: What this system is
- Stack: Technical foundation (language, dependencies, versions)
- Problem: What gap this addresses (optional)
- Approach: How problem is solved conceptually (optional)
- Scope: What's in, what's explicitly out (optional)

**Implementation** - Complete technical details (types, schemas, APIs, performance)

Both are necessary. Context alone is vague. Implementation alone is mechanically complete but semantically ambiguous. Together they enable one-shot AI implementation.

## Reading the Spec

Before implementing, read the compiled spec:

```bash
cca compile
```

This outputs the full specification with all attributes resolved (e.g., `{api-p99-latency}` becomes `100ms`).

For specific sections:

```bash
cca compile --section "API Endpoints"
cca compile --section "Database Schema"
cca compile --section core/types.adoc
```

## Writing Specs (When Assisting Spec Authors)

When helping a user write a specification:

**1. Start with Context Section**

Every `MANIFEST.adoc` begins with `== Context`. Guide the user to include:
- **Identity** (required): Name, paradigm
- **Stack** (required): Language, dependencies with exact versions
- **Problem** (valuable for non-trivial systems): What gap this addresses
- **Approach** (valuable): How problem is solved conceptually
- **Scope** (valuable): Explicit in/out boundaries

Example minimal Context:
```asciidoc
== Context

=== Identity
*Name:* [System Name]
*Paradigm:* [Core architectural pattern]

=== Stack
*Language:* [Language + version]
*Dependencies:* [Exact versions]
```

**2. Recognize Completeness Includes Intent**

A spec is incomplete if an implementing AI would need to guess:
- System purpose
- Trade-off priorities
- Scope boundaries

Validation checks for "Ambiguous Intent." If Problem/Approach/Scope are absent and system is non-trivial, validation will flag this.

**3. Run Validation During Writing**

```bash
cca validate --quick   # Fast structural checks
cca validate           # Full validation (checks intent clarity)
```

Fix validation failures before spec is done.

## Implementing Specs (When Building From Specification)

**1. Read Context First**

Before writing code, understand intent:

```bash
cca compile --section Context
```

The Context section tells you:
- What this system is for (Problem)
- How to think about the solution (Approach)
- What's explicitly excluded (Scope)

Use this when encountering edge cases not explicitly covered.

**2. Full Implementation Workflow**

```bash
cca compile                   # Read full spec
cca compile --section "API"   # Focus on specific concerns
```

## Checking Attribute Impact

Before changing a spec attribute value, check what's affected:

```bash
cca impact <attribute-name>
```

Example:
```bash
cca impact api-p99-latency
```

Shows all sections using that attribute so you understand the change scope.

**3. Implement Exactly As Specified**

The spec is complete - don't add unspecified features. If something seems missing, it's either:
- In the spec (search more carefully)
- Explicitly excluded (check Scope section)
- Actually missing (validation should have caught this - ask user)

**4. Use Context for Edge-Case Judgment**

When you encounter a scenario not explicitly covered:
- Check Context > Problem: What is this system for?
- Check Context > Approach: How should I think about solutions?
- Check Context > Scope: Is this in or out?

Example: Spec doesn't say how to handle empty field in PUT request. Context says "Internal tool prioritizing simplicity" → reject invalid input, don't add complex partial update logic.

## Listing Sections

To see the spec structure:

```bash
cca list
```
