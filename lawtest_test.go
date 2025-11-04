package lawtest_test

import (
	"math/rand"
	"testing"

	"github.com/alexshd/lawtest"
)

// Testing integer addition properties
func TestIntAddition(t *testing.T) {
	addOp := func(a, b int) int { return a + b }
	intGen := lawtest.IntGen(-100, 100)

	t.Run("Associative", func(t *testing.T) {
		lawtest.Associative(t, addOp, intGen)
	})

	t.Run("Commutative", func(t *testing.T) {
		lawtest.Commutative(t, addOp, intGen)
	})

	t.Run("Identity", func(t *testing.T) {
		lawtest.Identity(t, addOp, 0, intGen)
	})

	t.Run("Closure", func(t *testing.T) {
		lawtest.Closure(t, addOp, intGen)
	})

	t.Run("Inverse", func(t *testing.T) {
		negateOp := func(x int) int { return -x }
		lawtest.Inverse(t, addOp, negateOp, 0, intGen)
	})
}

// Testing integer multiplication properties
func TestIntMultiplication(t *testing.T) {
	mulOp := func(a, b int) int { return a * b }
	intGen := lawtest.IntGen(-50, 50)

	t.Run("Associative", func(t *testing.T) {
		lawtest.Associative(t, mulOp, intGen)
	})

	t.Run("Commutative", func(t *testing.T) {
		lawtest.Commutative(t, mulOp, intGen)
	})

	t.Run("Identity", func(t *testing.T) {
		lawtest.Identity(t, mulOp, 1, intGen)
	})
}

// Testing string concatenation
func TestStringConcat(t *testing.T) {
	concatOp := func(a, b string) string { return a + b }
	strGen := lawtest.StringGen(5)

	t.Run("Associative", func(t *testing.T) {
		lawtest.Associative(t, concatOp, strGen)
	})

	t.Run("Identity", func(t *testing.T) {
		lawtest.Identity(t, concatOp, "", strGen)
	})
}

// Testing boolean operations
func TestBooleanOr(t *testing.T) {
	orOp := func(a, b bool) bool { return a || b }
	boolGen := lawtest.BoolGen()

	t.Run("Associative", func(t *testing.T) {
		lawtest.Associative(t, orOp, boolGen)
	})

	t.Run("Commutative", func(t *testing.T) {
		lawtest.Commutative(t, orOp, boolGen)
	})

	t.Run("Identity", func(t *testing.T) {
		lawtest.Identity(t, orOp, false, boolGen)
	})

	t.Run("Idempotent", func(t *testing.T) {
		idempOp := func(b bool) bool { return orOp(b, b) }
		lawtest.Idempotent(t, idempOp, boolGen)
	})
}

// Testing Float64Gen generator
func TestFloat64Operations(t *testing.T) {
	mulOp := func(a, b float64) float64 { return a * b }
	floatGen := lawtest.Float64Gen(1.0, 10.0)

	t.Run("Commutative", func(t *testing.T) {
		lawtest.Commutative(t, mulOp, floatGen)
	})
}

// Custom generator for complex types
type Point struct {
	X, Y int
}

func TestPointOperations(t *testing.T) {
	pointGen := func() Point {
		intGen := lawtest.IntGen(-100, 100)
		return Point{X: intGen(), Y: intGen()}
	}

	addOp := func(a, b Point) Point {
		return Point{X: a.X + b.X, Y: a.Y + b.Y}
	}

	t.Run("Associative", func(t *testing.T) {
		lawtest.Associative(t, addOp, pointGen)
	})

	t.Run("Commutative", func(t *testing.T) {
		lawtest.Commutative(t, addOp, pointGen)
	})

	t.Run("Identity", func(t *testing.T) {
		origin := Point{X: 0, Y: 0}
		lawtest.Identity(t, addOp, origin, pointGen)
	})

	t.Run("Inverse", func(t *testing.T) {
		negateOp := func(p Point) Point {
			return Point{X: -p.X, Y: -p.Y}
		}
		origin := Point{X: 0, Y: 0}
		lawtest.Inverse(t, addOp, negateOp, origin, pointGen)
	})
}

// Test interface-based Group implementation
type IntModGroup struct {
	modulus int
}

func (g IntModGroup) Op(a, b int) int {
	return (a + b) % g.modulus
}

func (g IntModGroup) Identity() int {
	return 0
}

func (g IntModGroup) Inverse(a int) int {
	return (g.modulus - a) % g.modulus
}

func (g IntModGroup) Gen() int {
	return rand.Intn(g.modulus)
}

func TestIntModGroup(t *testing.T) {
	modGroup := IntModGroup{modulus: 12}
	lawtest.TestGroup[int](t, modGroup)
}

// Test interface-based Monoid implementation
type StringConcatMonoid struct{}

func (m StringConcatMonoid) Op(a, b string) string {
	return a + b
}

func (m StringConcatMonoid) Identity() string {
	return ""
}

func (m StringConcatMonoid) Gen() string {
	return lawtest.StringGen(5)()
}

func TestStringConcatMonoid(t *testing.T) {
	monoid := StringConcatMonoid{}
	lawtest.TestMonoid[string](t, monoid)
}

// Test interface-based Semigroup implementation
type MaxSemigroup struct{}

func (s MaxSemigroup) Op(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (s MaxSemigroup) Gen() int {
	return lawtest.IntGen(1, 1000)()
}

func TestMaxSemigroup(t *testing.T) {
	semigroup := MaxSemigroup{}
	lawtest.TestSemigroup[int](t, semigroup)
}

// Test interface-based IdempotentOp implementation
type AbsoluteValue struct{}

func (o AbsoluteValue) Apply(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (o AbsoluteValue) Gen() int {
	return lawtest.IntGen(-100, 100)()
}

func TestAbsoluteValueIdempotent(t *testing.T) {
	absOp := AbsoluteValue{}
	lawtest.TestIdempotentOp[int](t, absOp)
}

// Test Homomorphism - simple valid example
type IntAdditionGroup struct{}

func (g IntAdditionGroup) Op(a, b int) int   { return a + b }
func (g IntAdditionGroup) Identity() int     { return 0 }
func (g IntAdditionGroup) Inverse(a int) int { return -a }
func (g IntAdditionGroup) Gen() int          { return lawtest.IntGen(-50, 50)() }

type DoubleHomomorphism struct{}

func (h DoubleHomomorphism) Map(x int) int {
	return x * 2
}

func (h DoubleHomomorphism) SourceGroup() lawtest.Group[int] {
	return IntAdditionGroup{}
}

func (h DoubleHomomorphism) TargetGroup() lawtest.Group[int] {
	// Target is also addition (2x is homomorphism from + to +)
	return IntAdditionGroup{}
}

func TestDoubleHomomorphism(t *testing.T) {
	homo := DoubleHomomorphism{}
	lawtest.TestHomomorphism[int, int](t, homo)
}

// Test concurrency functions
type SafeCache struct {
	value int
}

func (c *SafeCache) Merge(other *SafeCache) *SafeCache {
	return &SafeCache{value: c.value + other.value}
}

func TestParallelSafety(t *testing.T) {
	mergeOp := func(a, b *SafeCache) *SafeCache {
		return a.Merge(b)
	}
	cacheGen := func() *SafeCache {
		return &SafeCache{value: rand.Intn(100)}
	}

	t.Run("ParallelSafe", func(t *testing.T) {
		isSafe := lawtest.ParallelSafe(t, mergeOp, cacheGen, 10)
		if !isSafe {
			t.Error("Expected cache merge to be parallel-safe")
		}
	})

	t.Run("ParallelAssociativity", func(t *testing.T) {
		lawtest.TestParallelAssociativity(t, mergeOp, cacheGen, 10)
	})

	t.Run("Immutability", func(t *testing.T) {
		lawtest.ImmutableOp(t, mergeOp, cacheGen)
	})
}

// Test custom configuration
func TestWithCustomConfig(t *testing.T) {
	addOp := func(a, b int) int { return a + b }
	intGen := lawtest.IntGen(-1000, 1000)

	cfg := lawtest.DefaultConfig()
	cfg.TestCases = 200

	t.Run("AssociativeWithConfig", func(t *testing.T) {
		lawtest.AssociativeWithConfig(t, addOp, intGen, cfg)
	})

	t.Run("CommutativeWithConfig", func(t *testing.T) {
		lawtest.CommutativeWithConfig(t, addOp, intGen, cfg)
	})

	t.Run("IdentityWithConfig", func(t *testing.T) {
		lawtest.IdentityWithConfig(t, addOp, 0, intGen, cfg)
	})
}

// Test generator edge cases
func TestGeneratorEdgeCases(t *testing.T) {
	t.Run("IntGen_SameMinMax", func(t *testing.T) {
		constGen := lawtest.IntGen(5, 5)
		for i := 0; i < 10; i++ {
			if val := constGen(); val != 5 {
				t.Errorf("Expected 5, got %d", val)
			}
		}
	})

	t.Run("IntGen_Panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for min > max")
			}
		}()
		_ = lawtest.IntGen(10, 5)
	})

	t.Run("Float64Gen_Panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for min > max")
			}
		}()
		_ = lawtest.Float64Gen(10.0, 5.0)
	})
}
