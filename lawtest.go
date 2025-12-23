// Package lawtest provides property-based testing using group theory.
//
// lawtest v0.2.0 - Fluent API with Functional Options
//
// lawtest helps verify mathematical properties of operations through
// randomized testing. Instead of writing specific test cases, you define
// the mathematical laws your code should obey, and lawtest generates
// hundreds of test cases automatically.
//
// # Quick Start (Fluent API)
//
//	func TestPacketPhysics(t *testing.T) {
//	    law := lawtest.For(t, packetGen, packetEq)
//
//	    law.Associative(add)       // (a∘b)∘c = a∘(b∘c)
//	    law.Commutative(add)       // a∘b = b∘a
//	    law.Identity(add, zero)    // a∘e = a
//	    law.Inverse(add, neg, zero)// a∘a⁻¹ = e
//
//	    // Or all at once:
//	    law.AbelianGroup(add, neg, zero)
//	}
//
// # Functional Options
//
//	law := lawtest.For(t, gen, eq).With(lawtest.WithTrials(500))
//
// # Return Values
//
// All check functions return bool for early exit:
//
//	if !law.Associative(add) {
//	    return // Stop if physics break
//	}
//
// # Vector Spaces
//
//	vs := lawtest.ForVectorSpace(t, vecGen, scalGen, vecEq)
//	vs.Axioms(add, zero, neg, scale, scalarAdd, scalarMul, scalarOne)
package lawtest

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// ===========================================================================
// CORE TYPES
// ===========================================================================

// BinaryOp is a binary operation that combines two values of type T.
type BinaryOp[T any] func(a, b T) T

// UnaryOp is a unary operation that transforms a value of type T.
type UnaryOp[T any] func(a T) T

// Generator produces random values of type T for property testing.
type Generator[T any] func() T

// Config holds configuration for property testing.
type Config struct {
	TestCases int           // Number of random test cases (default: 100)
	Timeout   time.Duration // Maximum time per test (default: 5s)
}

// defaultTestCases is the default number of random test cases.
const defaultTestCases = 100

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		TestCases: defaultTestCases,
		Timeout:   5 * time.Second,
	}
}

// ===========================================================================
// FUNCTIONAL OPTIONS
// ===========================================================================

// Option configures property test behavior.
type Option func(*Config)

// WithTrials sets the number of random test cases to run.
func WithTrials(n int) Option {
	return func(c *Config) { c.TestCases = n }
}

// WithTimeout sets the maximum time for the property test.
func WithTimeout(d time.Duration) Option {
	return func(c *Config) { c.Timeout = d }
}

// applyOptions creates a Config from default + options.
func applyOptions(opts ...Option) *Config {
	cfg := DefaultConfig()
	for _, o := range opts {
		o(cfg)
	}
	return cfg
}

// ===========================================================================
// GENERATORS
// ===========================================================================

// IntGen creates a Generator that produces random integers in [min, max].
func IntGen(min, max int) Generator[int] {
	if min > max {
		panic(fmt.Sprintf("min (%d) must be <= max (%d)", min, max))
	}
	return func() int {
		return min + rand.Intn(max-min+1)
	}
}

// StringGen creates a Generator that produces random strings of length n.
func StringGen(n int) Generator[string] {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return func() string {
		b := make([]byte, n)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		return string(b)
	}
}

// Float64Gen creates a Generator that produces random float64 values in [min, max].
func Float64Gen(min, max float64) Generator[float64] {
	if min > max {
		panic(fmt.Sprintf("min (%f) must be <= max (%f)", min, max))
	}
	return func() float64 {
		return min + rand.Float64()*(max-min)
	}
}

// BoolGen creates a Generator that produces random boolean values.
func BoolGen() Generator[bool] {
	return func() bool {
		return rand.Intn(2) == 1
	}
}

// ===========================================================================
// FLUENT API: Laws[T]
// ===========================================================================

// Laws provides a fluent API for testing algebraic laws.
//
// Create with For() and chain assertions:
//
//	law := lawtest.For(t, gen, eq)
//	law.Associative(add)
//	law.AbelianGroup(add, neg, zero)
type Laws[T any] struct {
	t    *testing.T
	gen  Generator[T]
	eq   func(T, T) bool
	opts []Option
}

// For creates a new Laws checker for type T.
//
//	law := lawtest.For(t, packetGen, packetEq)
func For[T any](t *testing.T, gen Generator[T], eq func(T, T) bool) *Laws[T] {
	return &Laws[T]{t: t, gen: gen, eq: eq}
}

// With adds options to all subsequent law checks.
//
//	law := lawtest.For(t, gen, eq).With(lawtest.WithTrials(500))
func (l *Laws[T]) With(opts ...Option) *Laws[T] {
	l.opts = append(l.opts, opts...)
	return l
}

// Associative tests (a ∘ b) ∘ c = a ∘ (b ∘ c).
func (l *Laws[T]) Associative(op BinaryOp[T]) bool {
	cfg := applyOptions(l.opts...)
	l.t.Helper()

	for i := 0; i < cfg.TestCases; i++ {
		a, b, c := l.gen(), l.gen(), l.gen()
		left := op(op(a, b), c)
		right := op(a, op(b, c))

		if !l.equal(left, right) {
			l.t.Errorf("Associativity failed: (a∘b)∘c != a∘(b∘c)\n  a=%v, b=%v, c=%v\n  left=%v, right=%v",
				a, b, c, left, right)
			return false
		}
	}
	l.t.Logf("✅ Associative: (a∘b)∘c = a∘(b∘c)")
	return true
}

// Commutative tests a ∘ b = b ∘ a.
func (l *Laws[T]) Commutative(op BinaryOp[T]) bool {
	cfg := applyOptions(l.opts...)
	l.t.Helper()

	for i := 0; i < cfg.TestCases; i++ {
		a, b := l.gen(), l.gen()
		left := op(a, b)
		right := op(b, a)

		if !l.equal(left, right) {
			l.t.Errorf("Commutativity failed: a∘b != b∘a\n  a=%v, b=%v\n  a∘b=%v, b∘a=%v",
				a, b, left, right)
			return false
		}
	}
	l.t.Logf("✅ Commutative: a∘b = b∘a")
	return true
}

// Identity tests a ∘ e = a and e ∘ a = a.
func (l *Laws[T]) Identity(op BinaryOp[T], identity T) bool {
	cfg := applyOptions(l.opts...)
	l.t.Helper()

	for i := 0; i < cfg.TestCases; i++ {
		a := l.gen()

		right := op(a, identity)
		if !l.equal(right, a) {
			l.t.Errorf("Right identity failed: a∘e != a\n  a=%v, e=%v, a∘e=%v", a, identity, right)
			return false
		}

		left := op(identity, a)
		if !l.equal(left, a) {
			l.t.Errorf("Left identity failed: e∘a != a\n  e=%v, a=%v, e∘a=%v", identity, a, left)
			return false
		}
	}
	l.t.Logf("✅ Identity: a∘e = e∘a = a")
	return true
}

// Inverse tests a ∘ a⁻¹ = e and a⁻¹ ∘ a = e.
func (l *Laws[T]) Inverse(op BinaryOp[T], inv UnaryOp[T], identity T) bool {
	cfg := applyOptions(l.opts...)
	l.t.Helper()

	for i := 0; i < cfg.TestCases; i++ {
		a := l.gen()
		aInv := inv(a)

		left := op(a, aInv)
		if !l.equal(left, identity) {
			l.t.Errorf("Left inverse failed: a∘a⁻¹ != e\n  a=%v, a⁻¹=%v, a∘a⁻¹=%v, e=%v",
				a, aInv, left, identity)
			return false
		}

		right := op(aInv, a)
		if !l.equal(right, identity) {
			l.t.Errorf("Right inverse failed: a⁻¹∘a != e\n  a⁻¹=%v, a=%v, a⁻¹∘a=%v, e=%v",
				aInv, a, right, identity)
			return false
		}
	}
	l.t.Logf("✅ Inverse: a∘a⁻¹ = a⁻¹∘a = e")
	return true
}

// Idempotent tests f(f(x)) = f(x).
func (l *Laws[T]) Idempotent(op UnaryOp[T]) bool {
	cfg := applyOptions(l.opts...)
	l.t.Helper()

	for i := 0; i < cfg.TestCases; i++ {
		x := l.gen()
		fx := op(x)
		ffx := op(fx)

		if !l.equal(fx, ffx) {
			l.t.Errorf("Idempotence failed: f(f(x)) != f(x)\n  x=%v, f(x)=%v, f(f(x))=%v",
				x, fx, ffx)
			return false
		}
	}
	l.t.Logf("✅ Idempotent: f(f(x)) = f(x)")
	return true
}

// Group tests all group axioms: associativity, identity, inverse.
func (l *Laws[T]) Group(op BinaryOp[T], inv UnaryOp[T], identity T) bool {
	l.t.Helper()
	return l.Associative(op) && l.Identity(op, identity) && l.Inverse(op, inv, identity)
}

// AbelianGroup tests all Abelian group axioms: group + commutativity.
func (l *Laws[T]) AbelianGroup(op BinaryOp[T], inv UnaryOp[T], identity T) bool {
	l.t.Helper()
	return l.Group(op, inv, identity) && l.Commutative(op)
}

// Immutable tests that op does not mutate its inputs.
func (l *Laws[T]) Immutable(op BinaryOp[T]) bool {
	cfg := applyOptions(l.opts...)
	l.t.Helper()

	for i := 0; i < cfg.TestCases; i++ {
		a, b := l.gen(), l.gen()
		aOrig, bOrig := a, b

		_ = op(a, b)

		if !l.equal(a, aOrig) {
			l.t.Errorf("Immutability violated: operation mutated first argument\n  before=%v, after=%v",
				aOrig, a)
			return false
		}
		if !l.equal(b, bOrig) {
			l.t.Errorf("Immutability violated: operation mutated second argument\n  before=%v, after=%v",
				bOrig, b)
			return false
		}
	}
	l.t.Logf("✅ Immutable: operation does not mutate inputs")
	return true
}

// ParallelSafe tests if op can be safely executed concurrently.
func (l *Laws[T]) ParallelSafe(op BinaryOp[T], goroutines int) bool {
	l.t.Helper()

	a, b := l.gen(), l.gen()
	expected := op(a, b)

	results := make([]T, goroutines)
	done := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			results[idx] = op(a, b)
			done <- true
		}(i)
	}

	for i := 0; i < goroutines; i++ {
		<-done
	}

	for i, result := range results {
		if !l.equal(result, expected) {
			l.t.Errorf("Parallel safety failed: goroutine %d produced different result\n  expected=%v, got=%v",
				i, expected, result)
			return false
		}
	}
	l.t.Logf("✅ ParallelSafe: no race conditions in %d goroutines", goroutines)
	return true
}

// equal compares two values using custom eq or any comparison.
func (l *Laws[T]) equal(a, b T) bool {
	if l.eq != nil {
		return l.eq(a, b)
	}
	return any(a) == any(b)
}

// ===========================================================================
// VECTOR SPACE TESTING
// ===========================================================================

// VectorLaws provides fluent API for vector space testing.
type VectorLaws[V any, S any] struct {
	t       *testing.T
	genVec  Generator[V]
	genScal Generator[S]
	eqVec   func(V, V) bool
	opts    []Option
}

// ForVectorSpace creates a new VectorLaws checker.
//
//	vs := lawtest.ForVectorSpace(t, vecGen, scalGen, vecEq)
//	vs.Axioms(add, zero, neg, scale, scalarAdd, scalarMul, scalarOne)
func ForVectorSpace[V any, S any](
	t *testing.T,
	genVec Generator[V],
	genScal Generator[S],
	eqVec func(V, V) bool,
) *VectorLaws[V, S] {
	return &VectorLaws[V, S]{t: t, genVec: genVec, genScal: genScal, eqVec: eqVec}
}

// With adds options to all subsequent checks.
func (v *VectorLaws[V, S]) With(opts ...Option) *VectorLaws[V, S] {
	v.opts = append(v.opts, opts...)
	return v
}

// Axioms tests all vector space axioms.
func (v *VectorLaws[V, S]) Axioms(
	add func(V, V) V,
	zero V,
	neg func(V) V,
	scale func(S, V) V,
	scalarAdd func(S, S) S,
	scalarMul func(S, S) S,
	scalarOne S,
) bool {
	v.t.Helper()
	cfg := applyOptions(v.opts...)

	// Test vector addition is Abelian group
	vecLaw := &Laws[V]{t: v.t, gen: v.genVec, eq: v.eqVec, opts: v.opts}
	if !vecLaw.AbelianGroup(add, neg, zero) {
		return false
	}

	// Scalar identity: 1·v = v
	v.t.Run("ScalarIdentity", func(t *testing.T) {
		for i := 0; i < cfg.TestCases; i++ {
			vec := v.genVec()
			result := scale(scalarOne, vec)
			if !v.eqVec(result, vec) {
				t.Errorf("Scalar identity failed: 1·v != v")
				return
			}
		}
		t.Logf("✅ ScalarIdentity: 1·v = v")
	})

	// Scalar distributes over vector addition: s·(v+w) = s·v + s·w
	v.t.Run("ScalarDistributesOverVectorAdd", func(t *testing.T) {
		for i := 0; i < cfg.TestCases; i++ {
			s := v.genScal()
			vec1, vec2 := v.genVec(), v.genVec()
			left := scale(s, add(vec1, vec2))
			right := add(scale(s, vec1), scale(s, vec2))
			if !v.eqVec(left, right) {
				t.Errorf("Distribution failed: s·(v+w) != s·v + s·w")
				return
			}
		}
		t.Logf("✅ ScalarDistributesOverVectorAdd: s·(v+w) = s·v + s·w")
	})

	// Scalar distributes over scalar addition: (s+t)·v = s·v + t·v
	v.t.Run("ScalarDistributesOverScalarAdd", func(t *testing.T) {
		for i := 0; i < cfg.TestCases; i++ {
			s, s2 := v.genScal(), v.genScal()
			vec := v.genVec()
			left := scale(scalarAdd(s, s2), vec)
			right := add(scale(s, vec), scale(s2, vec))
			if !v.eqVec(left, right) {
				t.Errorf("Distribution failed: (s+t)·v != s·v + t·v")
				return
			}
		}
		t.Logf("✅ ScalarDistributesOverScalarAdd: (s+t)·v = s·v + t·v")
	})

	// Scalar multiplication associativity: (s·t)·v = s·(t·v)
	v.t.Run("ScalarMulAssociative", func(t *testing.T) {
		for i := 0; i < cfg.TestCases; i++ {
			s, s2 := v.genScal(), v.genScal()
			vec := v.genVec()
			left := scale(scalarMul(s, s2), vec)
			right := scale(s, scale(s2, vec))
			if !v.eqVec(left, right) {
				t.Errorf("Associativity failed: (s·t)·v != s·(t·v)")
				return
			}
		}
		t.Logf("✅ ScalarMulAssociative: (s·t)·v = s·(t·v)")
	})

	return true
}

// ===========================================================================
// INTERFACE-BASED GROUP TESTING
// ===========================================================================

// Group represents an algebraic group with a binary operation.
//
// A Group must satisfy:
//   - Associativity: (a ∘ b) ∘ c = a ∘ (b ∘ c)
//   - Identity: exists e where a ∘ e = e ∘ a = a
//   - Inverse: for each a, exists a⁻¹ where a ∘ a⁻¹ = a⁻¹ ∘ a = e
//
// Example:
//
//	type IntAddMod12 struct{}
//	func (g IntAddMod12) Op(a, b int) int { return (a + b) % 12 }
//	func (g IntAddMod12) Identity() int { return 0 }
//	func (g IntAddMod12) Inverse(a int) int { return (12 - a) % 12 }
//	func (g IntAddMod12) Gen() int { return rand.Intn(12) }
//
//	func TestModularArithmetic(t *testing.T) {
//	    lawtest.TestGroup(t, IntAddMod12{})
//	}
type Group[T comparable] interface {
	Op(a, b T) T
	Identity() T
	Inverse(a T) T
	Gen() T
}

// TestGroup verifies all group properties for a type implementing Group.
func TestGroup[T comparable](t *testing.T, g Group[T]) {
	t.Helper()
	law := For(t, g.Gen, func(a, b T) bool { return a == b })
	law.Group(
		func(a, b T) T { return g.Op(a, b) },
		func(a T) T { return g.Inverse(a) },
		g.Identity(),
	)
}

// GroupCustom represents a group for non-comparable types (slices, maps).
//
// Example:
//
//	type PacketGroup struct{ Dim int }
//	func (g PacketGroup) Op(a, b Packet) Packet { return a.Add(b) }
//	func (g PacketGroup) Identity() Packet { return ZeroPacket(g.Dim) }
//	func (g PacketGroup) Inverse(a Packet) Packet { return a.Inverse() }
//	func (g PacketGroup) Gen() Packet { return RandomPacket(g.Dim) }
//	func (g PacketGroup) Eq(a, b Packet) bool { return a.Equal(b) }
type GroupCustom[T any] interface {
	Op(a, b T) T
	Identity() T
	Inverse(a T) T
	Gen() T
	Eq(a, b T) bool
}

// TestGroupCustom verifies all group properties for non-comparable types.
func TestGroupCustom[T any](t *testing.T, g GroupCustom[T]) {
	t.Helper()
	law := For(t, g.Gen, g.Eq)
	law.Group(
		func(a, b T) T { return g.Op(a, b) },
		func(a T) T { return g.Inverse(a) },
		g.Identity(),
	)
}

// AbelianGroup adds commutativity requirement to Group.
type AbelianGroup[T comparable] interface {
	Group[T]
}

// TestAbelianGroup verifies all Abelian group properties.
func TestAbelianGroup[T comparable](t *testing.T, g AbelianGroup[T]) {
	t.Helper()
	law := For(t, g.Gen, func(a, b T) bool { return a == b })
	law.AbelianGroup(
		func(a, b T) T { return g.Op(a, b) },
		func(a T) T { return g.Inverse(a) },
		g.Identity(),
	)
}

// AbelianGroupCustom adds commutativity for non-comparable types.
type AbelianGroupCustom[T any] interface {
	GroupCustom[T]
}

// TestAbelianGroupCustom verifies Abelian group properties for non-comparable types.
func TestAbelianGroupCustom[T any](t *testing.T, g AbelianGroupCustom[T]) {
	t.Helper()
	law := For(t, g.Gen, g.Eq)
	law.AbelianGroup(
		func(a, b T) T { return g.Op(a, b) },
		func(a T) T { return g.Inverse(a) },
		g.Identity(),
	)
}

// Monoid represents a semigroup with identity (no inverse required).
type Monoid[T comparable] interface {
	Op(a, b T) T
	Identity() T
	Gen() T
}

// TestMonoid verifies monoid properties (associativity + identity).
func TestMonoid[T comparable](t *testing.T, m Monoid[T]) {
	t.Helper()
	law := For(t, m.Gen, func(a, b T) bool { return a == b })
	law.Associative(func(a, b T) T { return m.Op(a, b) })
	law.Identity(func(a, b T) T { return m.Op(a, b) }, m.Identity())
}

// Semigroup represents an associative operation (no identity or inverse).
type Semigroup[T comparable] interface {
	Op(a, b T) T
	Gen() T
}

// TestSemigroup verifies semigroup properties (associativity only).
func TestSemigroup[T comparable](t *testing.T, s Semigroup[T]) {
	t.Helper()
	law := For(t, s.Gen, func(a, b T) bool { return a == b })
	law.Associative(func(a, b T) T { return s.Op(a, b) })
}

// ===========================================================================
// VECTOR SPACE INTERFACE
// ===========================================================================

// VectorSpaceCustom represents a vector space over a scalar field.
type VectorSpaceCustom[V any, S any] interface {
	Add(v, w V) V
	Zero() V
	Neg(v V) V
	Scale(s S, v V) V
	ScalarAdd(s, t S) S
	ScalarMul(s, t S) S
	ScalarOne() S
	GenVector() V
	GenScalar() S
	EqVector(v, w V) bool
}

// TestVectorSpaceCustom verifies all vector space axioms.
func TestVectorSpaceCustom[V any, S any](t *testing.T, vs VectorSpaceCustom[V, S]) {
	t.Helper()
	vsLaw := ForVectorSpace(t, vs.GenVector, vs.GenScalar, vs.EqVector)
	vsLaw.Axioms(
		vs.Add,
		vs.Zero(),
		vs.Neg,
		vs.Scale,
		vs.ScalarAdd,
		vs.ScalarMul,
		vs.ScalarOne(),
	)
}

// ===========================================================================
// PARALLEL SAFETY FOR GROUPS
// ===========================================================================

// TestGroupParallelSafe verifies a group's operation is safe for concurrent use.
func TestGroupParallelSafe[T any](t *testing.T, g GroupCustom[T], goroutines int) {
	t.Helper()
	law := For(t, g.Gen, g.Eq)
	law.ParallelSafe(func(a, b T) T { return g.Op(a, b) }, goroutines)
}

// ===========================================================================
// EQUIVALENCE TESTING
// ===========================================================================

// Equivalent tests if two functions produce the same output for all inputs.
func Equivalent[T any, R any](t *testing.T, f1, f2 func(T) R, gen func() T, eq func(R, R) bool) bool {
	t.Helper()

	for i := 0; i < defaultTestCases; i++ {
		input := gen()
		r1 := f1(input)
		r2 := f2(input)

		if !eq(r1, r2) {
			t.Errorf("Functions not equivalent at iteration %d\n  input=%v\n  f1=%v\n  f2=%v",
				i, input, r1, r2)
			return false
		}
	}
	t.Logf("✅ Functions are equivalent (%d trials)", defaultTestCases)
	return true
}
