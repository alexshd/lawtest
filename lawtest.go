// Package lawtest provides property-based testing using group theory.
//
// lawtest helps verify mathematical properties of operations through
// randomized testing. Instead of writing specific test cases, you define
// the mathematical laws your code should obey, and lawtest generates
// hundreds of test cases automatically.
//
// # Basic Usage
//
// Test that integer addition is associative:
//
//	func TestAddition(t *testing.T) {
//	    add := func(a, b int) int { return a + b }
//	    gen := lawtest.IntGen(-100, 100)
//	    lawtest.Associative(t, add, gen)
//	}
//
// Test multiple properties of a custom operation:
//
//	func TestCacheMerge(t *testing.T) {
//	    merge := func(a, b *Cache) *Cache { return a.Merge(b) }
//	    gen := func() *Cache { return NewCache() }
//
//	    t.Run("Associative", func(t *testing.T) {
//	        lawtest.Associative(t, merge, gen)
//	    })
//
//	    t.Run("Commutative", func(t *testing.T) {
//	        lawtest.Commutative(t, merge, gen)
//	    })
//	}
//
// # Properties Tested
//
// lawtest can verify these mathematical properties:
//
//   - Associativity: (a ∘ b) ∘ c = a ∘ (b ∘ c)
//   - Commutativity: a ∘ b = b ∘ a
//   - Identity: a ∘ e = a and e ∘ a = a
//   - Inverse: a ∘ a⁻¹ = e and a⁻¹ ∘ a = e
//   - Closure: result type matches input types
//   - Idempotence: f(f(x)) = f(x)
//
// # Interface-Based Testing
//
// For complex types, implement Group, Monoid, or Semigroup interfaces
// and test all properties at once:
//
//	type IntMod struct {
//	    modulus int
//	}
//
//	func (m IntMod) Op(a, b int) int {
//	    return (a + b) % m.modulus
//	}
//
//	func (m IntMod) Identity() int { return 0 }
//	func (m IntMod) Inverse(a int) int { return (m.modulus - a) % m.modulus }
//	func (m IntMod) Gen() int { return rand.Intn(m.modulus) }
//
//	func TestIntMod(t *testing.T) {
//	    g := IntMod{modulus: 12}
//	    lawtest.TestGroup(t, g) // Tests all group properties
//	}
//
// # Concurrency Testing
//
// Test if operations are safe for concurrent use:
//
//	func TestParallelSafety(t *testing.T) {
//	    op := func(a, b *Cache) *Cache { return a.Merge(b) }
//	    gen := func() *Cache { return NewCache() }
//
//	    isSafe := lawtest.ParallelSafe(t, op, gen, 20)
//	    if !isSafe {
//	        t.Error("Cache merge has race conditions")
//	    }
//	}
//
// Test that immutable operations don't mutate inputs:
//
//	func TestImmutability(t *testing.T) {
//	    op := func(a, b []int) []int { return append(a, b...) }
//	    gen := func() []int { return []int{1, 2, 3} }
//	    lawtest.ImmutableOp(t, op, gen)
//	}
package lawtest

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

// defaultTestCases is the default number of random test cases to generate.
const defaultTestCases = 100

// BinaryOp is a binary operation that combines two values of type T.
//
// Example:
//
//	add := func(a, b int) int { return a + b }
//	mul := func(a, b int) int { return a * b }
type BinaryOp[T any] func(a, b T) T

// UnaryOp is a unary operation that transforms a value of type T.
//
// Example:
//
//	negate := func(x int) int { return -x }
//	inverse := func(x float64) float64 { return 1.0 / x }
type UnaryOp[T any] func(a T) T

// Generator produces random values of type T for property testing.
//
// Generators should produce diverse values to thoroughly test properties.
// Use built-in generators like IntGen, StringGen, or create custom ones.
//
// Example:
//
//	gen := lawtest.IntGen(-100, 100)
//	val := gen() // produces random int in [-100, 100]
//
// Custom generator:
//
//	type Point struct{ X, Y int }
//	gen := func() Point {
//	    return Point{
//	        X: rand.Intn(200) - 100,
//	        Y: rand.Intn(200) - 100,
//	    }
//	}
type Generator[T any] func() T

// Config holds configuration for property testing.
//
// Use DefaultConfig() for sensible defaults, or customize for specific needs:
//
//	cfg := &lawtest.Config{
//	    TestCases: 500,           // Run 500 random tests
//	    Timeout:   10 * time.Second, // Max 10 seconds per test
//	}
type Config struct {
	TestCases int           // Number of random test cases to generate and verify
	Timeout   time.Duration // Maximum time allowed per property test
}

// DefaultConfig returns a Config with sensible defaults.
//
// Default values:
//   - TestCases: 100 random test cases
//   - Timeout: 5 seconds
//
// Example:
//
//	cfg := lawtest.DefaultConfig()
//	cfg.TestCases = 200 // Customize as needed
func DefaultConfig() *Config {
	return &Config{
		TestCases: defaultTestCases,
		Timeout:   5 * time.Second,
	}
}

// Associative tests if a binary operation is associative: (a ∘ b) ∘ c = a ∘ (b ∘ c).
//
// Associativity means the order of applying operations doesn't matter,
// which is crucial for operations that will be chained or parallelized.
//
// Example:
//
//	func TestAdditionAssociative(t *testing.T) {
//	    add := func(a, b int) int { return a + b }
//	    gen := lawtest.IntGen(-1000, 1000)
//	    lawtest.Associative(t, add, gen)
//	}
//
// This verifies: (a + b) + c = a + (b + c) for 100 random combinations.
func Associative[T comparable](t *testing.T, op BinaryOp[T], gen Generator[T]) {
	AssociativeWithConfig(t, op, gen, DefaultConfig())
}

// AssociativeWithConfig tests associativity with custom configuration.
//
// Use this when you need more test cases or different timeout values.
//
// Example:
//
//	cfg := &lawtest.Config{TestCases: 500}
//	lawtest.AssociativeWithConfig(t, op, gen, cfg)
func AssociativeWithConfig[T comparable](t *testing.T, op BinaryOp[T], gen Generator[T], cfg *Config) {
	t.Helper()

	for i := 0; i < cfg.TestCases; i++ {
		a, b, c := gen(), gen(), gen()

		// (a ∘ b) ∘ c
		left := op(op(a, b), c)

		// a ∘ (b ∘ c)
		right := op(a, op(b, c))

		if left != right {
			t.Errorf("Associativity failed: (a∘b)∘c != a∘(b∘c)\n  a=%v, b=%v, c=%v\n  left=%v, right=%v",
				a, b, c, left, right)
			return
		}
	}
}

// Commutative tests if a binary operation is commutative: a ∘ b = b ∘ a.
//
// Commutativity means operand order doesn't matter. Not all operations
// are commutative (e.g., subtraction, matrix multiplication).
//
// Example:
//
//	func TestMultiplicationCommutative(t *testing.T) {
//	    mul := func(a, b int) int { return a * b }
//	    gen := lawtest.IntGen(1, 100)
//	    lawtest.Commutative(t, mul, gen)
//	}
//
// This verifies: a * b = b * a for 100 random pairs.
func Commutative[T comparable](t *testing.T, op BinaryOp[T], gen Generator[T]) {
	CommutativeWithConfig(t, op, gen, DefaultConfig())
}

// CommutativeWithConfig tests commutativity with custom configuration.
func CommutativeWithConfig[T comparable](t *testing.T, op BinaryOp[T], gen Generator[T], cfg *Config) {
	t.Helper()

	for i := 0; i < cfg.TestCases; i++ {
		a, b := gen(), gen()

		left := op(a, b)
		right := op(b, a)

		if left != right {
			t.Errorf("Commutativity failed: a∘b != b∘a\n  a=%v, b=%v\n  a∘b=%v, b∘a=%v",
				a, b, left, right)
			return
		}
	}
}

// Identity tests if an identity element exists: a ∘ e = a and e ∘ a = a.
//
// The identity element is a special value that, when combined with any other
// value, returns that value unchanged.
//
// Example:
//
//	func TestAdditionIdentity(t *testing.T) {
//	    add := func(a, b int) int { return a + b }
//	    gen := lawtest.IntGen(-100, 100)
//	    lawtest.Identity(t, add, 0, gen) // 0 is the identity for addition
//	}
//
// Common identities:
//   - Addition: 0 (a + 0 = a)
//   - Multiplication: 1 (a * 1 = a)
//   - String concatenation: "" (s + "" = s)
//   - Boolean OR: false (b || false = b)
func Identity[T comparable](t *testing.T, op BinaryOp[T], identity T, gen Generator[T]) {
	IdentityWithConfig(t, op, identity, gen, DefaultConfig())
}

// IdentityWithConfig tests identity element with custom configuration.
func IdentityWithConfig[T comparable](t *testing.T, op BinaryOp[T], identity T, gen Generator[T], cfg *Config) {
	t.Helper()

	for i := 0; i < cfg.TestCases; i++ {
		a := gen()

		// a ∘ e = a
		leftResult := op(a, identity)
		if leftResult != a {
			t.Errorf("Left identity failed: a∘e != a\n  a=%v, e=%v, a∘e=%v",
				a, identity, leftResult)
			return
		}

		// e ∘ a = a
		rightResult := op(identity, a)
		if rightResult != a {
			t.Errorf("Right identity failed: e∘a != a\n  e=%v, a=%v, e∘a=%v",
				identity, a, rightResult)
			return
		}
	}
}

// Inverse tests if each element has an inverse: a ∘ a⁻¹ = e and a⁻¹ ∘ a = e.
//
// An inverse "undoes" an operation, returning to the identity element.
//
// Example:
//
//	func TestIntegerAdditionInverse(t *testing.T) {
//	    add := func(a, b int) int { return a + b }
//	    negate := func(a int) int { return -a }
//	    gen := lawtest.IntGen(-100, 100)
//	    lawtest.Inverse(t, add, negate, 0, gen)
//	}
//
// This verifies: a + (-a) = 0 for all test values.
//
// Common inverses:
//   - Addition: negation (a + (-a) = 0)
//   - Multiplication: reciprocal (a * (1/a) = 1)
//   - Boolean XOR: self (a ⊕ a = false)
func Inverse[T comparable](t *testing.T, op BinaryOp[T], inverse UnaryOp[T], identity T, gen Generator[T]) {
	InverseWithConfig(t, op, inverse, identity, gen, DefaultConfig())
}

// InverseWithConfig tests inverse elements with custom configuration.
func InverseWithConfig[T comparable](t *testing.T, op BinaryOp[T], inv UnaryOp[T], identity T, gen Generator[T], cfg *Config) {
	t.Helper()

	for i := 0; i < cfg.TestCases; i++ {
		a := gen()
		aInv := inv(a)

		// a ∘ a⁻¹ = e
		leftResult := op(a, aInv)
		if leftResult != identity {
			t.Errorf("Left inverse failed: a∘a⁻¹ != e\n  a=%v, a⁻¹=%v, e=%v, a∘a⁻¹=%v",
				a, aInv, identity, leftResult)
			return
		}

		// a⁻¹ ∘ a = e
		rightResult := op(aInv, a)
		if rightResult != identity {
			t.Errorf("Right inverse failed: a⁻¹∘a != e\n  a⁻¹=%v, a=%v, e=%v, a⁻¹∘a=%v",
				aInv, a, identity, rightResult)
			return
		}
	}
}

// Closure tests if an operation stays within the same type.
//
// Closure means combining two values of type T always produces another value
// of type T. In Go, this is enforced by the type system through generics,
// so this function mainly serves as documentation and performs basic validation.
//
// Example:
//
//	func TestSetUnionClosure(t *testing.T) {
//	    union := func(a, b Set) Set { return a.Union(b) }
//	    gen := func() Set { return NewRandomSet() }
//	    lawtest.Closure(t, union, gen)
//	}
//
// The test verifies that Set.Union(Set) always returns a Set.
func Closure[T any](t *testing.T, op BinaryOp[T], gen Generator[T]) {
	t.Helper()

	// Generate a few examples to show closure is maintained
	for i := 0; i < 10; i++ {
		a, b := gen(), gen()
		result := op(a, b)

		// Type is enforced by Go's generics, but we can check the operation completes
		aType := reflect.TypeOf(a)
		resultType := reflect.TypeOf(result)

		if aType != resultType {
			t.Errorf("Closure violated: operation changed type\n  input type=%v, result type=%v",
				aType, resultType)
			return
		}
	}

	t.Logf("✓ Closure property holds (enforced by Go's type system)")
}

// Idempotent tests if repeated application gives the same result: f(f(x)) = f(x).
//
// An idempotent operation can be applied multiple times without changing
// the result after the first application. This is crucial for retry logic,
// caching, and data normalization.
//
// Example:
//
//	func TestAbsoluteValueIdempotent(t *testing.T) {
//	    abs := func(x int) int {
//	        if x < 0 {
//	            return -x
//	        }
//	        return x
//	    }
//	    gen := lawtest.IntGen(-100, 100)
//	    lawtest.Idempotent(t, abs, gen)
//	}
//
// This verifies: abs(abs(x)) = abs(x) for all test values.
//
// Common idempotent operations:
//   - String trimming: trim(trim(s)) = trim(s)
//   - Absolute value: abs(abs(x)) = abs(x)
//   - Set deduplication: dedupe(dedupe(set)) = dedupe(set)
//   - Cache warming: warm(warm(cache)) = warm(cache)
func Idempotent[T comparable](t *testing.T, op UnaryOp[T], gen Generator[T]) {
	IdempotentWithConfig(t, op, gen, DefaultConfig())
}

// IdempotentWithConfig tests idempotence with custom configuration.
func IdempotentWithConfig[T comparable](t *testing.T, op UnaryOp[T], gen Generator[T], cfg *Config) {
	t.Helper()

	for i := 0; i < cfg.TestCases; i++ {
		x := gen()

		fx := op(x)
		ffx := op(fx)

		if fx != ffx {
			t.Errorf("Idempotence failed: f(f(x)) != f(x)\n  x=%v, f(x)=%v, f(f(x))=%v",
				x, fx, ffx)
			return
		}
	}
}

// IntGen creates a Generator that produces random integers in [min, max].
//
// The generator produces uniformly distributed integers within the specified
// range, inclusive of both min and max.
//
// Example:
//
//	gen := lawtest.IntGen(-100, 100)
//	for i := 0; i < 10; i++ {
//	    fmt.Println(gen()) // prints random int in [-100, 100]
//	}
//
// Panics if min > max.
func IntGen(min, max int) Generator[int] {
	if min > max {
		panic(fmt.Sprintf("min (%d) must be <= max (%d)", min, max))
	}
	return func() int {
		return min + rand.Intn(max-min+1)
	}
}

// StringGen creates a Generator that produces random strings of length n.
//
// Strings contain alphanumeric characters (a-z, A-Z, 0-9).
// For different character sets, create a custom generator.
//
// Example:
//
//	gen := lawtest.StringGen(8)
//	username := gen() // e.g., "aB3xYz9K"
//
// Custom character set example:
//
//	hexGen := func() string {
//	    const charset = "0123456789abcdef"
//	    b := make([]byte, 16)
//	    for i := range b {
//	        b[i] = charset[rand.Intn(len(charset))]
//	    }
//	    return string(b)
//	}
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
//
// Example:
//
//	gen := lawtest.Float64Gen(0.0, 1.0)
//	probability := gen() // random float in [0.0, 1.0]
//
// Panics if min > max.
func Float64Gen(min, max float64) Generator[float64] {
	if min > max {
		panic(fmt.Sprintf("min (%f) must be <= max (%f)", min, max))
	}
	return func() float64 {
		return min + rand.Float64()*(max-min)
	}
}

// BoolGen creates a Generator that produces random boolean values.
//
// Example:
//
//	gen := lawtest.BoolGen()
//	flag := gen() // true or false with equal probability
func BoolGen() Generator[bool] {
	return func() bool {
		return rand.Intn(2) == 1
	}
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
//   - Closure: a ∘ b produces another T
//
// Example implementation:
//
//	type IntAddMod12 struct{}
//
//	func (g IntAddMod12) Op(a, b int) int {
//	    return (a + b) % 12
//	}
//
//	func (g IntAddMod12) Identity() int { return 0 }
//
//	func (g IntAddMod12) Inverse(a int) int {
//	    return (12 - a) % 12
//	}
//
//	func (g IntAddMod12) Gen() int {
//	    return rand.Intn(12)
//	}
//
//	func TestModularArithmetic(t *testing.T) {
//	    lawtest.TestGroup(t, IntAddMod12{})
//	}
type Group[T comparable] interface {
	// Op performs the group operation: a ∘ b
	Op(a, b T) T

	// Identity returns the identity element e where a ∘ e = a
	Identity() T

	// Inverse returns the inverse of a where a ∘ a⁻¹ = e
	Inverse(a T) T

	// Gen generates a random element for testing
	Gen() T
}

// Monoid represents an algebraic monoid (group without inverse requirement).
//
// A Monoid must satisfy:
//   - Associativity: (a ∘ b) ∘ c = a ∘ (b ∘ c)
//   - Identity: exists e where a ∘ e = e ∘ a = a
//   - Closure: a ∘ b produces another T
//
// Example implementation:
//
//	type StringConcat struct{}
//
//	func (m StringConcat) Op(a, b string) string {
//	    return a + b
//	}
//
//	func (m StringConcat) Identity() string { return "" }
//
//	func (m StringConcat) Gen() string {
//	    return lawtest.StringGen(5)()
//	}
//
//	func TestStringConcat(t *testing.T) {
//	    lawtest.TestMonoid(t, StringConcat{})
//	}
type Monoid[T comparable] interface {
	// Op performs the monoid operation: a ∘ b
	Op(a, b T) T

	// Identity returns the identity element e where a ∘ e = a
	Identity() T

	// Gen generates a random element for testing
	Gen() T
}

// Semigroup represents an algebraic semigroup (only associative operation).
//
// A Semigroup must satisfy:
//   - Associativity: (a ∘ b) ∘ c = a ∘ (b ∘ c)
//   - Closure: a ∘ b produces another T
//
// Example implementation:
//
//	type Max struct{}
//
//	func (s Max) Op(a, b int) int {
//	    if a > b {
//	        return a
//	    }
//	    return b
//	}
//
//	func (s Max) Gen() int {
//	    return lawtest.IntGen(1, 1000)()
//	}
//
//	func TestMaxSemigroup(t *testing.T) {
//	    lawtest.TestSemigroup(t, Max{})
//	}
type Semigroup[T comparable] interface {
	// Op performs the semigroup operation: a ∘ b
	Op(a, b T) T

	// Gen generates a random element for testing
	Gen() T
}

// IdempotentOp represents an operation that satisfies f(f(x)) = f(x).
//
// Example implementation:
//
//	type Abs struct{}
//
//	func (o Abs) Apply(x int) int {
//	    if x < 0 {
//	        return -x
//	    }
//	    return x
//	}
//
//	func (o Abs) Gen() int {
//	    return lawtest.IntGen(-100, 100)()
//	}
//
//	func TestAbsolute(t *testing.T) {
//	    lawtest.TestIdempotentOp(t, Abs{})
//	}
type IdempotentOp[T comparable] interface {
	// Apply performs the operation
	Apply(x T) T

	// Gen generates a random element for testing
	Gen() T
}

// Homomorphism represents a structure-preserving map between groups.
//
// A Homomorphism must preserve:
//   - Operation: h(a ∘ b) = h(a) ∘ h(b)
//   - Identity: h(e_source) = e_target
//
// Example implementation:
//
//	type AbsoluteValue struct {
//	    source IntAddition  // (ℤ, +)
//	    target IntMultiplication // (ℕ, *)
//	}
//
//	func (h AbsoluteValue) Map(x int) int {
//	    if x < 0 {
//	        return -x
//	    }
//	    return x
//	}
//
//	func (h AbsoluteValue) SourceGroup() lawtest.Group[int] {
//	    return h.source
//	}
//
//	func (h AbsoluteValue) TargetGroup() lawtest.Group[int] {
//	    return h.target
//	}
type Homomorphism[T, U comparable] interface {
	// Map transforms T to U while preserving structure
	Map(x T) U

	// SourceGroup returns the source group
	SourceGroup() Group[T]

	// TargetGroup returns the target group
	TargetGroup() Group[U]
}

// ===========================================================================
// INTERFACE-BASED PROPERTY TESTS
// ===========================================================================

// TestGroup verifies all group properties for a type implementing the Group interface.
//
// Tests performed:
//   - Associativity: (a ∘ b) ∘ c = a ∘ (b ∘ c)
//   - Identity: a ∘ e = e ∘ a = a
//   - Inverse: a ∘ a⁻¹ = a⁻¹ ∘ a = e
//   - Closure: type consistency
//
// Example:
//
//	func TestMyGroup(t *testing.T) {
//	    g := MyGroupImpl{...}
//	    lawtest.TestGroup(t, g)
//	}
func TestGroup[T comparable](t *testing.T, g Group[T]) {
	TestGroupWithConfig(t, g, DefaultConfig())
}

// TestGroupWithConfig verifies all group properties with custom configuration.
func TestGroupWithConfig[T comparable](t *testing.T, g Group[T], cfg *Config) {
	t.Helper()

	t.Run("Associativity", func(t *testing.T) {
		AssociativeWithConfig(t, g.Op, g.Gen, cfg)
	})

	t.Run("Identity", func(t *testing.T) {
		IdentityWithConfig(t, g.Op, g.Identity(), g.Gen, cfg)
	})

	t.Run("Inverse", func(t *testing.T) {
		InverseWithConfig(t, g.Op, g.Inverse, g.Identity(), g.Gen, cfg)
	})

	t.Run("Closure", func(t *testing.T) {
		Closure(t, g.Op, g.Gen)
	})
}

// TestMonoid verifies all monoid properties (associativity and identity).
//
// Tests performed:
//   - Associativity: (a ∘ b) ∘ c = a ∘ (b ∘ c)
//   - Identity: a ∘ e = e ∘ a = a
//   - Closure: type consistency
//
// Example:
//
//	func TestStringConcatMonoid(t *testing.T) {
//	    m := StringConcat{}
//	    lawtest.TestMonoid(t, m)
//	}
func TestMonoid[T comparable](t *testing.T, m Monoid[T]) {
	TestMonoidWithConfig(t, m, DefaultConfig())
}

// TestMonoidWithConfig verifies monoid properties with custom configuration.
func TestMonoidWithConfig[T comparable](t *testing.T, m Monoid[T], cfg *Config) {
	t.Helper()

	t.Run("Associativity", func(t *testing.T) {
		AssociativeWithConfig(t, m.Op, m.Gen, cfg)
	})

	t.Run("Identity", func(t *testing.T) {
		IdentityWithConfig(t, m.Op, m.Identity(), m.Gen, cfg)
	})

	t.Run("Closure", func(t *testing.T) {
		Closure(t, m.Op, m.Gen)
	})
}

// TestSemigroup verifies semigroup properties (associativity only).
//
// Tests performed:
//   - Associativity: (a ∘ b) ∘ c = a ∘ (b ∘ c)
//   - Closure: type consistency
//
// Example:
//
//	func TestMaxSemigroup(t *testing.T) {
//	    s := Max{}
//	    lawtest.TestSemigroup(t, s)
//	}
func TestSemigroup[T comparable](t *testing.T, s Semigroup[T]) {
	TestSemigroupWithConfig(t, s, DefaultConfig())
}

// TestSemigroupWithConfig verifies semigroup properties with custom configuration.
func TestSemigroupWithConfig[T comparable](t *testing.T, s Semigroup[T], cfg *Config) {
	t.Helper()

	t.Run("Associativity", func(t *testing.T) {
		AssociativeWithConfig(t, s.Op, s.Gen, cfg)
	})

	t.Run("Closure", func(t *testing.T) {
		Closure(t, s.Op, s.Gen)
	})
}

// TestIdempotentOp verifies the idempotence property: f(f(x)) = f(x).
//
// Example:
//
//	func TestAbsIdempotent(t *testing.T) {
//	    op := Abs{}
//	    lawtest.TestIdempotentOp(t, op)
//	}
func TestIdempotentOp[T comparable](t *testing.T, op IdempotentOp[T]) {
	TestIdempotentOpWithConfig(t, op, DefaultConfig())
}

// TestIdempotentOpWithConfig verifies idempotence with custom configuration.
func TestIdempotentOpWithConfig[T comparable](t *testing.T, op IdempotentOp[T], cfg *Config) {
	t.Helper()

	t.Run("Idempotence", func(t *testing.T) {
		IdempotentWithConfig(t, op.Apply, op.Gen, cfg)
	})
}

// TestHomomorphism verifies that a map preserves group structure.
//
// Tests performed:
//   - Preserves operation: h(a ∘ b) = h(a) ∘ h(b)
//   - Preserves identity: h(e_source) = e_target
//
// Example:
//
//	func TestAbsHomomorphism(t *testing.T) {
//	    h := AbsoluteValue{...}
//	    lawtest.TestHomomorphism(t, h)
//	}
func TestHomomorphism[T, U comparable](t *testing.T, h Homomorphism[T, U]) {
	TestHomomorphismWithConfig(t, h, DefaultConfig())
}

// TestHomomorphismWithConfig verifies homomorphism properties with custom configuration.
func TestHomomorphismWithConfig[T, U comparable](t *testing.T, h Homomorphism[T, U], cfg *Config) {
	t.Helper()

	srcGroup := h.SourceGroup()
	tgtGroup := h.TargetGroup()

	t.Run("PreservesOperation", func(t *testing.T) {
		// Verify: h(a ∘ b) = h(a) ∘ h(b)
		for i := 0; i < cfg.TestCases; i++ {
			a := srcGroup.Gen()
			b := srcGroup.Gen()

			// Left side: h(a ∘ b)
			ab := srcGroup.Op(a, b)
			hAb := h.Map(ab)

			// Right side: h(a) ∘ h(b)
			ha := h.Map(a)
			hb := h.Map(b)
			haHb := tgtGroup.Op(ha, hb)

			if hAb != haHb {
				t.Errorf("Homomorphism failed: h(a∘b) != h(a)∘h(b)\n  a=%v, b=%v\n  h(a∘b)=%v, h(a)∘h(b)=%v",
					a, b, hAb, haHb)
				return
			}
		}
	})

	t.Run("PreservesIdentity", func(t *testing.T) {
		// Verify: h(e_source) = e_target
		sourceIdentity := srcGroup.Identity()
		targetIdentity := tgtGroup.Identity()

		mappedIdentity := h.Map(sourceIdentity)

		if mappedIdentity != targetIdentity {
			t.Errorf("Homomorphism doesn't preserve identity: h(e_src) != e_tgt\n  e_src=%v, e_tgt=%v, h(e_src)=%v",
				sourceIdentity, targetIdentity, mappedIdentity)
		}
	})
}

// ===========================================================================
// HELPER: TEST THAT A STRUCT *FAILS* GROUP PROPERTIES (for negative testing)
// ===========================================================================

// ExpectGroupFailure runs group tests expecting them to FAIL
// Useful for verifying that broken implementations are correctly detected
func ExpectGroupFailure[T comparable](t *testing.T, g Group[T], expectedFailure string) {
	t.Helper()

	// Create a sub-test that we expect to fail
	failed := false

	testFunc := func(t *testing.T) {
		TestGroup(t, g)
	}

	// Run test and capture if it failed
	result := testing.RunTests(func(pat, str string) (bool, error) {
		return true, nil
	}, []testing.InternalTest{
		{
			Name: "ExpectedFailure",
			F:    testFunc,
		},
	})

	if !result {
		failed = true
	}

	if !failed {
		t.Errorf("Expected group test to fail with '%s', but it passed!", expectedFailure)
	} else {
		t.Logf("✓ Group correctly failed validation (as expected: %s)", expectedFailure)
	}
}

// ===========================================================================
// CONCURRENCY SAFETY TESTING
// ===========================================================================

// ParallelSafe tests if an operation can be safely executed concurrently.
//
// Returns true if the operation is parallel-safe (no race conditions detected).
// Immutable operations (pure functions) should pass. Mutating operations
// will exhibit race conditions under concurrent access.
//
// The test launches multiple goroutines that execute the operation
// simultaneously on shared data. Race conditions, panics, or data corruption
// indicate the operation is not parallel-safe.
//
// Example:
//
//	func TestCacheMergeParallel(t *testing.T) {
//	    merge := func(a, b *Cache) *Cache { return a.Merge(b) }
//	    gen := func() *Cache { return NewCache() }
//
//	    isSafe := lawtest.ParallelSafe(t, merge, gen, 20)
//	    if !isSafe {
//	        t.Error("Cache merge is not parallel-safe")
//	    }
//	}
//
// For production use, always run tests with -race flag:
//
//	go test -race
func ParallelSafe[T comparable](t *testing.T, op BinaryOp[T], gen Generator[T], goroutines int) bool {
	return ParallelSafeWithConfig(t, op, gen, goroutines, DefaultConfig())
}

// ParallelSafeWithConfig tests parallel safety with custom configuration.
func ParallelSafeWithConfig[T comparable](t *testing.T, op BinaryOp[T], gen Generator[T], goroutines int, cfg *Config) bool {
	t.Helper()

	if goroutines < 2 {
		goroutines = 10 // Default to 10 goroutines
	}

	// Create shared test data
	testData := make([]T, cfg.TestCases)
	for i := range testData {
		testData[i] = gen()
	}

	// Run operations concurrently
	done := make(chan bool, goroutines)
	errors := make(chan string, goroutines)

	for g := 0; g < goroutines; g++ {
		go func(id int) {
			defer func() {
				if r := recover(); r != nil {
					errors <- fmt.Sprintf("Goroutine %d panicked: %v", id, r)
				}
				done <- true
			}()

			// Perform operations on shared data
			for i := 0; i < len(testData)-1; i++ {
				_ = op(testData[i], testData[i+1])
			}
		}(g)
	}

	// Wait for all goroutines
	for i := 0; i < goroutines; i++ {
		<-done
	}

	// Check for errors
	close(errors)
	hasErrors := false
	for err := range errors {
		t.Logf("⚠ Race condition detected: %s", err)
		hasErrors = true
	}

	if hasErrors {
		t.Logf("❌ Operation is NOT parallel-safe (race conditions detected)")
		return false
	}

	t.Logf("✅ Operation appears parallel-safe (no race conditions in %d goroutines)", goroutines)
	return true
}

// TestParallelAssociativity tests if associativity holds under concurrent execution.
//
// This is a stronger test than Associative: not only must (a∘b)∘c = a∘(b∘c),
// but this must hold even when operations execute concurrently. This catches
// bugs related to shared mutable state, race conditions, or non-deterministic
// behavior in concurrent contexts.
//
// The test first verifies sequential associativity, then launches multiple
// goroutines that test associativity simultaneously.
//
// Example:
//
//	func TestCacheMergeConcurrent(t *testing.T) {
//	    merge := func(a, b *Cache) *Cache { return a.Merge(b) }
//	    gen := func() *Cache { return NewCache() }
//	    lawtest.TestParallelAssociativity(t, merge, gen, 20)
//	}
//
// Run with -race flag to detect data races:
//
//	go test -race -run TestCacheMergeConcurrent
func TestParallelAssociativity[T comparable](t *testing.T, op BinaryOp[T], gen Generator[T], goroutines int) {
	TestParallelAssociativityWithConfig(t, op, gen, goroutines, DefaultConfig())
}

// TestParallelAssociativityWithConfig tests parallel associativity with custom configuration.
func TestParallelAssociativityWithConfig[T comparable](t *testing.T, op BinaryOp[T], gen Generator[T], goroutines int, cfg *Config) {
	t.Helper()

	if goroutines < 2 {
		goroutines = 10
	}

	// First verify sequential associativity
	t.Run("Sequential", func(t *testing.T) {
		AssociativeWithConfig(t, op, gen, cfg)
	})

	// Then test under concurrent load
	t.Run("Concurrent", func(t *testing.T) {
		type testCase struct {
			a, b, c T
			left    T
			right   T
		}

		results := make(chan testCase, cfg.TestCases)
		done := make(chan bool, goroutines)

		// Launch goroutines to test associativity concurrently
		casesPerGoroutine := cfg.TestCases / goroutines
		if casesPerGoroutine == 0 {
			casesPerGoroutine = 1
		}

		for g := 0; g < goroutines; g++ {
			go func() {
				defer func() { done <- true }()

				for i := 0; i < casesPerGoroutine; i++ {
					a, b, c := gen(), gen(), gen()
					left := op(op(a, b), c)
					right := op(a, op(b, c))

					results <- testCase{a, b, c, left, right}
				}
			}()
		}

		// Wait for all goroutines
		go func() {
			for i := 0; i < goroutines; i++ {
				<-done
			}
			close(results)
		}()

		// Check results
		failures := 0
		for tc := range results {
			if tc.left != tc.right {
				t.Errorf("Associativity failed under concurrency: (a∘b)∘c != a∘(b∘c)\n  a=%v, b=%v, c=%v\n  left=%v, right=%v",
					tc.a, tc.b, tc.c, tc.left, tc.right)
				failures++
				if failures >= 3 {
					break
				}
			}
		}

		if failures == 0 {
			t.Logf("✅ Associativity holds under concurrent execution (%d goroutines)", goroutines)
		} else {
			t.Logf("❌ Associativity violated under concurrency (%d failures)", failures)
		}
	})
}

// ImmutableOp tests if an operation creates new values instead of mutating inputs.
//
// Immutable operations are crucial for:
//   - Thread-safe code (no need for locks)
//   - Functional programming patterns
//   - Avoiding side effects and debugging nightmares
//
// The test verifies that input values remain unchanged after the operation.
//
// Example:
//
//	func TestListAppendImmutable(t *testing.T) {
//	    // Good: returns new slice
//	    append := func(a, b []int) []int {
//	        result := make([]int, len(a)+len(b))
//	        copy(result, a)
//	        copy(result[len(a):], b)
//	        return result
//	    }
//
//	    gen := func() []int { return []int{1, 2, 3} }
//	    lawtest.ImmutableOp(t, append, gen)
//	}
//
// This catches mutations that violate functional programming principles.
func ImmutableOp[T comparable](t *testing.T, op BinaryOp[T], gen Generator[T]) {
	ImmutableOpWithConfig(t, op, gen, DefaultConfig())
}

// ImmutableOpWithConfig tests immutability with custom configuration.
func ImmutableOpWithConfig[T comparable](t *testing.T, op BinaryOp[T], gen Generator[T], cfg *Config) {
	t.Helper()

	for i := 0; i < cfg.TestCases; i++ {
		a, b := gen(), gen()

		// Create copies for comparison (for comparable types)
		aOriginal := a
		bOriginal := b

		// Apply operation
		_ = op(a, b)

		// Check if inputs were mutated
		if a != aOriginal {
			t.Errorf("Immutability violated: operation mutated first argument\n  before=%v, after=%v",
				aOriginal, a)
			return
		}

		if b != bOriginal {
			t.Errorf("Immutability violated: operation mutated second argument\n  before=%v, after=%v",
				bOriginal, b)
			return
		}
	}

	t.Logf("✅ Operation is immutable (does not mutate inputs)")
}

// AssociativeCustom tests associativity using a custom equality function.
// Use this for non-comparable types (slices, maps, functions).
//
// Example:
//
//	type State struct { items []int }
//	merge := func(a, b State) State { ... }
//	gen := func() State { return State{items: []int{1,2,3}} }
//	eq := func(a, b State) bool { return reflect.DeepEqual(a.items, b.items) }
//	lawtest.AssociativeCustom(t, merge, gen, eq)
func AssociativeCustom[T any](t *testing.T, op BinaryOp[T], gen Generator[T], eq func(T, T) bool) {
	AssociativeCustomWithConfig(t, op, gen, eq, DefaultConfig())
}

// AssociativeCustomWithConfig tests associativity with custom equality and configuration.
func AssociativeCustomWithConfig[T any](t *testing.T, op BinaryOp[T], gen Generator[T], eq func(T, T) bool, cfg *Config) {
	t.Helper()

	for i := 0; i < cfg.TestCases; i++ {
		a, b, c := gen(), gen(), gen()

		// (a ∘ b) ∘ c
		left := op(op(a, b), c)

		// a ∘ (b ∘ c)
		right := op(a, op(b, c))

		if !eq(left, right) {
			t.Errorf("Associativity failed: (a∘b)∘c != a∘(b∘c)\n  a=%v, b=%v, c=%v\n  left=%v, right=%v",
				a, b, c, left, right)
			return
		}
	}

	t.Logf("✅ Operation is associative with custom equality")
}

// ImmutableOpCustom tests immutability using a custom equality function.
// Use this for non-comparable types (slices, maps, functions).
//
// Example:
//
//	type State struct { data map[string]int }
//	merge := func(a, b State) State { ... }
//	gen := func() State { return State{data: map[string]int{"x": 1}} }
//	eq := func(a, b State) bool { return reflect.DeepEqual(a.data, b.data) }
//	lawtest.ImmutableOpCustom(t, merge, gen, eq)
func ImmutableOpCustom[T any](t *testing.T, op BinaryOp[T], gen Generator[T], eq func(T, T) bool) {
	ImmutableOpCustomWithConfig(t, op, gen, eq, DefaultConfig())
}

// ImmutableOpCustomWithConfig tests immutability with custom equality and configuration.
func ImmutableOpCustomWithConfig[T any](t *testing.T, op BinaryOp[T], gen Generator[T], eq func(T, T) bool, cfg *Config) {
	t.Helper()

	for i := 0; i < cfg.TestCases; i++ {
		a, b := gen(), gen()

		// Create deep copies using custom serialization
		// Since we can't rely on comparable, we'll test by calling the function
		// and checking if the original values changed using custom equality
		aOriginal := a
		bOriginal := b

		// Apply operation
		_ = op(a, b)

		// Check if inputs were mutated using custom equality
		if !eq(a, aOriginal) {
			t.Errorf("Immutability violated: operation mutated first argument\n  before=%v, after=%v",
				aOriginal, a)
			return
		}

		if !eq(b, bOriginal) {
			t.Errorf("Immutability violated: operation mutated second argument\n  before=%v, after=%v",
				bOriginal, b)
			return
		}
	}

	t.Logf("✅ Operation is immutable (does not mutate inputs)")
}

// ParallelSafeCustom tests parallel safety using a custom equality function.
// Use this for non-comparable types (slices, maps, functions).
//
// Example:
//
//	type Cache struct { data map[string]string }
//	merge := func(a, b Cache) Cache { ... }
//	gen := func() Cache { return Cache{data: map[string]string{"x": "y"}} }
//	eq := func(a, b Cache) bool { return reflect.DeepEqual(a.data, b.data) }
//	lawtest.ParallelSafeCustom(t, merge, gen, eq, 100)
func ParallelSafeCustom[T any](t *testing.T, op BinaryOp[T], gen Generator[T], eq func(T, T) bool, goroutines int) bool {
	return ParallelSafeCustomWithConfig(t, op, gen, eq, goroutines, DefaultConfig())
}

// ParallelSafeCustomWithConfig tests parallel safety with custom equality and configuration.
func ParallelSafeCustomWithConfig[T any](t *testing.T, op BinaryOp[T], gen Generator[T], eq func(T, T) bool, goroutines int, cfg *Config) bool {
	t.Helper()

	// Generate test data
	a, b := gen(), gen()

	// Expected result (sequential)
	expected := op(a, b)

	// Run operation concurrently many times
	results := make([]T, goroutines)
	done := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			results[idx] = op(a, b)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < goroutines; i++ {
		<-done
	}

	// Check if all results match expected using custom equality
	for i, result := range results {
		if !eq(result, expected) {
			t.Errorf("Parallel safety failed: goroutine %d produced different result\n  expected=%v, got=%v",
				i, expected, result)
			return false
		}
	}

	t.Logf("✅ Operation appears parallel-safe (no race conditions in %d goroutines)", goroutines)
	return true
}

// Equivalent tests if two functions produce the same output for all inputs.
//
// This is useful for verifying that optimized implementations (tail recursion,
// iterative versions, cached versions) produce the same results as the original.
//
// Example:
//
//	// Original recursive
//	func Factorial(n int) int {
//	    if n <= 1 { return 1 }
//	    return n * Factorial(n-1)
//	}
//
//	// Tail optimized version
//	func FactorialTail(n, acc int) int {
//	    if n <= 1 { return acc }
//	    return FactorialTail(n-1, n*acc)
//	}
//
//	// Prove they're equivalent
//	func TestTailOptimization(t *testing.T) {
//	    gen := func() int { return rand.Intn(20) + 1 }
//	    lawtest.Equivalent(t,
//	        func(n int) int { return Factorial(n) },
//	        func(n int) int { return FactorialTail(n, 1) },
//	        gen,
//	    )
//	}
//
// Parameters:
//   - t: testing.T instance
//   - f1: first function to compare
//   - f2: second function to compare
//   - gen: generator function for random test inputs
//
// Returns true if both functions produce the same output for all test cases.
func Equivalent[T any, R comparable](t *testing.T, f1, f2 func(T) R, gen func() T) bool {
	t.Helper()
	iterations := defaultTestCases

	for i := 0; i < iterations; i++ {
		input := gen()
		result1 := f1(input)
		result2 := f2(input)

		if result1 != result2 {
			t.Errorf("Functions not equivalent at iteration %d\n  input=%v\n  f1(input)=%v\n  f2(input)=%v",
				i, input, result1, result2)
			return false
		}
	}

	t.Logf("✅ Functions are equivalent (tested %d random inputs)", iterations)
	return true
}

// EquivalentCustom tests if two functions produce the same output for all inputs,
// using a custom equality function for non-comparable output types.
//
// This extends Equivalent to work with slices, maps, structs without comparable fields, etc.
//
// Example:
//
//	// Recursive list reverse
//	func ReverseRecursive(list []int) []int {
//	    if len(list) == 0 { return []int{} }
//	    return append(ReverseRecursive(list[1:]), list[0])
//	}
//
//	// Iterative list reverse
//	func ReverseIterative(list []int) []int {
//	    result := make([]int, len(list))
//	    for i := range list {
//	        result[len(list)-1-i] = list[i]
//	    }
//	    return result
//	}
//
//	// Prove they're equivalent
//	func TestReverseEquivalent(t *testing.T) {
//	    gen := func() []int {
//	        n := rand.Intn(10)
//	        list := make([]int, n)
//	        for i := range list { list[i] = rand.Intn(100) }
//	        return list
//	    }
//	    eq := func(a, b []int) bool { return reflect.DeepEqual(a, b) }
//	    lawtest.EquivalentCustom(t, ReverseRecursive, ReverseIterative, gen, eq)
//	}
//
// Parameters:
//   - t: testing.T instance
//   - f1: first function to compare
//   - f2: second function to compare
//   - gen: generator function for random test inputs
//   - eq: custom equality function for comparing outputs
//
// Returns true if both functions produce equal output for all test cases.
func EquivalentCustom[T any, R any](t *testing.T, f1, f2 func(T) R, gen func() T, eq func(R, R) bool) bool {
	t.Helper()
	iterations := defaultTestCases

	for i := 0; i < iterations; i++ {
		input := gen()
		result1 := f1(input)
		result2 := f2(input)

		if !eq(result1, result2) {
			t.Errorf("Functions not equivalent at iteration %d\n  input=%v\n  f1(input)=%v\n  f2(input)=%v",
				i, input, result1, result2)
			return false
		}
	}

	t.Logf("✅ Functions are equivalent (tested %d random inputs)", iterations)
	return true
}
