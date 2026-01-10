# Architecture Specification Authoring System

**Enable one-shot AI implementation from written specifications.**

Target AI: Claude Code | Target Completion: 90-95% production-ready

## What Is This?

A system for writing specifications that AI can implement correctly on the first try.

Not documentation. Not planning. ***Executable architecture as code.***

Specifications become compiler input - complete, decided, quantified - that produce working software at 90-95% completion without iteration.

## Why Use This?

**Bad spec:** "Use a database. Make it fast. Add caching if needed."

**Result:** AI asks questions, makes wrong choices, misses requirements.

**Good spec:** `PostgreSQL 16. P99 latency <100ms. Cache with 5min TTL (Redis 7).`

**Result:** AI implements exactly what you need, first time.

The difference is completeness, decision-making, and quantification.

## Quick Example

### Before (Incomplete)
```
Build a user API.
- Handle authentication
- Store users
- Make it scalable
```

**Problems:** What auth? Which database? How scalable?

### After (Complete)
```
User Authentication API
- OAuth 2.0 + JWT (RS256)
- PostgreSQL 16 (users table: id uuid, email varchar(255) unique, password_hash text)
- Support 1000 req/sec, P99 <100ms
- Rate limit: 100 requests/min per IP
```

**Result:** AI knows exactly what to build.

## Get Started

### 1. Install

**Requirements:** Claude Code >= 2.1, Go 1.21+

```bash
npm install -g @asciidoctor/cli
go install github.com/emontenegr/ClaudeCodeArchitect/cmd/cca@latest
```

### 2. Try It

```bash
cd your-project
cca skill              # Install Claude Code integration
cca validate --quick   # Check your spec
cca compile | head -50 # See compiled output
```

## Core Principles

1. **Decided, Not Deciding** - Every choice made before writing spec
2. **Zero Conditionals** - No "if/when/maybe", single implementation path
3. **Complete Types** - Every field, every type, all constraints
4. **Quantified Performance** - Numbers, not adjectives (`P99 <100ms` not "fast")
5. **Mathematical Derivation** - Why every constant has its value
6. **Modular Structure** - File-per-concern using AsciiDoc includes

## Commands

| Command | Purpose |
|---------|---------|
| `cca compile` | Full spec to Markdown |
| `cca compile --section <name>` | Single section with attributes resolved |
| `cca validate` | Structural + semantic completeness check |
| `cca validate --quick` | Structural checks only (fast, no Claude) |
| `cca diff [commit]` | Compiled output diff between commits |
| `cca impact <attr>` | Show sections using an attribute |
| `cca list` | List all sections |
| `cca skill` | Install Claude Code skill |

## How It Works

### Compilation Pipeline

1. Finds spec via `.spec.yaml` or convention (`MANIFEST.adoc`, `spec/MANIFEST.adoc`)
2. Compiles AsciiDoc to HTML via asciidoctor
3. Converts HTML to Markdown
4. Outputs to stdout

Key: `{api-p99-latency}` becomes `100ms` — Claude sees actual values, not placeholders.

### Validation Strategy

Two-phase validation:
1. **Structural (Go)**: Compiles? Has sections? Has key components?
2. **Semantic (Claude)**: 18-point completeness checklist

Large specs (>20KB) prompt for confirmation before Claude analysis. Use `--quick` for structural checks only.

### Requirements

- **asciidoctor**: AsciiDoc compilation — `npm install -g @asciidoctor/cli`
- **Claude CLI**: Semantic validation (optional, skip with `--quick`)

## Writing Specifications

Learn the methodology:

1. **[Principles](docs/principles.md)** - Decided not deciding, zero conditionals, completeness
2. **[Structure](docs/structure.md)** - File organization, types, algorithms, data formats
3. **[Frameworks](docs/frameworks.md)** - When to be explicit vs let AI decide
4. **[Anti-Patterns](docs/anti-patterns.md)** - Common mistakes to avoid
5. **[Validation](docs/validation.md)** - 18-point completeness check

See it in practice: **[Simple API Example](examples/simple-api/)**

## Example Structure

```
myproject/
├── MANIFEST.adoc              # Entry point, composes all parts
├── core/
│   ├── metadata.adoc          # System context
│   └── types.adoc             # Complete type definitions
├── concerns/
│   ├── performance.adoc       # Quantified requirements
│   └── security.adoc          # Security constraints
├── interfaces/
│   └── api.adoc               # Complete API spec
├── operations/
│   ├── deployment.adoc        # Deployment config
│   └── testing.adoc           # Test requirements
└── (implementation files)     # AI implements here
```

## Full Documentation

- **[Principles](docs/principles.md)** - Core concepts
- **[Structure](docs/structure.md)** - Modular organization
- **[Frameworks](docs/frameworks.md)** - Decision-making
- **[Anti-Patterns](docs/anti-patterns.md)** - What not to do
- **[Workflow](docs/workflow.md)** - Git-based evolution
- **[Testing](docs/testing.md)** - Test decision framework
- **[Validation](docs/validation.md)** - Completeness checking
- **[Implementation](docs/implementation.md)** - What to specify, what to omit

---

**License:** MIT
**Target AI:** Claude Code
**Completion Rate:** 90-95%
