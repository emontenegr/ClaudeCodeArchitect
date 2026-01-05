# Simple User API Example

This is a complete, minimal specification demonstrating the ClaudeCodeArchitect principles.

## What This Example Shows

A user CRUD API with authentication that demonstrates:

1. **Decided, Not Deciding** - Every choice made (`PostgreSQL 16`, `Redis 7`, `JWT auth`, `bcrypt cost=12`)
2. **Zero Conditionals** - No "if needed" or "maybe add" - single implementation path
3. **Complete Types** - Full `User` and `Session` structs with all fields
4. **Quantified Performance** - `P99 <100ms`, `P50 <20ms`, `1000 req/s`, `100 req/min` rate limit
5. **Mathematical Derivation** - Connection pool size *calculated* from throughput requirements
6. **Modular Structure** - Single file for simplicity, but shows all required sections

## The Specification

See [MANIFEST.adoc](MANIFEST.adoc) for the complete specification.

This spec is intentionally kept in a single file for readability in this simple example. For larger projects, you would split it into:

```
simple-api/
├── MANIFEST.adoc
├── core/
│   ├── metadata.adoc
│   └── types.adoc
├── interfaces/
│   └── api.adoc
├── concerns/
│   └── performance.adoc
└── operations/
    ├── deployment.adoc
    └── testing.adoc
```

## What Makes This Complete

**Type Definitions:**
- All struct fields with exact types
- All invariants specified
- Database schema complete

**API Specification:**
- All routes with methods
- Request/response examples with real data
- All error cases (`400`, `401`, `404`, `409`)
- Authentication and authorization specified

**Performance:**
- Quantified latency (`P99/P50`)
- Throughput target (`1000 req/s`)
- Rate limiting specified
- Connection pool size derived from throughput

**Dependencies:**
- All libraries with exact versions
- Database with exact version
- Cache with exact version

**Deployment:**
- Platform specified (`Docker Compose`)
- Complete `docker-compose.yml`
- Environment variables defined

**Testing:**
- Coverage targets specified
- What to test defined
- What NOT to test defined
- Test framework specified

## How to Use This

1. **Read the spec** - [MANIFEST.adoc](MANIFEST.adoc)
2. **Notice completeness** - Every decision is made
3. **See quantification** - All performance as numbers
4. **Observe single path** - No conditionals or options

If you gave this spec to Claude Code, you would get a working implementation with ~90-95% completion.

## Key Learnings

**Before (Incomplete):**
```
Build a user API
- CRUD operations
- Authentication
- Use a database
```

**After (Complete):**
```
User API with JWT auth
- PostgreSQL 16 (exact schema provided)
- Redis 7 for sessions (TTL: 300s)
- bcrypt cost=12 for passwords
- P99 latency <100ms
- Rate limit: 100/min per IP
- All 5 endpoints fully specified
```

***The difference is completeness, quantification, and decision-making.***

## Try It Yourself

Compile this spec and give it to Claude:

```bash
# Compile to markdown
asciidoctor -b markdown -o compiled-spec.md MANIFEST.adoc

# Give to Claude Code
claude-code "Read compiled-spec.md and implement the user API"
```

You should get a working implementation without needing to answer questions.
