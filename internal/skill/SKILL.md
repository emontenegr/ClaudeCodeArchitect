---
name: adoc
description: Handle AsciiDoc (.adoc) files. Use when you see .adoc files, MANIFEST.adoc, or spec/ folders. Never read .adoc files directly - use caca to compile them first.
---

# Architecture Specification Workflow

This project uses `caca` to manage architecture specifications. The spec is the source of truth - implement from it, don't guess.

## Reading the Spec

Before implementing, read the compiled spec:

```bash
caca compile
```

This outputs the full specification with all attributes resolved (e.g., `{api-p99-latency}` becomes `100ms`).

For specific sections:

```bash
caca compile --section "API Endpoints"
caca compile --section "Database Schema"
caca compile --section core/types.adoc
```

## Validation

Before starting implementation, validate the spec is complete:

```bash
caca validate --quick   # Fast structural checks
caca validate           # Full validation with semantic analysis
```

If validation fails, the spec needs fixes before implementation.

## Checking Attribute Impact

Before changing a spec attribute value, check what's affected:

```bash
caca impact <attribute-name>
```

Example:
```bash
caca impact api-p99-latency
```

Shows all sections using that attribute so you understand the change scope.

## Listing Sections

To see the spec structure:

```bash
caca list
```

## Implementation Workflow

1. **Read the spec first** - Run `caca compile` to understand what you're building
2. **Check specific sections** - Use `--section` for targeted reading
3. **Implement exactly as specified** - The spec is complete; don't add unspecified features
4. **Validate before finishing** - Run `caca validate` to confirm spec compliance

## Key Principles

- **Spec is truth**: If something isn't in the spec, ask before implementing it
- **Attributes are resolved**: You see actual values, not placeholders
- **Sections are navigable**: Use `--section` to focus on relevant parts
- **Changes have impact**: Use `impact` before modifying attribute values
