# lawtest

**Property-based testing for Go using group theory**

Make mathematical properties (associativity, commutativity, identity, inverse, closure) easy to test in Go code.

## v0.2.0: Clean API

### Fluent API (New!)

For complex types, use the fluent `Laws` wrapper:

```go
func TestPacketPhysics(t *testing.T) {
    law := lawtest.For(t, packetGen, packetEq)

    law.Associative(add)      // ✅ (a∘b)∘c = a∘(b∘c)
    law.Commutative(add)      // ✅ a∘b = b∘a
    law.Identity(add, zero)   // ✅ a∘e = a
    law.Inverse(add, neg, zero) // ✅ a∘a⁻¹ = e

    // Or test all group axioms at once:
    law.AbelianGroup(add, neg, zero)
}
```

### Functional Options (New!)

Configure tests with options instead of separate functions:

```go
lawtest.For(t, gen, eq).With(lawtest.WithTrials(500)).Associative(add)
```

### Return Values

All methods return `bool` for early exit:

```go
if !law.Associative(add) {
    return // Stop if physics break
}
```

## Why?

Traditional unit tests check specific inputs. Property tests check **rules that should always hold**:

```go
// Unit test: Does 2 + 3 = 5?
assert.Equal(t, 5, add(2, 3))

// Property test: Is addition associative for ALL numbers?
lawtest.For(t, intGen, intEq).Associative(add)
```

The property test generates hundreds of random inputs and verifies the mathematical law holds.

## Installation

```bash
go get github.com/alexshd/lawtest
```

## Quick Start

### Simple Types (Fluent API)

```go
func TestIntegerAddition(t *testing.T) {
    gen := func() int { return rand.Intn(200) - 100 }
    eq := func(a, b int) bool { return a == b }
    add := func(a, b int) int { return a + b }
    neg := func(a int) int { return -a }

    law := lawtest.For(t, gen, eq)
    law.Group(add, neg, 0)  // Tests associativity, identity, inverse
}
```

### Vector Spaces

```go
func TestZ3VectorSpace(t *testing.T) {
    vs := lawtest.ForVectorSpace(t, vecGen, scalGen, vecEq)
    vs.Axioms(add, zero, neg, scale, scalarAdd, scalarMul, scalarOne)
}
```

## Properties Supported

### Core Properties (via `Laws[T]`)

- **Associative**: `(a ∘ b) ∘ c = a ∘ (b ∘ c)`
- **Commutative**: `a ∘ b = b ∘ a`
- **Identity**: `a ∘ e = a` (e is identity element)
- **Inverse**: `a ∘ a⁻¹ = e` (inverse exists)
- **Idempotent**: `f(f(x)) = f(x)`
- **Group**: All group axioms at once
- **AbelianGroup**: Group + commutativity

### Concurrency

- **Immutable**: Operation doesn't mutate inputs
- **ParallelSafe**: No race conditions

### Vector Space (via `VectorLaws[V, S]`)

- **Axioms**: All 8 vector space axioms

### Legacy API (still available)

- **AssociativeCustom**, **CommutativeCustom**, etc. for backward compatibility
- **GroupCustom[T any]** interface for non-comparable types
- **TestGroupCustom**, **TestVectorSpaceCustom** for interface-based testing

## API Reference

### Fluent API (v0.2.0)

```go
// Create a law checker
law := lawtest.For(t, gen, eq)

// Individual checks (all return bool)
law.Associative(op)           // (a·b)·c = a·(b·c)
law.Commutative(op)           // a·b = b·a
law.Identity(op, e)           // a·e = a
law.Inverse(op, inv, e)       // a·a⁻¹ = e
law.Idempotent(op)            // f(f(x)) = f(x)

// Convenience methods
law.Group(op, inv, e)         // All group axioms
law.AbelianGroup(op, inv, e)  // Group + commutativity

// Concurrency tests
law.Immutable(op)             // op doesn't mutate inputs
law.ParallelSafe(op, 8)       // op is thread-safe

// With options
law.With(lawtest.WithTrials(500)).Associative(op)
law.With(lawtest.WithTimeout(time.Second)).Group(op, inv, e)
```

### Vector Space API

```go
vs := lawtest.ForVectorSpace(t, vecGen, scalGen, vecEq)
vs.Axioms(add, zero, neg, scale, sAdd, sMul, sOne)
```

### Options

| Option           | Description                                |
| ---------------- | ------------------------------------------ |
| `WithTrials(n)`  | Number of random test cases (default: 100) |
| `WithTimeout(d)` | Timeout for all trials (default: 30s)      |

### Legacy API

For backward compatibility, these are still available:

- `Associative`, `AssociativeCustom`, `AssociativeWithConfig`
- `Commutative`, `CommutativeCustom`, `CommutativeWithConfig`
- `Identity`, `IdentityCustom`
- `Inverse`, `InverseCustom`
- `TestGroupCustom`, `TestAbelianGroupCustom`, `TestVectorSpaceCustom`

## Requirements

- Go 1.18 or higher (uses generics)
- [examples/interface_groups/CONCURRENCY_TESTING.md](examples/interface_groups/CONCURRENCY_TESTING.md) - Concurrency testing guide

## License

Apache 2.0 - See [LICENSE](LICENSE) file for details.

Copyright 2025 Alex Shadrin ([@alexshd](https://github.com/alexshd))
