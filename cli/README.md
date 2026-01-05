# spec-cli

CLI tool for compiling architecture specifications for AI consumption.

## Purpose

Provides a simple interface for AI assistants (like Claude Code) to read architecture specifications written in AsciiDoc.

## Installation

```bash
# From ClaudeCodeArchitect repo
task build-cli

# Or install to system
task install-cli
```

## Usage

### Compile Specification

Compiles AsciiDoc spec to Markdown and outputs to stdout:

```bash
spec-cli compile
```

Claude Code workflow:
```bash
# Claude runs this to read your spec
spec-cli compile
# Markdown output appears in bash result, Claude reads it directly
```

### Validate Specification

Checks if spec compiles without errors:

```bash
spec-cli validate
```

### Help

```bash
spec-cli help
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

Or in a subdirectory:

```
myproject/
├── spec/
│   └── MANIFEST.adoc
└── src/
```

spec-cli automatically finds:
- `MANIFEST.adoc`
- `spec/MANIFEST.adoc`
- `plan/MANIFEST.adoc`

### Option 2: Explicit Configuration

Create `.spec.yaml` in your project root:

```yaml
spec: ./path/to/MANIFEST.adoc
```

Paths can be:
- Relative to project root
- Absolute paths
- Remote repos (future feature)

## How It Works

1. **Finds spec** via `.spec.yaml` or convention
2. **Compiles AsciiDoc** to HTML (libasciidoc)
3. **Converts HTML** to Markdown (html-to-markdown)
4. **Outputs Markdown** to stdout
5. **AI reads** Markdown from bash output

No intermediate files needed - Claude reads compilation output directly.

## Requirements

- Go 1.25+ (for building)

No external dependencies - pure Go implementation.

## AI Workflow

When Claude Code needs to read your spec:

```bash
# Claude runs:
cd myproject
spec-cli compile

# Gets Markdown output with:
# - All includes resolved
# - All attributes substituted
# - Clean, readable Markdown
# - No HTML noise
```

## Development

```bash
# Build
cd cli
go build -o ../bin/spec-cli ./cmd/spec-cli

# Test
cd ../examples/simple-api
../../bin/spec-cli validate
../../bin/spec-cli compile > test.md
```
