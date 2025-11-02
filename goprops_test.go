package goprops_test

import (
"testing"

"github.com/alexshd/goprops"
)

// Example: Testing integer addition properties
func TestIntAddition(t *testing.T) {
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

	t.Run("Closure", func(t *testing.T) {
goprops.Closure(t, add, gen)
})
}

// Example: Testing integer multiplication properties
func TestIntMultiplication(t *testing.T) {
	mul := func(a, b int) int { return a * b }
	gen := goprops.IntGen(-50, 50)

	t.Run("Associative", func(t *testing.T) {
goprops.Associative(t, mul, gen)
})

	t.Run("Commutative", func(t *testing.T) {
goprops.Commutative(t, mul, gen)
})

	t.Run("Identity", func(t *testing.T) {
goprops.Identity(t, mul, 1, gen)
})
}

// Example: Testing string concatenation
func TestStringConcat(t *testing.T) {
	concat := func(a, b string) string { return a + b }
	gen := goprops.StringGen(5)

	t.Run("Associative", func(t *testing.T) {
goprops.Associative(t, concat, gen)
})

	t.Run("Identity", func(t *testing.T) {
goprops.Identity(t, concat, "", gen)
})

	// Note: String concatenation is NOT commutative
	// "abc" + "def" != "def" + "abc"
}

// Example: Testing boolean operations
func TestBooleanOr(t *testing.T) {
	or := func(a, b bool) bool { return a || b }
	gen := goprops.BoolGen()

	t.Run("Associative", func(t *testing.T) {
goprops.Associative(t, or, gen)
})

	t.Run("Commutative", func(t *testing.T) {
goprops.Commutative(t, or, gen)
})

	t.Run("Identity", func(t *testing.T) {
goprops.Identity(t, or, false, gen)
})

	t.Run("Idempotent", func(t *testing.T) {
idempOp := func(b bool) bool { return or(b, b) }
goprops.Idempotent(t, idempOp, gen)
})
}

// Example: Custom generator for complex types
type Point struct {
	X, Y int
}

func PointGen() goprops.Generator[Point] {
	intGen := goprops.IntGen(-100, 100)
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
goprops.Associative(t, add, gen)
})

	t.Run("Commutative", func(t *testing.T) {
goprops.Commutative(t, add, gen)
})

	t.Run("Identity", func(t *testing.T) {
origin := Point{X: 0, Y: 0}
goprops.Identity(t, add, origin, gen)
})

	t.Run("Inverse", func(t *testing.T) {
negate := func(p Point) Point {
return Point{X: -p.X, Y: -p.Y}
}
origin := Point{X: 0, Y: 0}
goprops.Inverse(t, add, negate, origin, gen)
})
}
