# Go-Specific Patterns

> **Note:**
> This section contains ***OPTIONAL*** language-specific guidance.
> Apply only if using Go.
> Core principles are universal.

## Text Templates

### When Applicable

When building dynamic text output (prompts, reports, templates).

### Principle

Use `text/template` instead of manual string concatenation.

### Rationale

* Separates content from code
* Enables non-programmers to modify templates
* Version control tracks template changes separately
* Testing template variations easier

### Decision Framework

Ask: Am I building strings with dynamic values?

| Scenario | Approach |
|----------|----------|
| YES + values come from struct | `text/template` |
| YES + simple interpolation | `fmt.Sprintf` acceptable |
| NO | Just use string literal |

## Embed Directive

### Principle

Use `//go:embed` ***ONLY*** for content the binary needs at ***RUNTIME***.

### Critical Question

Does running binary read this content?

### Decision Tree

**Runtime (binary execution):**

* Content read during execution → `//go:embed`
* Static, invariant → `//go:embed`
* Part of application logic → `//go:embed`

**Build-time (compilation):**

* Codegen input (GraphQL, OpenAPI, Proto) → External file
* Build tools read, generate code → External file

**Test-time (testing):**

* Test fixtures, vectors, golden files → External file
* Test suite loads from filesystem → External file

**Deploy-time (deployment):**

* K8s manifests, docker-compose → External file
* Deployment tools read these → External file

**Documentation:**

* Developer reference, examples → External file or omit
* Not used by code/build/test/deploy → Omit from spec

**Default:** External file unless clearly runtime-needed.

### Examples by Category

**Runtime candidates:**

* Static text needed by binary (system prompts, templates)
* SQL migrations if binary runs migrator at startup
* Static assets if binary serves them (web server, game assets)

**Build-time (NOT embedded):**

* Schema definitions that generate code
* IDL files for code generation
* Build configuration

**Test-time (NOT embedded):**

* Test data files
* Validation datasets
* Reference outputs

Your domain determines what fits each category.

### Anti-Pattern

❌ Don't blindly embed everything in references/

❌ Don't embed test data

❌ Don't embed build-time codegen inputs

❌ Don't embed deployment configs

✅ Ask: "Does binary read this at runtime?"

✅ Only embed if answer is clearly YES

## Constants/Config/Secrets

### Principle

Three distinct concerns with different lifecycles. ***MUST NOT mix.***

### Constants

* **What:** Algorithmic values that never change
* **Where:** Go const declarations in source code
* **Characteristics:**
  - Mathematical invariants
  - Derived from theory/validation
  - Part of algorithm itself
  - Never vary by environment
* **Decision:** Is this value algorithmic or operational?
  - Algorithmic → const in code

### Config

* **What:** Operational parameters varying by environment
* **Where:** Kubernetes ConfigMap (YAML)
* **Characteristics:**
  - Infrastructure endpoints
  - Resource limits
  - Feature flags
  - Different dev/staging/prod
* **Decision:** Does this change between environments?
  - Yes → K8s ConfigMap

### Secrets

* **What:** Sensitive data requiring secure storage
* **Where:** Kubernetes Secret (YAML, base64)
* **Characteristics:**
  - API keys, passwords, tokens
  - Security-critical
  - Never in version control
  - RBAC-controlled access
* **Decision:** Is this sensitive?
  - Yes → K8s Secret

### Anti-Pattern

❌ WRONG: Mixing all three in .env file

❌ WRONG: Embedding configuration

❌ WRONG: Secrets in ConfigMap

❌ WRONG: Constants in environment variables

✅ CORRECT: Separate concerns completely

✅ Constants in code

✅ Config in ConfigMap

✅ Secrets in Secret

### Note

This pattern specific to K8s deployments.

Other deployment targets have equivalent separation:

* Docker Compose: env_file vs secrets
* Serverless: environment vs secrets manager
* Traditional: config file vs secrets vault

Principle is universal: separate constants/config/secrets.

Mechanism varies by platform.
