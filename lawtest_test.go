package lawtest_test

import (
	"testing"

	"github.com/alexshd/lawtest"
)

// Example: Testing integer addition properties
func TestIntAddition(t *testing.T) {
	add := func(a, b int) int { return a + b }
	gen := lawtest.IntGen(-100, 100)

	t.Run("Associative", func(t *testing.T) {
		lawtest.Associative(t, add, gen)
	})

	t.Run("Commutative", func(t *testing.T) {
		lawtest.Commutative(t, add, gen)
	})

	t.Run("Identity", func(t *testing.T) {
		lawtest.Identity(t, add, 0, gen)
	})

	t.Run("Closure", func(t *testing.T) {
		lawtest.Closure(t, add, gen)
	})
}

// Example: Testing integer multiplication properties
func TestIntMultiplication(t *testing.T) {
	mul := func(a, b int) int { return a * b }
	gen := lawtest.IntGen(-50, 50)

	t.Run("Associative", func(t *testing.T) {
		lawtest.Associative(t, mul, gen)
	})

	t.Run("Commutative", func(t *testing.T) {
		lawtest.Commutative(t, mul, gen)
	})

	t.Run("Identity", func(t *testing.T) {
		lawtest.Identity(t, mul, 1, gen)
	})
}

// Example: Testing string concatenation
func TestStringConcat(t *testing.T) {
	concat := func(a, b string) string { return a + b }
	gen := lawtest.StringGen(5)

	t.Run("Associative", func(t *testing.T) {
		lawtest.Associative(t, concat, gen)
	})

	t.Run("Identity", func(t *testing.T) {
		lawtest.Identity(t, concat, "", gen)
	})

	// Note: String concatenation is NOT commutative
	// "abc" + "def" != "def" + "abc"
}

// Example: Testing boolean operations
func TestBooleanOr(t *testing.T) {
	or := func(a, b bool) bool { return a || b }
	gen := lawtest.BoolGen()

	t.Run("Associative", func(t *testing.T) {
		lawtest.Associative(t, or, gen)
	})

	t.Run("Commutative", func(t *testing.T) {
		lawtest.Commutative(t, or, gen)
	})

	t.Run("Identity", func(t *testing.T) {
		lawtest.Identity(t, or, false, gen)
	})

	t.Run("Idempotent", func(t *testing.T) {
		idempOp := func(b bool) bool { return or(b, b) }
		lawtest.Idempotent(t, idempOp, gen)
	})
}

// Example: Custom generator for complex types
type Point struct {
	X, Y int
}

func PointGen() lawtest.Generator[Point] {
	intGen := lawtest.IntGen(-100, 100)
	return func() Point {
		return Point{X: intGen(), Y: intGen()}
	}
}

func TestPointOperations(t *testing.T) {
	// Vector addition
	add := func(a, b Point) Point {
		return Point{X: a.X + b.X, Y: a.Y + b.Y}
	}

	gen := PointGen()

	t.Run("Associative", func(t *testing.T) {
		lawtest.Associative(t, add, gen)
	})

	t.Run("Commutative", func(t *testing.T) {
		lawtest.Commutative(t, add, gen)
	})

	t.Run("Identity", func(t *testing.T) {
		origin := Point{X: 0, Y: 0}
		lawtest.Identity(t, add, origin, gen)
	})

	t.Run("Inverse", func(t *testing.T) {
		negate := func(p Point) Point {
			return Point{X: -p.X, Y: -p.Y}
		}
		origin := Point{X: 0, Y: 0}
		lawtest.Inverse(t, add, negate, origin, gen)
	})
}
