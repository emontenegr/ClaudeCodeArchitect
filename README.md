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

## Core Principles

1. **Decided, Not Deciding** - Every choice made before writing spec
2. **Zero Conditionals** - No "if/when/maybe", single implementation path
3. **Complete Types** - Every field, every type, all constraints
4. **Quantified Performance** - Numbers, not adjectives (`P99 <100ms` not "fast")
5. **Mathematical Derivation** - Why every constant has its value
6. **Modular Structure** - File-per-concern using AsciiDoc includes

## Getting Started

### 1. Read the Principles
Start with foundational concepts:
- [Foundational Principles](docs/principles.md) - Decided not deciding, zero conditionals, completeness
- [Structural Requirements](docs/structure.md) - File organization, types, algorithms, data formats

### 2. Learn the Frameworks
Apply decision-making:
- [Decision Frameworks](docs/frameworks.md) - When to be explicit vs let AI decide
- [Anti-Patterns](docs/anti-patterns.md) - Common mistakes to avoid

### 3. Master Validation
Ensure completeness:
- [Validation Checklist](docs/validation.md) - 18-point completeness check
- [Testing Framework](docs/testing.md) - What to test, what not to test

### 4. Study Examples
See it in practice:
- [Simple API Example](examples/simple-api/) - Complete spec for user CRUD API

### 5. Check Language-Specific Patterns (Optional)
- [Go Patterns](docs/languages/go.md) - text/template, embed, constants/config/secrets

## Full Documentation

- **[Principles](docs/principles.md)** - Core concepts that drive everything
- **[Structure](docs/structure.md)** - How to organize specs modularly
- **[Frameworks](docs/frameworks.md)** - Decision-making for specification
- **[Anti-Patterns](docs/anti-patterns.md)** - What not to do
- **[Workflow](docs/workflow.md)** - Git-based evolution and compilation
- **[Testing](docs/testing.md)** - Test decision framework
- **[Validation](docs/validation.md)** - Completeness checking
- **[Implementation](docs/implementation.md)** - What to specify, what to omit

## Quick Start

### 1. Install Tools

```bash
# Install asciidoctor (required)
npm install -g @asciidoctor/cli

# Install cca
go install github.com/elijahmont3x/ClaudeCodeArchitect/cli/cmd/cca@v0.1.0
```

### 2. Install Claude Code Skill

```bash
# Install to current project
cca skill

# Or install globally (all projects)
cca skill --global
```

The skill teaches Claude Code to use cca automatically when working with AsciiDoc specs.

### 3. Create Your Spec

```bash
mkdir -p myproject/spec
vim myproject/spec/MANIFEST.adoc
```

### 4. Validate and Implement

```bash
cd myproject

# Validate spec completeness
cca validate

# Give to Claude Code
cca compile  # Claude reads the output
```

## What Makes This Different

**Traditional docs:** Explain what exists
**Planning docs:** Explore possibilities
**This system:** Define what to build

Specifications are:
- **Complete** (all decisions made)
- **Quantified** (numbers, not words)
- **Modular** (file-per-concern)
- **Executable** (AI implements directly)

***Not planning. Implementation blueprint.***

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

## Success Metrics

- 90-95% completion rate without iteration
- Zero "what should I do here?" questions from AI
- Quantified, measurable requirements throughout
- Complete on first implementation attempt

## Domain Agnostic

Works for any domain:
- Web APIs
- Game engines
- ML pipelines
- IoT firmware
- Financial systems
- Mobile apps

Principles are universal. Your domain determines specifics.

## Philosophy

> **"Specifications are compiler input, not human documentation."**

Write for the machine that will implement it.
Make every decision before writing.
Leave nothing ambiguous.
Quantify everything.

***The result:*** AI that builds exactly what you need, first time.

---

**License:** MIT
**Target AI:** Claude Code
**Completion Rate:** 90-95%
