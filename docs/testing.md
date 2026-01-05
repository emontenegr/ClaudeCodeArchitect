# Testing Decision Framework

## The Problem

Testing *everything* wastes effort.

Testing *nothing* risks bugs.

**What should specifications require testing?**

## Decision Framework: Test YOUR Code

### Test YOUR Code

**Business logic you implemented:**

* Data transformations you wrote
* Validation rules you defined
* Algorithms you created
* State machines you designed

**Example:**

```asciidoc
Test: Email validation rejects invalid formats

Input: "invalid-email"
Expected: ValidationError("Invalid email format")

Input: "user@example.com"
Expected: Valid
```

**Integration points:**

* How your code connects to external systems
* Test the CONTRACT, not the library

**Example:**

```asciidoc
Test: User creation inserts to database

Given: Valid user data
When: createUser() called
Then:
- User record exists in database
- User ID is valid UUID v4
- createdAt timestamp is recent
```

**Edge cases in your algorithms:**

* Boundary conditions
* Error paths
* Unusual inputs

**Example:**

```asciidoc
Test: Search handles empty results

Given: Query "zzzzzzz" (no matches)
When: searchUsers() called
Then: Returns empty array, not null
```

### Do NOT Test Library Code

**Third-party libraries:**

❌ Don't test that Redis stores values correctly

❌ Don't test that bcrypt hashes correctly

❌ Don't test that PostgreSQL enforces constraints

❌ Don't test that Express routes correctly

**Why:** Libraries are tested by their maintainers. Redundant effort.

**Standard library:**

❌ Don't test that array.sort() works

❌ Don't test that JSON.parse() works

❌ Don't test that http.Get() works

**Why:** Standard library is battle-tested.

**Framework internals:**

❌ Don't test that Next.js renders pages

❌ Don't test that Django handles requests

❌ Don't test that React manages state

**Why:** Framework testing is framework's responsibility.

### Test Usage, Not Implementation

**Test how you USE libraries:**

**Do** test that you CALL bcrypt correctly:

```asciidoc
Test: Password hashing uses cost factor 12

Given: Password "SecurePass123!"
When: hashPassword() called
Then: bcrypt.compare() verifies with cost=12
```

**Do** test that you CONFIGURE Redis correctly:

```asciidoc
Test: Cache TTL is 300 seconds

Given: Cached value
When: 301 seconds elapse
Then: Cache returns miss
```

**The distinction:**

* Testing bcrypt library = ❌ (library concern)
* Testing you use bcrypt correctly = ✅ (your concern)

## Specifying Tests in PLAN

### Test Requirements Format

`operations/testing.adoc`:

```asciidoc
# Testing Requirements

## Coverage Targets

* Line coverage: 80% minimum
* Public API coverage: 100%
* Critical paths coverage: 100%

## What to Test

**Business Logic:**

* User validation (email format, username constraints)
* Authentication flow (login, token generation, refresh)
* Authorization logic (role-based access control)

**Integration Contracts:**

* Database operations (CRUD operations, transactions)
* Cache operations (set, get, TTL behavior)
* External API calls (request format, response handling)

**Edge Cases:**

* Empty inputs
* Null values
* Boundary conditions (max length, min value)
* Concurrent operations
* Error conditions

## What NOT to Test

* PostgreSQL query execution (database concern)
* Redis storage reliability (cache concern)
* bcrypt hashing correctness (library concern)
* Express routing mechanism (framework concern)

## Test Framework

* Framework: Jest 29.7.0
* Mocking: jest.mock() for external services
* Integration: Testcontainers for database/cache

## Example Test Cases

**User Creation:**

```typescript
describe('createUser', () => {
  it('creates user with valid data', async () => {
    const user = await createUser({
      username: 'john_doe',
      email: 'john@example.com',
      password: 'SecurePass123!'
    });

    expect(user.id).toMatch(/^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/);
    expect(user.username).toBe('john_doe');
    expect(user.passwordHash).not.toBe('SecurePass123!');
  });

  it('rejects invalid email', async () => {
    await expect(createUser({
      username: 'john_doe',
      email: 'invalid',
      password: 'SecurePass123!'
    })).rejects.toThrow('Invalid email format');
  });
});
```

### Reasoning

Clear testing requirements prevent:

* Testing Redis (wastes time)
* Missing critical business logic tests (risks bugs)
* Ambiguous coverage expectations (is 60% enough?)
* Inconsistent testing approaches

Specify what to test, what not to test, and how to test it.
