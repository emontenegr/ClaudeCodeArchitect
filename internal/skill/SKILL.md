---
name: adoc
description: Handle AsciiDoc (.adoc) files. Use when you see .adoc files, MANIFEST.adoc, or spec/ folders. Never read .adoc files directly - use cca to compile them first.
---

# Claude Code Architect (CCA)

## What This Is

This project uses **CCA** - a methodology and CLI for writing architecture specifications that AI can implement correctly on the first try.

**Goal:** 90-95% production-ready implementation from written spec, zero iteration.

**How:** Specifications are complete, decided, quantified - like compiler input. They contain:
- Context (why, scope, approach)
- Complete implementation details (types, schemas, APIs, performance)

## Command Reference

| Command | Purpose |
|---------|---------|
| `cca compile` | Compile spec to readable Markdown (resolves includes/attributes) |
| `cca compile --section <name>` | Compile specific section only |
| `cca validate` | Full validation (structural + semantic via Claude) |
| `cca validate --quick` | Fast structural checks only |
| `cca validate --ultra` | Enhanced validation (3x parallel + synthesis) |
| `cca diff [commit]` | Diff compiled output vs git commit |
| `cca impact <attribute>` | Show which sections use an attribute |
| `cca list` | List all sections in spec |
| `cca skill` | Install/update this Claude Code skill |

## Spec Structure

CCA specs are modular AsciiDoc files with two layers:

**Context** - Intent and foundation:
- Identity: What this system is (name, paradigm)
- Stack: Technical foundation (language, dependencies, versions)
- Abstract: What gap this addresses (optional but valuable)
- Approach: How problem is solved conceptually (optional)
- Scope: What's in, what's explicitly out (optional)

**Implementation** - Complete technical details:
- Types, schemas, APIs, algorithms
- Performance requirements (quantified)
- Deployment, testing specifications

Both are necessary. Context provides intent for edge-case judgment. Implementation provides mechanical completeness.

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
- **Abstract** (valuable for non-trivial systems): What gap this addresses
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

Validation checks for "Ambiguous Intent." If Abstract/Approach/Scope are absent and system is non-trivial, validation will flag this.

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
- What this system is for (Abstract)
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

**3. When Spec Code Doesn't Compile: Research, Don't Patch**

CRITICAL: If example code in the spec doesn't match the library API:

**DON'T:**
- Delete the code to make it compile
- Comment it out and move on
- Declare success when build passes

**DO:**
- Research the actual library API (use `go doc`, check docs, read source)
- Understand the INTENT of the spec code (what is it trying to achieve?)
- Implement that intent using the correct API
- If the API doesn't support it, ask the user

**Example:**
```
Spec shows: ta.PromptStyle = titleStyle
Build error: PromptStyle field doesn't exist
WRONG: Delete the line → build passes → done
RIGHT: Research bubbles/textarea styling API → find correct method → implement prompt styling
```

Build passing is NOT the success criterion. Correct implementation of intent is.

**4. Implement Exactly As Specified**

The spec is complete - don't add unspecified features. If something seems missing, it's either:
- In the spec (search more carefully)
- Explicitly excluded (check Scope section)
- Actually missing (validation should have caught this - ask user)

**5. Use Context for Edge-Case Judgment**

When you encounter a scenario not explicitly covered:
- Check Context > Abstract: What is this system for?
- Check Context > Approach: How should I think about solutions?
- Check Context > Scope: Is this in or out?

Example: Spec doesn't say how to handle empty field in PUT request. Context says "Internal tool prioritizing simplicity" → reject invalid input, don't add complex partial update logic.

**6. Commit Atomically As You Implement**

CRITICAL: Don't implement everything then commit once. Commit after each spec section.

**Workflow:**
```
Read spec section → Implement component → Verify it works → Commit → Next section
```

**Commit granularity:**
- Per spec section (Types, Database, API, Deployment)
- Or per logical component if section is large

**Commit message format:**
```
Implement [component name] (spec: [section name])

[1-2 sentences on what was implemented]

CCA-Spec: [section name]
```

**Example:**
```
Implement core type definitions (spec: Core Types)

Added User, Session, and Config structs with all fields,
validation rules, and invariants per specification.

CCA-Spec: Core Types
```

**Why this matters:**

1. **Verification:** Each commit is a checkpoint - you verify component works before moving on
2. **Traceability:** Git history maps to spec structure - easy to find which code implements which spec section
3. **Review:** User can review commit-by-commit instead of massive final diff
4. **Debugging:** If something breaks, git log shows exactly which component/section
5. **Spec evolution:** When spec changes, commit messages show what implementation needs updating

**What to verify before committing:**
- Component compiles (if applicable)
- Component tests pass (if you wrote tests for it)
- Integrates with previous commits (no breaking changes)

**Don't:**
- Commit with "WIP", "fix", "temp" messages
- Commit broken code to "fix later"
- Make one giant commit at the end
- Omit `CCA-Spec:` trailer (needed for traceability)

## Listing Sections

To see the spec structure:

```bash
cca list
```
