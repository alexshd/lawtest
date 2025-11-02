// Package goprops_examples shows advanced usage patternspackage examples

package main

import (
	"fmt"
	"testing"

	"github.com/alexshd/goprops"
)

// Example 1: Testing a custom Set type
type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func (s Set[T]) Add(item T) Set[T] {
	result := NewSet[T]()
	for k := range s {
		result[k] = struct{}{}
	}
	result[item] = struct{}{}
	return result
}

func (s Set[T]) Union(other Set[T]) Set[T] {
	result := NewSet[T]()
	for k := range s {
		result[k] = struct{}{}
	}
	for k := range other {
		result[k] = struct{}{}
	}
	return result
}

// Example 2: Monoid - a set with an associative binary operation and identity
type Monoid[T comparable] struct {
	Op       goprops.BinaryOp[T]
	Identity T
}

// TestMonoidLaws verifies all monoid properties at once
func TestMonoidLaws(t *testing.T, m Monoid[int], gen goprops.Generator[int]) {
	t.Run("Associative", func(t *testing.T) {
		goprops.Associative(t, m.Op, gen)
	})

	t.Run("Identity", func(t *testing.T) {
		goprops.Identity(t, m.Op, m.Identity, gen)
	})
}

// Example 3: Group - a monoid where every element has an inverse
type Group[T comparable] struct {
	Monoid[T]
	Inverse goprops.UnaryOp[T]
}

// TestGroupLaws verifies all group properties
func TestGroupLaws(t *testing.T, g Group[int], gen goprops.Generator[int]) {
	TestMonoidLaws(t, g.Monoid, gen)

	t.Run("Inverse", func(t *testing.T) {
		goprops.Inverse(t, g.Op, g.Inverse, g.Identity, gen)
	})
}

// Example 4: Testing that float addition is "almost" associative (within epsilon)
type ApproxFloat struct {
	Value   float64
	Epsilon float64
}

func (a ApproxFloat) Equals(b ApproxFloat) bool {
	diff := a.Value - b.Value
	if diff < 0 {
		diff = -diff
	}
	return diff < a.Epsilon
}

// Example 5: Testing concurrent operations preserve properties
type ConcurrentCounter struct {
	value int
}

func (c *ConcurrentCounter) Increment() int {
	c.value++
	return c.value
}

// Example 6: Matrix operations (demonstrating non-commutative multiplication)
type Matrix2x2 struct {
	A, B, C, D int
}

func (m Matrix2x2) Multiply(other Matrix2x2) Matrix2x2 {
	return Matrix2x2{
		A: m.A*other.A + m.B*other.C,
		B: m.A*other.B + m.B*other.D,
		C: m.C*other.A + m.D*other.C,
		D: m.C*other.B + m.D*other.D,
	}
}

func MatrixGen() goprops.Generator[Matrix2x2] {
	intGen := goprops.IntGen(-10, 10)
	return func() Matrix2x2 {
		return Matrix2x2{
			A: intGen(),
			B: intGen(),
			C: intGen(),
			D: intGen(),
		}
	}
}

func TestMatrixMultiplication(t *testing.T) {
	mul := func(a, b Matrix2x2) Matrix2x2 { return a.Multiply(b) }
	gen := MatrixGen()

	t.Run("Associative", func(t *testing.T) {
		// Matrix multiplication IS associative: (AB)C = A(BC)
		goprops.Associative(t, mul, gen)
	})

	t.Run("NOT Commutative", func(t *testing.T) {
		// This test SHOULD FAIL - matrix multiplication is NOT commutative
		// Uncomment to see it fail:
		// goprops.Commutative(t, mul, gen)
		t.Log("Matrix multiplication is NOT commutative (AB != BA)")
	})
}

// Example 7: List concatenation forms a monoid
type List[T any] []T

func (l List[T]) Concat(other List[T]) List[T] {
	result := make(List[T], len(l)+len(other))
	copy(result, l)
	copy(result[len(l):], other)
	return result
}

// Since List contains slices, we need to make it comparable
// We'll use a wrapper that implements equality
type IntList struct {
	Items []int
}

func (a IntList) Equals(b IntList) bool {
	if len(a.Items) != len(b.Items) {
		return false
	}
	for i := range a.Items {
		if a.Items[i] != b.Items[i] {
			return false
		}
	}
	return true
}

// Example 8: Property-based testing for caching
// Tests that adding same key twice keeps last value (idempotent-ish)
type Cache[K comparable, V any] struct {
	data map[K]V
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{data: make(map[K]V)}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.data[key] = value
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	v, ok := c.data[key]
	return v, ok
}

// Example: Testing last-write-wins property
func TestCacheLastWriteWins(t *testing.T) {
	cache := NewCache[string, int]()
	key := "test"

	for i := 0; i < 100; i++ {
		vals := []int{
			goprops.IntGen(1, 1000)(),
			goprops.IntGen(1, 1000)(),
		}

		cache.Set(key, vals[0])
		cache.Set(key, vals[1])

		result, ok := cache.Get(key)
		if !ok {
			t.Fatal("Key not found")
		}
		if result != vals[1] {
			t.Errorf("Last write didn't win: expected %d, got %d", vals[1], result)
		}
	}
}

func main() {
	fmt.Println("goprops examples - run with: go test -v examples.go")
	fmt.Println()
	fmt.Println("Key concepts to explore:")
	fmt.Println("1. Monoid: Associative operation + identity element")
	fmt.Println("2. Group: Monoid + every element has an inverse")
	fmt.Println("3. Non-commutative operations: Matrix multiplication")
	fmt.Println("4. Idempotence: f(f(x)) = f(x)")
	fmt.Println("5. Approximate equality: Floating point with epsilon")
}
