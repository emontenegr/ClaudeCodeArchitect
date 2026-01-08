# cca

CLI tool for compiling and validating architecture specifications for AI consumption.

## Purpose

Provides an interface for AI assistants (like Claude Code) to read, validate, and navigate architecture specifications written in AsciiDoc.

## Installation

```bash
# If you have Go installed
go install github.com/elijahmont3x/ClaudeCodeArchitect/cli/cmd/cca@latest

# Or build from source
cd cli
go build -o cca ./cmd/cca
```

## Commands

### compile

Compiles AsciiDoc spec to Markdown:

```bash
cca compile                      # Full spec to stdout
cca compile --section "API Spec" # Single section with attrs resolved
cca compile --section core/types.adoc
```

### validate

Validates spec completeness using structural checks + Claude semantic analysis:

```bash
cca validate           # Full validation (structural + Claude)
cca validate --quick   # Structural checks only (fast, no Claude)
cca validate --yes     # Skip size confirmation for large specs
```

Requires Claude CLI for semantic validation. Install from: https://claude.ai/code

### diff

Compares compiled output between commits:

```bash
cca diff           # Compare with HEAD~1
cca diff HEAD~3    # Compare with 3 commits ago
cca diff main      # Compare with main branch
```

### impact

Shows which sections use an attribute:

```bash
cca impact api-p99-latency
```

Output:
```
Attribute: api-p99-latency
Defined in: MANIFEST.adoc:4 = "100ms"

Used in:
  - MANIFEST.adoc:93 (Section: "Performance Requirements")
  - MANIFEST.adoc:454 (Section: "Performance Specifications")
```

### list

Lists all sections in the spec:

```bash
cca list
```

## Configuration

### Option 1: Convention

Place `MANIFEST.adoc` in your project:

```
myproject/
├── MANIFEST.adoc
├── core/
├── concerns/
└── src/
```

cca automatically finds:
- `MANIFEST.adoc`
- `spec/MANIFEST.adoc`
- `plan/MANIFEST.adoc`

### Option 2: Explicit Configuration

Create `.spec.yaml` in your project root:

```yaml
spec: ./path/to/MANIFEST.adoc
```

## How It Works

1. **Finds spec** via `.spec.yaml` or convention
2. **Compiles AsciiDoc** to HTML (asciidoctor)
3. **Converts HTML** to Markdown (html-to-markdown)
4. **Outputs** to stdout for AI consumption

Key feature: Attributes like `{api-p99-latency}` are resolved during compilation, so Claude sees actual values, not placeholders.

## Validation

The `validate` command runs two phases:

**Phase 1: Structural Checks (Go)**
- Spec compiles
- Structure parseable
- Has sections
- Has key sections (types, API, deployment, etc.)

**Phase 2: Semantic Validation (Claude)**
- 18-point completeness checklist
- Context-aware issue detection
- No false positives on error message examples

Large specs (>20KB) prompt for confirmation. Use `--yes` to skip in CI.

## Requirements

- **asciidoctor** (for AsciiDoc compilation)
  ```bash
  # Windows/Mac/Linux with Node.js
  npm install -g @asciidoctor/cli

  # Or with Ruby
  gem install asciidoctor
  ```
- **Claude CLI** (for semantic validation) - https://claude.ai/code
- Go 1.21+ (only if building from source)

## AI Workflow

```bash
# Claude reads your spec
cca compile

# Claude checks a specific section
cca compile --section "Database Schema"

# Claude checks attribute impact before changes
cca impact cache-ttl

# Validate before implementation
cca validate
```
