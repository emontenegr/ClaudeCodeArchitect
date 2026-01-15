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

**Installation:**

If `cca` is not available:
```bash
go install github.com/emontenegr/ClaudeCodeArchitect/cmd/cca@latest
```

Requirements: Go 1.21+, Claude Code >= 2.1, asciidoctor CLI

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

## Core Principle: SPEC IS TRUTH

**CONSTITUTIONAL RULE:** Never implement before the spec. Always adhere to the spec.

This means:
- Spec changes BEFORE code changes (never reverse)
- Code implements what spec says (never what you think it should say)
- Spec-code divergence is a critical failure

**The Order:**
1. Change spec files (if needed)
2. Commit spec changes
3. Verify spec with `cca compile`
4. Then implement code

Never reverse this order.

**If you find yourself thinking:**
- "I'll implement the code, then update spec to match"
- "Spec is wrong, but I know what's right, I'll just implement"
- "I'll update spec later, let me code first"

**STOP.** You are violating CCA's core principle.

**When Plans Involve Spec Changes:**

If your work requires updating the spec:
```bash
# 1. Edit spec files first
vim spec/stages/scoping.adoc

# 2. Commit spec changes
git commit -m "Update scoping stage to use CLI input"

# 3. Verify spec was actually changed
cca compile | grep [keyword]

# 4. Only then write code
```

Implementing code before spec changes creates spec-code divergence.

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

**7. Verification Protocol Before Declaring Complete**

CONSTITUTIONAL RULE: You may NOT declare implementation complete until passing this verification.

**IMPORTANT:** This is an INTERNAL checklist.

**Use thinking blocks for verification.** Do NOT output verification steps to user.

Your thinking block should contain:
- Section mapping (spec sections → implementation commits)
- Divergence check (spec vs implementation comparison)
- Gap detection (scenarios not covered)
- Completion determination

Your user-facing output should ONLY be:
- "Implementation complete. All spec sections implemented and verified." OR
- "Spec gap found: [specific scenario]. Section [name] doesn't specify [behavior]. Need clarification." OR
- "Spec divergence: Spec shows [X], implemented [Y] because [reason]. Spec should be updated."

**Step 1: Spec Faithfulness - Section Mapping**

**CRITICAL:** Verify against the ACTUAL spec content, not what you think spec should say.

Run: `cca compile` and read the current spec.

List every spec section with implementation status:

```
✓ Context - N/A (informational)
✓ Core Types - implemented in commit abc123
✓ Database Schema - implemented in commit def456
✗ Error States - NOT IMPLEMENTED
```

If ANY section shows "NOT IMPLEMENTED":
- Re-read that section carefully in the COMPILED spec
- Check Scope Out (is it explicitly excluded?)
- If required and not excluded: implement it now
- If unclear: STOP and ask user

**Spec-Code Alignment Check:**

Compare what spec says vs what you implemented:

```
Spec says: "TUI with Bubble Tea framework"
Code has: "CLI with stdin/stdout"
→ CRITICAL DIVERGENCE - code doesn't match spec AT ALL
```

If you find divergence:
- DON'T rationalize ("but CLI is better")
- DON'T continue ("I'll update spec later")
- STOP and report: "Spec describes [X], implementation is [Y]. Which is correct? If [Y], update spec first."

**Step 2: Spec Divergence Check**

Answer:

```
Q1: Did I implement differently than spec shows?
[List any place where code differs from spec examples/descriptions]

Q2: Did I make any choices the spec leaves unspecified?
[List any decisions you made that spec doesn't cover]
```

**If Q1 has items (divergence from spec):**
- STOP and report: "Spec divergence: [what spec says] vs [what I implemented] because [reason]. Spec should be updated."
- Even if justified (technical necessity), spec must stay accurate
- User updates spec, then you continue

**If Q2 has items (spec doesn't specify):**
- Check if Context (Abstract/Approach/Scope) guided the decision
- If Context supports it: OK (document in commit message)
- If not in spec or Context: STOP and ask user

**Example:**

❌ WRONG:
"Spec shows value receivers but Bubble Tea requires pointers. Changed to pointers. Justified by technical necessity."
→ You created spec divergence without reporting it.

✅ CORRECT:
"Spec divergence found: Spec shows value receivers (line 45), but Bubble Tea requires pointer receivers for state mutation. Implementation uses pointers. **Spec should be updated to show pointer receivers.** Stopping until spec is corrected."

**Step 3: Deferral Detection**

Search your implementation for:
- TODO or FIXME comments
- "Known gap" or "deferred" in your thinking
- Incomplete implementations

If found: These are NOT allowed. Either:
- Implement now (if spec requires it)
- Remove (if spec doesn't require it)
- STOP and ask user (if spec is unclear)

**Step 4: Spec Insufficiency Check**

While implementing, did you encounter scenarios where the spec was silent?

```
Scenario: [describe what spec doesn't cover]
Resolution: [how you handled it]
Source: [Context section that guided you, or "UNCLEAR"]
```

If any resolutions say "UNCLEAR":
- STOP implementing
- Report to user: "Spec doesn't specify [scenario]. Need clarification in [section]."
- Don't guess or invent behavior

**Step 5: Build + Behavior Verification**

- Code compiles/builds successfully
- Tests pass (if spec requires tests)
- Behavior matches spec EXACTLY (not "close enough")
- No runtime errors on basic usage

**Completion Criteria:**

✅ ALL must be true:
- Every spec section implemented (Step 1)
- Zero invented features (Step 2)
- Zero deferrals/TODOs (Step 3)
- Zero unresolved spec gaps (Step 4)
- Build + tests pass (Step 5)

❌ If ANY criterion fails: You are NOT complete.

**The Critical Distinction:**

**Incomplete implementation** = You didn't implement what spec says
→ Fix: Implement the missing spec section

**Incomplete spec** = Spec doesn't say what to do for scenario X
→ Fix: STOP, report to user, don't invent

Your job is faithful implementation, not gap-filling.

**Anti-Pattern Examples:**

❌ WRONG: "Parse error retry not implemented - deferring to phase 2"
→ Spec requires it. No phase 2. Implement now.

❌ WRONG: "Spec doesn't say how to handle timeout - implementing exponential backoff"
→ You just invented a feature. STOP and ask user.

❌ WRONG: "Known gap in error handling but critical features done"
→ There are no "critical vs non-critical" features in spec. All are required.

✅ CORRECT: "Spec section Error States requires retry logic. Implementing now per spec."

✅ CORRECT: "Encountered scenario: concurrent user creation. Spec doesn't address this. Stopping to ask user: should spec cover race conditions?"

✅ CORRECT: "All spec sections implemented. Zero TODOs. Build passes. Tests pass. Behavior verified against spec. Complete."

**How to Report Completion:**

❌ DON'T narrate the verification:
"Step 1: Section Coverage... ✓ Context implemented... ✓ Types implemented..."

✅ DO report concisely:
"Implementation complete. All spec sections implemented and verified."

OR if gap found:
"Spec gap found: Error handling for concurrent writes not specified in Database Schema section. Need clarification before proceeding."

## Listing Sections

To see the spec structure:

```bash
cca list
```
