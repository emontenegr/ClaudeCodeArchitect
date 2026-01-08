---
name: adoc
version: 0.1.0
description: Handle AsciiDoc (.adoc) files. Use when you see .adoc files, MANIFEST.adoc, or spec/ folders. Never read .adoc files directly - use spec-cli to compile them first.
---

# Architecture Specification Workflow

This project uses `spec-cli` to manage architecture specifications. The spec is the source of truth - implement from it, don't guess.

## Reading the Spec

Before implementing, read the compiled spec:

```bash
spec-cli compile
```

This outputs the full specification with all attributes resolved (e.g., `{api-p99-latency}` becomes `100ms`).

For specific sections:

```bash
spec-cli compile --section "API Endpoints"
spec-cli compile --section "Database Schema"
spec-cli compile --section core/types.adoc
```

## Validation

Before starting implementation, validate the spec is complete:

```bash
spec-cli validate --quick   # Fast structural checks
spec-cli validate           # Full validation with semantic analysis
```

If validation fails, the spec needs fixes before implementation.

## Checking Attribute Impact

Before changing a spec attribute value, check what's affected:

```bash
spec-cli impact <attribute-name>
```

Example:
```bash
spec-cli impact api-p99-latency
```

Shows all sections using that attribute so you understand the change scope.

## Listing Sections

To see the spec structure:

```bash
spec-cli list
```

## Implementation Workflow

1. **Read the spec first** - Run `spec-cli compile` to understand what you're building
2. **Check specific sections** - Use `--section` for targeted reading
3. **Implement exactly as specified** - The spec is complete; don't add unspecified features
4. **Validate before finishing** - Run `spec-cli validate` to confirm spec compliance

## Key Principles

- **Spec is truth**: If something isn't in the spec, ask before implementing it
- **Attributes are resolved**: You see actual values, not placeholders
- **Sections are navigable**: Use `--section` to focus on relevant parts
- **Changes have impact**: Use `impact` before modifying attribute values
