# goprops Quick Start

## Installation

```bash
go get github.com/alexshd/goprops
```

## 30-Second Example

```go
package mycode_test

import (
    "testing"
    "github.com/alexshd/goprops"
)

func TestAdditionIsAssociative(t *testing.T) {
    add := func(a, b int) int { return a + b }
    gen := goprops.IntGen(-1000, 1000)
    
    goprops.Associative(t, add, gen)
}
```

Run: `go test -v`

**What just happened?** 
- Generated 100 random integers
- Tested that `(a + b) + c = a + (b + c)` for all combinations
- If it finds a counterexample, test fails with the exact values

## Five-Minute Tutorial

Test all basic properties of integer addition:

```go
func TestIntegerAddition(t *testing.T) {
    add := func(a, b int) int { return a + b }
    gen := goprops.IntGen(-100, 100)
    
    t.Run("Associative", func(t *testing.T) {
        goprops.Associative(t, add, gen)
    })
    
    t.Run("Commutative", func(t *testing.T) {
        goprops.Commutative(t, add, gen)
    })
    
    t.Run("Identity", func(t *testing.T) {
        goprops.Identity(t, add, 0, gen)
    })
}
```

## Real-World Example: Cache Testing

```go
type Cache struct {
    data map[string]int
}

func (c *Cache) Set(key string, val int) { c.data[key] = val }
func (c *Cache) Get(key string) int { return c.data[key] }

func TestCacheLastWriteWins(t *testing.T) {
    cache := &Cache{data: make(map[string]int)}
    
    for i := 0; i < 100; i++ {
        key := "test"
        a := goprops.IntGen(1, 1000)()
        b := goprops.IntGen(1, 1000)()
        
        cache.Set(key, a)
        cache.Set(key, b)
        
        if cache.Get(key) != b {
            t.Error("Last write should win")
        }
    }
}
```

## Custom Types

```go
type Point struct {
    X, Y int
}

func PointGen() goprops.Generator[Point] {
    intGen := goprops.IntGen(-100, 100)
    return func() Point {
        return Point{X: intGen(), Y: intGen()}
    }
}

func TestVectorAddition(t *testing.T) {
    add := func(a, b Point) Point {
        return Point{X: a.X + b.X, Y: a.Y + b.Y}
    }
    
    goprops.Associative(t, add, PointGen())
    goprops.Commutative(t, add, PointGen())
    goprops.Identity(t, add, Point{0, 0}, PointGen())
}
```

## Available Properties

- `Associative(t, op, gen)` - Tests `(a∘b)∘c = a∘(b∘c)`
- `Commutative(t, op, gen)` - Tests `a∘b = b∘a`
- `Identity(t, op, id, gen)` - Tests `a∘e = a`
- `Inverse(t, op, inv, id, gen)` - Tests `a∘a⁻¹ = e`
- `Closure(t, op, gen)` - Tests result type matches input
- `Idempotent(t, op, gen)` - Tests `f(f(x)) = f(x)`

## Built-in Generators

- `IntGen(min, max)` - Random integers in range
- `Float64Gen(min, max)` - Random floats in range
- `StringGen(length)` - Random strings
- `BoolGen()` - Random booleans

## What This Catches

**Unit tests miss:**
- Edge cases you didn't think of
- Structural properties of operations
- Bugs that only appear with specific value combinations

**Property tests catch:**
- plyGO's bug where all slices map to same key ✓
- Cache consistency issues ✓
- Merge operations that aren't associative ✓
- Operations that claim to be commutative but aren't ✓

## Next Steps

1. Read `LEARNING_GUIDE.md` for deep dive
2. Check `examples/` for advanced patterns
3. Run your own code through property tests
4. Watch bugs you didn't know existed get caught

**Philosophy:** One mathematical law replaces hundreds of example tests.
