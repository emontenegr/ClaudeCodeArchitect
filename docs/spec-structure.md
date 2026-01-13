# Specification Structure Guide

## Minimum Viable Spec

Every spec needs:

1. **Summary** - One paragraph. What this system is and does.
2. **Scope** - What's included, what's explicitly excluded.
3. **Implementation Details** - The 18-point checklist material.

## Optional Sections

Add these when they provide value:

### Context
- **Problem Statement** - What gap or pain point this addresses
- **Constraints** - External limitations (budget, timeline, compliance, existing systems)

### Architecture
- **Conceptual Model** - How to think about the system (not implementation, but mental model)
- **Component Relationships** - How parts interact, data flow narrative

### Decisions
- **Key Trade-offs** - Choices made and why (e.g., "chose consistency over availability because X")
- **Alternatives Considered** - What was rejected and why

### Boundaries
- **Non-Goals** - What this explicitly won't do
- **Future Phases** - Concrete next steps (not vague "maybe later")

## When to Use What

| System Complexity | Recommended Sections |
|-------------------|---------------------|
| Script/utility | Summary, Scope, Implementation |
| Single service | + Problem Statement, Non-Goals |
| Multi-component | + Conceptual Model, Component Relationships |
| Distributed/critical | + Key Trade-offs, Alternatives Considered |

## Example: Minimal Spec

```asciidoc
= User API

== Summary

REST API for user CRUD operations. Supports authentication via JWT.

== Scope

Included: User registration, login, profile management.
Excluded: Password reset (separate service), admin operations.

== Implementation

[... 18-point details ...]
```

## Example: Full Spec

```asciidoc
= Payment Processing Service

== Summary

Processes credit card payments for checkout flow. Integrates with Stripe.
Prioritizes correctness over latency - failed charges must never result in
fulfilled orders.

== Problem Statement

Current system uses synchronous payment calls in the checkout flow.
Timeout failures cause order state inconsistencies requiring manual
reconciliation (~50 cases/month).

== Scope

Included: Card payments, refunds, webhook handling.
Excluded: Alternative payment methods, subscription billing, invoicing.

== Conceptual Model

Payments are processed as a state machine:
  pending -> authorized -> captured -> settled
           -> declined
           -> refunded

Each transition is idempotent. Duplicate requests return existing state.

== Key Trade-offs

Eventual consistency for settlement status.
- Stripe webhooks may arrive out of order
- We accept up to 5min delay in settlement confirmation
- Chose this over polling (cost) or synchronous confirmation (latency)

== Non-Goals

- PCI compliance for card storage (Stripe handles this)
- Multi-currency (USD only for v1)
- Partial refunds

== Implementation

[... 18-point details ...]
```

## Validation

The 18-point checklist validates **implementation completeness**.

Optional sections are not validated - their presence and depth is author judgment based on system complexity.
