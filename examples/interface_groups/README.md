# lawtest Examples: Interface-Based Property Testing with Concurrency

This directory demonstrates **user-defined interface testing** with lawtest, including **concurrency safety verification**.

## What's in this directory

### `demonstration.go` - User-Defined Interface Example

Shows how to define YOUR OWN interface (not from lawtest) and test implementations:

```go
// USER-DEFINED INTERFACE (you create this!)
type Cacher interface {
    Set(key string, value int)
    Get(key string) (int, bool)
    Merge(other Cacher) Cacher
    Clear()
}

// ✅ CORRECT: Immutable implementation
type GoodCache struct { ... }

// ❌ BROKEN: Mutating implementation  
type BrokenCache struct { ... }
```

**Run it:**
```bash
go run demonstration.go
```

### `demonstration_test.go` - Comprehensive Property Tests

Tests both implementations with lawtest to reveal bugs:

1. **Semigroup properties**: Associativity, closure
2. **Parallel safety**: Can it run concurrently?
3. **Immutability**: Does it mutate inputs?
4. **Concurrency stress tests**: Properties under load

**Run tests:**
```bash
# All tests
go test -v

# Specific property
go test -v -run TestGoodCacheSemigroup

# Concurrency tests
go test -v -run "Immutability|ParallelSafety"

# With race detector
go test -race -v
```

### `interface_examples.go` - lawtest Interface Implementations

Examples of implementing lawtest built-in interfaces:
- `Group[T]`: Groups with identity and inverse
- `Monoid[T]`: Monoids with identity
- `Semigroup[T]`: Associative operations
- `IdempotentOp[T]`: Idempotent functions

### `interface_examples_test.go` - Testing Interface Implementations

Tests for all the lawtest interface implementations.

## Key Concepts Demonstrated

### 1. User-Defined Interfaces

You can test **ANY** interface you create with lawtest:

```go
// Your interface
type MyInterface interface {
    Operation(other MyInterface) MyInterface
}

// Implement Semigroup for testing
type MySemigroup struct{}

func (s MySemigroup) Op(a, b *MyImpl) *MyImpl {
    return a.Operation(b).(*MyImpl)
}

func (s MySemigroup) Gen() *MyImpl {
    return NewMyImpl() // your generator
}

// Test it!
lawtest.TestSemigroup(t, MySemigroup{})
```

### 2. Immutability = Concurrency Safety

**GoodCache** (immutable):
```go
func (c *GoodCache) Merge(other Cacher) Cacher {
    result := NewGoodCache()      // ✅ New instance
    maps.Copy(result.data, c.data) // ✅ Copy, don't mutate
    return result                   // ✅ Pure function
}
```

**Properties:**
- ✅ Parallel-safe (no shared mutable state)
- ✅ Associative
- ✅ Deterministic
- ✅ Can be parallelized

**BrokenCache** (mutating):
```go
func (c *BrokenCache) Merge(other Cacher) Cacher {
    for k, v := range other.data {
        c.data[k] = v  // ❌ Mutates receiver!
    }
    return c           // ❌ Returns mutated self
}
```

**Problems:**
- ❌ NOT parallel-safe (race conditions)
- ❌ NOT associative (first-write-wins)
- ❌ Non-deterministic
- ❌ Cannot be parallelized

### 3. New Concurrency Testing Functions

#### `ParallelSafe(t, op, gen, goroutines)`

Tests if operation can run concurrently:

```go
isSafe := lawtest.ParallelSafe(t, op, gen, 20)
// true  = can parallelize
// false = has race conditions
```

#### `TestParallelAssociativity(t, op, gen, goroutines)`

Tests if associativity survives concurrent execution:

```go
lawtest.TestParallelAssociativity(t, op, gen, 10)
// Verifies: (a∘b)∘c = a∘(b∘c) even under concurrent load
```

#### `ImmutableOp(t, op, gen)`

Tests if operation mutates inputs:

```go
lawtest.ImmutableOp(t, op, gen)
// Verifies: op(a,b) doesn't change a or b
```

## Test Results

```bash
$ go test -v -run "Immutability|ParallelSafety"

=== RUN   TestGoodCacheParallelSafety
✅ Operation appears parallel-safe (no race conditions in 20 goroutines)
✅ GoodCache.Merge is parallel-safe (immutable operations)
--- PASS: TestGoodCacheParallelSafety (0.00s)

=== RUN   TestGoodCacheImmutability
✅ cache1 not mutated (immutable)
✅ cache2 not mutated (immutable)
✅ Result is new instance (immutable merge)
--- PASS: TestGoodCacheImmutability (0.00s)

=== RUN   TestBrokenCacheImmutability
✅ Correctly detected: result is SAME instance as input (violates immutability)
--- PASS: TestBrokenCacheImmutability (0.00s)
```

## What You Learn

1. **Can this operation run in parallel?** → `ParallelSafe()`
2. **Does this operation mutate?** → `ImmutableOp()`
3. **Are mathematical properties preserved under concurrency?** → `TestParallelAssociativity()`
4. **Is my architecture thread-safe by design?** → All of the above

## Mathematical Properties → Real-World Architecture

| Mathematical Property | Real-World Impact |
|----------------------|-------------------|
| **Associativity**: `(a∘b)∘c = a∘(b∘c)` | Operations can be reordered/parallelized |
| **Immutability**: `op(a,b)` doesn't mutate | Lock-free concurrency possible |
| **Idempotence**: `f(f(x)) = f(x)` | Retry-safe, duplicate-safe operations |
| **Commutativity**: `a∘b = b∘a` | Order-independent, CRDTs possible |

## Use Cases

- **Cache implementations**: Test merge strategies
- **Data structures**: Verify concurrent safety
- **Distributed systems**: Test if operations can be replicated
- **Functional programming**: Validate purity guarantees
- **Performance**: Identify parallelization opportunities

## Further Reading

- `CONCURRENCY_TESTING.md` - Detailed guide to concurrency testing
- `../../lawtest.go` - Core property testing functions
- `../../README.md` - Main lawtest documentation

## Quick Start

```bash
# See the demonstration
go run demonstration.go

# Run all tests
go test -v

# Test concurrency safety
go test -v -run Parallel

# With race detector
go test -race -v
```

## The Big Idea

**Mathematical properties aren't just abstract theory - they determine if your code can run in parallel, handle failures gracefully, and scale to distributed systems.**

lawtest lets you test these properties automatically, catching architectural issues before they become production bugs.
