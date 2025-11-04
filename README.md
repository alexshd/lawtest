# lawtest

**Property-based testing for Go using group theory**

Make mathematical properties (associativity, commutativity, identity, inverse, closure) easy to test in Go code.

## Why?

Traditional unit tests check specific inputs. Property tests check **rules that should always hold**:

```go
// Unit test: Does 2 + 3 = 5?
assert.Equal(t, 5, add(2, 3))

// Property test: Is addition associative for ALL numbers?
lawtest.Associative(t, add)
```

The property test generates hundreds of random inputs and verifies the mathematical law holds.

## Installation

```bash
go get github.com/alexshd/lawtest
```

## Quick Start

```go
import (
    "testing"
    "github.com/alexshd/lawtest"
)

func TestCacheProperties(t *testing.T) {
    cache := NewCache()

    // Test that cache operations are associative
    lawtest.Associative(t, func(a, b int) int {
        cache.Set("key", a)
        cache.Set("key", b)
        return cache.Get("key")
    })
}

func TestMergeCommutative(t *testing.T) {
    // Test that merge order doesn't matter
    lawtest.Commutative(t, func(a, b map[string]int) map[string]int {
        return merge(a, b)
    })
}
```

## Properties Supported

### Core Properties

- **Associative**: `(a ∘ b) ∘ c = a ∘ (b ∘ c)`
- **Commutative**: `a ∘ b = b ∘ a`
- **Identity**: `a ∘ e = a` (e is identity element)
- **Inverse**: `a ∘ a⁻¹ = e` (inverse exists)
- **Closure**: `a ∘ b` produces same type as inputs
- **Idempotent**: `f(f(x)) = f(x)`

### Concurrency Safety (New!)

- **ParallelSafe**: Can operations run concurrently without race conditions?
- **ImmutableOp**: Does the operation mutate its inputs?
- **TestParallelAssociativity**: Do properties hold under concurrent execution?

## Requirements

- Go 1.18 or higher (uses generics)

## Examples

See the `examples/` directory for comprehensive examples including:

- Interface-based testing (Group, Monoid, Semigroup)
- Concurrency safety testing
- User-defined interface testing

## Documentation

- [QUICKSTART.md](QUICKSTART.md) - Quick start guide
- [examples/interface_groups/README.md](examples/interface_groups/README.md) - Interface testing examples
- [examples/interface_groups/CONCURRENCY_TESTING.md](examples/interface_groups/CONCURRENCY_TESTING.md) - Concurrency testing guide

## License

Apache 2.0 - See [LICENSE](LICENSE) file for details.

Copyright 2025 Alex Shadrin ([@alexshd](https://github.com/alexshd))
