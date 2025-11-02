# Concurrency Safety Testing with lawtest

## Overview

We've enhanced `lawtest` with concurrency safety testing capabilities. This allows you to test if your operations are:

1. **Parallel-safe**: Can be executed concurrently without race conditions
2. **Immutable**: Don't mutate their inputs (pure functions)
3. **Associative under concurrency**: Mathematical properties hold even under concurrent load

## New Functions in lawtest.go

### 1. `ParallelSafe[T comparable](t *testing.T, op BinaryOp[T], gen Generator[T], goroutines int) bool`

Tests if an operation can be safely executed concurrently.

```go
op := func(a, b *GoodCache) *GoodCache {
    return a.Merge(b).(*GoodCache)
}

gen := func() *GoodCache {
    cache := NewGoodCache()
    cache.Set("x", 1)
    return cache
}

// Test with 20 concurrent goroutines
isSafe := lawtest.ParallelSafe(t, op, gen, 20)
```

**What it tests:**

- Launches multiple goroutines executing the operation
- Detects panics and race conditions
- Returns `true` if parallel-safe, `false` otherwise

### 2. `TestParallelAssociativity[T comparable](t *testing.T, op BinaryOp[T], gen Generator[T], goroutines int)`

Tests if associativity holds under concurrent execution.

```go
lawtest.TestParallelAssociativity(t, op, gen, 10)
```

**What it tests:**

- First verifies sequential associativity: `(a∘b)∘c = a∘(b∘c)`
- Then tests the same property under concurrent load
- Ensures mathematical properties survive concurrent execution

### 3. `ImmutableOp[T comparable](t *testing.T, op BinaryOp[T], gen Generator[T])`

Tests if an operation creates new values instead of mutating inputs.

```go
lawtest.ImmutableOp(t, op, gen)
```

**What it tests:**

- Verifies that applying the operation doesn't change input values
- Ensures operations are pure functions
- Critical for parallel-safe operations

## Why This Matters

### Immutable Operations → Parallel-Safe

```go
// ✅ GOOD: Creates new instance (immutable)
func (c *GoodCache) Merge(other Cacher) Cacher {
    result := NewGoodCache()
    maps.Copy(result.data, c.data)        // Copy, don't mutate
    maps.Copy(result.data, otherGood.data)
    return result  // New instance
}
```

**Properties:**

- ✅ Can run in parallel (no shared mutable state)
- ✅ Associative: `(A∪B)∪C = A∪(B∪C)`
- ✅ Deterministic: same inputs → same output
- ✅ Thread-safe by design

### Mutating Operations → Race Conditions

```go
// ❌ BROKEN: Mutates receiver (side effects)
func (c *BrokenCache) Merge(other Cacher) Cacher {
    for k, v := range otherBroken.data {
        if _, exists := c.data[k]; !exists {
            c.data[k] = v  // BUG: Mutates c!
        }
    }
    return c  // Returns mutated self
}
```

**Problems:**

- ❌ Cannot run in parallel (race conditions)
- ❌ Not associative: `(A∪B)∪C ≠ A∪(B∪C)`
- ❌ Non-deterministic: result depends on execution order
- ❌ NOT thread-safe

## Demonstration Tests

### Test 1: Parallel Safety

```go
func TestGoodCacheParallelSafety(t *testing.T) {
    isSafe := lawtest.ParallelSafe(t, op, gen, 20)
    // ✅ Returns true - GoodCache is immutable
}

func TestBrokenCacheParallelSafety(t *testing.T) {
    isSafe := lawtest.ParallelSafe(t, op, gen, 20)
    // May return false - BrokenCache mutates
}
```

### Test 2: Immutability

```go
func TestGoodCacheImmutability(t *testing.T) {
    cache1 := NewGoodCache()
    cache1.Set("x", 1)

    result := cache1.Merge(cache2)

    // ✅ cache1 unchanged (immutable)
    // ✅ result is new instance
}

func TestBrokenCacheImmutability(t *testing.T) {
    cache1 := NewBrokenCache()
    cache1.Set("x", 1)

    result := cache1.Merge(cache2)

    // ❌ cache1 modified (mutation)
    // ❌ result == cache1 (same instance)
}
```

### Test 3: Parallel Associativity

```go
func TestParallelAssociativityComparison(t *testing.T) {
    // GoodCache: Associativity holds even under concurrent load
    lawtest.TestParallelAssociativity(t, goodOp, goodGen, 10)

    // BrokenCache: May fail due to mutation + concurrency
    lawtest.TestParallelAssociativity(t, brokenOp, brokenGen, 10)
}
```

## Mathematical Properties & Concurrency

| Property          | Sequential                          | Concurrent                            | Parallel-Safe                 |
| ----------------- | ----------------------------------- | ------------------------------------- | ----------------------------- |
| **Associativity** | `(a∘b)∘c = a∘(b∘c)`                 | Must hold under concurrent execution  | Required for parallelization  |
| **Immutability**  | `op(a,b)` doesn't mutate `a` or `b` | No shared mutable state               | Enables lock-free concurrency |
| **Determinism**   | Same inputs → same output           | Same inputs → same output (any order) | Reproducible results          |

## Key Insights

1. **Immutability enables parallelism**: Pure functions can run concurrently without locks
2. **Algebraic properties matter**: Associativity isn't just math - it determines if operations can be reordered/parallelized
3. **Testing reveals design**: lawtest tests expose architectural issues (mutation, side effects)
4. **Property-based → architecture-based**: Mathematical properties map directly to concurrency safety

## Running Tests

```bash
# Run all concurrency tests
go test -v -run Parallel

# With race detector (detects actual data races)
go test -race -v -run Parallel

# Test immutability
go test -v -run Immutability

# Full test suite
go test -v
```

## Use Cases

### When to use these tests:

1. **Cache implementations**: Verify merge operations are parallel-safe
2. **Data structures**: Ensure operations don't corrupt under concurrent load
3. **Functional code**: Validate immutability guarantees
4. **Distributed systems**: Test if operations can be safely replicated/partitioned
5. **Performance optimization**: Identify which operations can be parallelized

### What you learn:

- ✅ Can this operation run in parallel?
- ✅ Does this operation mutate its inputs?
- ✅ Will mathematical properties hold under concurrency?
- ✅ Is this architecture thread-safe by design?

## Example Output

```
=== Testing GoodCache Parallel Safety ===
✅ Operation appears parallel-safe (no race conditions in 20 goroutines)
✅ GoodCache.Merge is parallel-safe (immutable operations)

=== Testing GoodCache Immutability ===
✅ cache1 not mutated (immutable)
✅ cache2 not mutated (immutable)
✅ Result is new instance (immutable merge)

=== Testing BrokenCache Immutability ===
✅ Correctly detected mutation! cache1 length changed: 1 → 2
✅ Correctly detected: result is SAME instance as input (violates immutability)
```

## Conclusion

By combining group theory with concurrency testing, we can:

1. **Verify mathematical properties** (associativity, idempotence, etc.)
2. **Test concurrency safety** (parallel execution, race conditions)
3. **Validate architectural decisions** (immutability, pure functions)
4. **Identify parallelization opportunities** (which operations can be concurrent)

This is the power of property-based testing: **mathematical laws reveal real-world architectural properties**.
