// Package main shows basic usage patterns for lawtest
package main

import (
	"testing"

	"github.com/alexshd/lawtest"
)

func TestMatrixMultiplication(t *testing.T) {
	mul := func(a, b Matrix2x2) Matrix2x2 { return a.Multiply(b) }
	gen := MatrixGen()

	t.Run("Associative", func(t *testing.T) {
		// Matrix multiplication IS associative: (AB)C = A(BC)
		lawtest.Associative(t, mul, gen)
	})

	t.Run("NOT Commutative", func(t *testing.T) {
		// This test SHOULD FAIL - matrix multiplication is NOT commutative
		// Uncomment to see it fail:
		// lawtest.Commutative(t, mul, gen)
		t.Log("Matrix multiplication is NOT commutative (AB != BA)")
	})
}

// Example: Testing last-write-wins property
func TestCacheLastWriteWins(t *testing.T) {
	cache := NewCache[string, int]()
	key := "test"

	for i := 0; i < 100; i++ {
		vals := []int{
			lawtest.IntGen(1, 1000)(),
			lawtest.IntGen(1, 1000)(),
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

// Test Set Union is associative
func TestSetUnionAssociative(t *testing.T) {
	intGen := lawtest.IntGen(1, 10)

	setGen := func() Set[int] {
		s := NewSet[int]()
		n := lawtest.IntGen(1, 5)()
		for i := 0; i < n; i++ {
			s = s.Add(intGen())
		}
		return s
	}

	t.Run("Union is Associative - Manual Test", func(t *testing.T) {
		// Manual test since Set[int] is not comparable
		for i := 0; i < 100; i++ {
			a, b, c := setGen(), setGen(), setGen()

			// (a ∪ b) ∪ c
			left := a.Union(b).Union(c)

			// a ∪ (b ∪ c)
			right := a.Union(b.Union(c))

			// Check if they have same elements
			if !setsEqual(left, right) {
				t.Errorf("Associativity failed for set union")
			}
		}
	})
}

// Helper to compare sets
func setsEqual(a, b Set[int]) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if _, ok := b[k]; !ok {
			return false
		}
	}
	return true
}
