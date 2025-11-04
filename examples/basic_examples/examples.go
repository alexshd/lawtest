// Package main shows basic usage patterns for lawtest
package main

import (
	"fmt"

	"github.com/alexshd/lawtest"
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

// Example 2: Matrix operations (demonstrating non-commutative multiplication)
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

func MatrixGen() lawtest.Generator[Matrix2x2] {
	intGen := lawtest.IntGen(-10, 10)
	return func() Matrix2x2 {
		return Matrix2x2{
			A: intGen(),
			B: intGen(),
			C: intGen(),
			D: intGen(),
		}
	}
}

// Example 3: List concatenation forms a monoid
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

// Example 4: Property-based testing for caching
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

func main() {
	fmt.Println("lawtest basic examples - run with: go test -v")
	fmt.Println()
	fmt.Println("Key concepts to explore:")
	fmt.Println("1. Monoid: Associative operation + identity element")
	fmt.Println("2. Group: Monoid + every element has an inverse")
	fmt.Println("3. Non-commutative operations: Matrix multiplication")
	fmt.Println("4. Idempotence: f(f(x)) = f(x)")
	fmt.Println("5. Approximate equality: Floating point with epsilon")
}
