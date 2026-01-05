# Git-Based Evolution Workflow

## The Problem

Specifications evolve during and after implementation:

* Implementation reveals *incomplete spec*
* Requirements *change*
* Performance targets need *adjustment*
* New constraints *discovered*

How do you efficiently communicate spec changes to Claude?

## The Solution: Git as Communication Channel

Use git diffs to show exactly what changed.

### Workflow

**1. Make spec changes**

```bash
$ vim concerns/performance.adoc
# Change :api-p99-latency: 100ms to 50ms
```

**2. Commit with clear message**

```bash
$ git add concerns/performance.adoc
$ git commit -m "Tighten API latency requirement from 100ms to 50ms P99

Based on load testing, current 100ms target is too loose.
Real-world P99 is 30ms, setting target to 50ms provides headroom."
```

**3. Show Claude the diff**

```bash
$ git diff HEAD~1 concerns/performance.adoc
```

**4. Communicate to Claude**

```
Read git diff HEAD~1 concerns/performance.adoc

I tightened the API latency requirements. Update implementation to meet new P99 target of 50ms.
```

### Benefits

**Precise:** Claude sees exactly what changed, no ambiguity

**Contextualized:** Commit message explains why

**Minimal:** Only reads changed files, not entire spec

**Trackable:** Git history shows spec evolution

**Auditable:** Can trace implementation changes to spec changes

### Best Practices

**Atomic commits:**

One logical change per commit.

```
✅ Good:
git commit -m "Add caching requirement for user queries"

❌ Bad:
git commit -m "Update performance, fix typos, add new endpoint"
```

**Clear commit messages:**

Explain WHAT changed and WHY.

```
✅ Good:
"Increase connection pool from 10 to 25 connections

Load testing showed pool exhaustion under 500 concurrent users.
New limit provides 2x headroom."

❌ Bad:
"Update config"
```

**Tag sections being updated:**

```bash
git commit -m "concerns/performance: Tighten API latency SLA"
```

Prefix shows which concern changed.

### Recompiling After Changes

After updating modular sources:

```bash
# Recompile for Claude
$ asciidoctor -b docbook -o compiled-spec.xml MANIFEST.adoc

# Or markdown
$ asciidoctor -b markdown -o compiled-spec.md MANIFEST.adoc

# Verify
$ wc -l compiled-spec.md
5234 compiled-spec.md

# Give to Claude
$ claude-code "Read compiled-spec.md and adjust implementation for updated requirements"
```

# Compilation for Claude

## AsciiDoctor Compilation

After writing modular specification, compile to single file:

**Markdown output (recommended for Claude):**

```bash
asciidoctor -b markdown -o compiled-spec.md MANIFEST.adoc
```

**DocBook output (alternative):**

```bash
asciidoctor -b docbook -o compiled-spec.xml MANIFEST.adoc
```

**HTML output (for human review):**

```bash
asciidoctor -b html5 -o compiled-spec.html MANIFEST.adoc
```

## Verification

Check compiled output size:

```bash
wc -l compiled-spec.md
```

Review for completeness:

```bash
less compiled-spec.md
# Verify all includes resolved
# Verify all attributes substituted
# Verify no broken references
```

## Giving to Claude

```bash
claude-code "Read compiled-spec.md and implement the system according to specification"
```

Or within Claude Code session:

```
Read compiled-spec.md - this is the complete specification. Implement the system exactly as specified.
```
