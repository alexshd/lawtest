package lawtest

import (
	"testing"
	"time"
)

// ===========================================================================
// FLUENT API TESTS
// ===========================================================================

func TestFluentAPI_IntegerAddition(t *testing.T) {
	gen := func() int { return IntGen(-100, 100)() }
	eq := func(a, b int) bool { return a == b }
	add := func(a, b int) int { return a + b }

	law := For(t, gen, eq)
	law.Associative(add)
	law.Commutative(add)
}

func TestFluentAPI_Group(t *testing.T) {
	gen := func() int { return IntGen(-100, 100)() }
	eq := func(a, b int) bool { return a == b }
	add := func(a, b int) int { return a + b }
	neg := func(a int) int { return -a }

	law := For(t, gen, eq)
	law.Group(add, neg, 0)
}

func TestFluentAPI_AbelianGroup(t *testing.T) {
	gen := func() int { return IntGen(-100, 100)() }
	eq := func(a, b int) bool { return a == b }
	add := func(a, b int) int { return a + b }
	neg := func(a int) int { return -a }

	law := For(t, gen, eq)
	law.AbelianGroup(add, neg, 0)
}

func TestFluentAPI_WithOptions(t *testing.T) {
	gen := func() int { return IntGen(-100, 100)() }
	eq := func(a, b int) bool { return a == b }
	add := func(a, b int) int { return a + b }
	neg := func(a int) int { return -a }

	law := For(t, gen, eq).With(WithTrials(50))
	law.AbelianGroup(add, neg, 0)
}

func TestFluentAPI_Idempotent(t *testing.T) {
	gen := func() int { return IntGen(-100, 100)() }
	eq := func(a, b int) bool { return a == b }
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}

	law := For(t, gen, eq)
	law.Idempotent(abs)
}

func TestFluentAPI_Immutable(t *testing.T) {
	gen := func() int { return IntGen(-100, 100)() }
	eq := func(a, b int) bool { return a == b }
	add := func(a, b int) int { return a + b }

	law := For(t, gen, eq)
	law.Immutable(add)
}

func TestFluentAPI_ParallelSafe(t *testing.T) {
	gen := func() int { return IntGen(-100, 100)() }
	eq := func(a, b int) bool { return a == b }
	add := func(a, b int) int { return a + b }

	law := For(t, gen, eq)
	law.ParallelSafe(add, 10)
}

// ===========================================================================
// VECTOR SPACE TESTS
// ===========================================================================

// Z3 scalar type for vector space tests
type z3 int8

const (
	z3Zero z3 = 0
	z3One  z3 = 1
	z3Two  z3 = 2
)

func z3Add(a, b z3) z3 { return z3((int(a) + int(b)) % 3) }
func z3Mul(a, b z3) z3 { return z3((int(a) * int(b)) % 3) }
func z3Neg(a z3) z3    { return z3((3 - int(a)) % 3) }

// Simple 2D Z3 vector
type vec2 struct{ x, y z3 }

func vec2Add(a, b vec2) vec2 { return vec2{z3Add(a.x, b.x), z3Add(a.y, b.y)} }
func vec2Neg(a vec2) vec2    { return vec2{z3Neg(a.x), z3Neg(a.y)} }
func vec2Scale(s z3, v vec2) vec2 {
	return vec2{z3Mul(s, v.x), z3Mul(s, v.y)}
}
func vec2Eq(a, b vec2) bool { return a.x == b.x && a.y == b.y }
func vec2Gen() vec2         { return vec2{z3(IntGen(0, 2)()), z3(IntGen(0, 2)())} }
func z3Gen() z3             { return z3(IntGen(0, 2)()) }

func TestFluentAPI_VectorSpace(t *testing.T) {
	vs := ForVectorSpace(t, vec2Gen, z3Gen, vec2Eq)
	vs.Axioms(
		vec2Add,
		vec2{z3Zero, z3Zero},
		vec2Neg,
		vec2Scale,
		z3Add,
		z3Mul,
		z3One,
	)
}

// ===========================================================================
// OPTIONS TESTS
// ===========================================================================

func TestOptions_WithTrials(t *testing.T) {
	gen := func() int { return IntGen(-100, 100)() }
	eq := func(a, b int) bool { return a == b }
	add := func(a, b int) int { return a + b }

	// Low trials
	law := For(t, gen, eq).With(WithTrials(10))
	law.Associative(add)

	// High trials
	law2 := For(t, gen, eq).With(WithTrials(500))
	law2.Associative(add)
}

func TestOptions_VectorSpaceWith(t *testing.T) {
	vs := ForVectorSpace(t, vec2Gen, z3Gen, vec2Eq).With(WithTrials(50))
	vs.Axioms(
		vec2Add,
		vec2{z3Zero, z3Zero},
		vec2Neg,
		vec2Scale,
		z3Add,
		z3Mul,
		z3One,
	)
}

func TestOptions_WithTimeout(t *testing.T) {
	gen := func() int { return IntGen(-100, 100)() }
	eq := func(a, b int) bool { return a == b }
	add := func(a, b int) int { return a + b }

	law := For(t, gen, eq).With(WithTimeout(10 * time.Second))
	law.Associative(add)
}

// ===========================================================================
// GENERATOR TESTS
// ===========================================================================

func TestIntGen(t *testing.T) {
	gen := IntGen(-10, 10)
	for i := 0; i < 100; i++ {
		v := gen()
		if v < -10 || v > 10 {
			t.Errorf("IntGen out of range: %d", v)
		}
	}
}

func TestStringGen(t *testing.T) {
	gen := StringGen(8)
	for i := 0; i < 10; i++ {
		s := gen()
		if len(s) != 8 {
			t.Errorf("StringGen wrong length: %d", len(s))
		}
	}
}

func TestFloat64Gen(t *testing.T) {
	gen := Float64Gen(0.0, 1.0)
	for i := 0; i < 100; i++ {
		v := gen()
		if v < 0.0 || v > 1.0 {
			t.Errorf("Float64Gen out of range: %f", v)
		}
	}
}

func TestBoolGen(t *testing.T) {
	gen := BoolGen()
	trueCount := 0
	for i := 0; i < 100; i++ {
		if gen() {
			trueCount++
		}
	}
	// Should have some of each
	if trueCount == 0 || trueCount == 100 {
		t.Errorf("BoolGen not random: %d trues", trueCount)
	}
}

// ===========================================================================
// EQUIVALENCE TESTS
// ===========================================================================

func TestEquivalent(t *testing.T) {
	// Two ways to compute absolute value
	abs1 := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}
	abs2 := func(x int) int {
		if x >= 0 {
			return x
		}
		return -x
	}

	gen := func() int { return IntGen(-100, 100)() }
	eq := func(a, b int) bool { return a == b }

	Equivalent(t, abs1, abs2, gen, eq)
}
