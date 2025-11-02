// Package goprops provides property-based testing using group theory
package goprops

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

const (
	defaultTestCases = 100 // Number of random test cases to generate
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// BinaryOp represents a binary operation: (a, b) -> result
type BinaryOp[T any] func(a, b T) T

// UnaryOp represents a unary operation: a -> result
type UnaryOp[T any] func(a T) T

// Generator creates random values of type T
type Generator[T any] func() T

// Config holds testing configuration
type Config struct {
	TestCases int           // Number of test cases to run
	Timeout   time.Duration // Max time per property test
}

// DefaultConfig returns default testing configuration
func DefaultConfig() *Config {
	return &Config{
		TestCases: defaultTestCases,
		Timeout:   5 * time.Second,
	}
}

// Associative tests if operation satisfies: (a ∘ b) ∘ c = a ∘ (b ∘ c)
func Associative[T comparable](t *testing.T, op BinaryOp[T], gen Generator[T]) {
	AssociativeWithConfig(t, op, gen, DefaultConfig())
}

// AssociativeWithConfig tests associativity with custom config
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

// Commutative tests if operation satisfies: a ∘ b = b ∘ a
func Commutative[T comparable](t *testing.T, op BinaryOp[T], gen Generator[T]) {
	CommutativeWithConfig(t, op, gen, DefaultConfig())
}

// CommutativeWithConfig tests commutativity with custom config
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

// Identity tests if identity element exists: a ∘ e = a and e ∘ a = a
func Identity[T comparable](t *testing.T, op BinaryOp[T], identity T, gen Generator[T]) {
	IdentityWithConfig(t, op, identity, gen, DefaultConfig())
}

// IdentityWithConfig tests identity with custom config
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

// Inverse tests if each element has an inverse: a ∘ a⁻¹ = e and a⁻¹ ∘ a = e
func Inverse[T comparable](t *testing.T, op BinaryOp[T], inverse UnaryOp[T], identity T, gen Generator[T]) {
	InverseWithConfig(t, op, inverse, identity, gen, DefaultConfig())
}

// InverseWithConfig tests inverse with custom config
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

// Closure tests if operation stays within the same type (checked at compile time with generics)
// This function mainly serves as documentation - Go's type system enforces closure
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

// Idempotent tests if repeated application gives same result: f(f(x)) = f(x)
func Idempotent[T comparable](t *testing.T, op UnaryOp[T], gen Generator[T]) {
	IdempotentWithConfig(t, op, gen, DefaultConfig())
}

// IdempotentWithConfig tests idempotence with custom config
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

// Common generators for built-in types

// IntGen generates random integers in range [min, max]
func IntGen(min, max int) Generator[int] {
	if min > max {
		panic(fmt.Sprintf("min (%d) must be <= max (%d)", min, max))
	}
	return func() int {
		return min + rand.Intn(max-min+1)
	}
}

// StringGen generates random strings of length n
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

// Float64Gen generates random float64 in range [min, max]
func Float64Gen(min, max float64) Generator[float64] {
	if min > max {
		panic(fmt.Sprintf("min (%f) must be <= max (%f)", min, max))
	}
	return func() float64 {
		return min + rand.Float64()*(max-min)
	}
}

// BoolGen generates random booleans
func BoolGen() Generator[bool] {
	return func() bool {
		return rand.Intn(2) == 1
	}
}
