# goprops

**Property-based testing for Go using group theory**

Make mathematical properties (associativity, commutativity, identity, inverse, closure) easy to test in Go code.

## Why?

Traditional unit tests check specific inputs. Property tests check **rules that should always hold**:

```go
// Unit test: Does 2 + 3 = 5?
assert.Equal(t, 5, add(2, 3))

// Property test: Is addition associative for ALL numbers?
goprops.Associative(t, add)
```

The property test generates hundreds of random inputs and verifies the mathematical law holds.

## Installation

```bash
go get github.com/alexshd/goprops
```

## Quick Start

```go
import (
    "testing"
    "github.com/alexshd/goprops"
)

func TestCacheProperties(t *testing.T) {
    cache := NewCache()

    // Test that cache operations are associative
    goprops.Associative(t, func(a, b int) int {
        cache.Set("key", a)
        cache.Set("key", b)
        return cache.Get("key")
    })
}

func TestMergeCommutative(t *testing.T) {
    // Test that merge order doesn't matter
    goprops.Commutative(t, func(a, b map[string]int) map[string]int {
        return merge(a, b)
    })
}
```

## Properties Supported

- **Associative**: `(a ∘ b) ∘ c = a ∘ (b ∘ c)`
- **Commutative**: `a ∘ b = b ∘ a`
- **Identity**: `a ∘ e = a` (e is identity element)
- **Inverse**: `a ∘ a⁻¹ = e` (inverse exists)
- **Closure**: `a ∘ b` produces same type as inputs

## Status

🚧 **In Development** (3-4 day build project)

## License

MIT
