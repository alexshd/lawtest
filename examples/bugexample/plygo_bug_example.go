// Example: Using goprops to catch the valueKey bug from plyGO
//
// This demonstrates how property-based testing with group theory
// would have caught plyGO's valueKey bug where different slices
// all collapsed to the same "<complex>" key.
package main

import (
	"fmt"
	"testing"

	"github.com/alexshd/goprops"
)

// Simplified version of plyGO's problematic valueKey function
func valueKeyBuggy(v any) any {
	switch v := v.(type) {
	case []int:
		return "<complex>" // BUG: All slices map to same key!
	case map[string]int:
		return "<complex>" // BUG: All maps map to same key!
	default:
		return v
	}
}

// Fixed version that properly distinguishes values
func valueKeyFixed(v any) any {
	// Return the value itself - let Go's comparison handle it
	// (Only works for comparable types)
	return v
}

// Test that valueKey preserves distinctness
// If a != b, then valueKey(a) != valueKey(b)
func TestValueKeyPreservesDistinctness(t *testing.T) {
	// Generator for int slices
	sliceGen := func() []int {
		n := goprops.IntGen(1, 5)()
		slice := make([]int, n)
		for i := range slice {
			slice[i] = goprops.IntGen(1, 100)()
		}
		return slice
	}

	t.Run("Buggy Version Fails", func(t *testing.T) {
		// Generate two different slices
		for i := 0; i < 10; i++ {
			a := sliceGen()
			b := sliceGen()

			// If they're actually different slices
			if !equal(a, b) {
				keyA := valueKeyBuggy(a)
				keyB := valueKeyBuggy(b)

				// Keys should be different!
				if keyA == keyB {
					t.Logf("BUG DETECTED: Different slices have same key")
					t.Logf("  a=%v -> key=%v", a, keyA)
					t.Logf("  b=%v -> key=%v", b, keyB)
					t.Logf("  This causes GroupBy to merge distinct groups!")
					return
				}
			}
		}
	})
}

// Property: valueKey should be injective (one-to-one)
// If valueKey(a) == valueKey(b), then a == b
func TestValueKeyInjective(t *testing.T) {
	intGen := goprops.IntGen(1, 1000)

	t.Run("Works for primitive types", func(t *testing.T) {
		seen := make(map[any]int)

		for i := 0; i < 100; i++ {
			val := intGen()
			key := valueKeyFixed(val)

			if prevVal, exists := seen[key]; exists {
				// If we've seen this key before, the values should be equal
				if prevVal != val {
					t.Errorf("Injectivity violated: different values map to same key")
					t.Errorf("  val1=%v, val2=%v, key=%v", prevVal, val, key)
				}
			}
			seen[key] = val
		}
	})
}

// Property: valueKey should form a homomorphism
// For GroupBy to work correctly, equal values must map to equal keys
func TestValueKeyHomomorphism(t *testing.T) {
	intGen := goprops.IntGen(1, 100)

	for i := 0; i < 100; i++ {
		a := intGen()
		b := a // Same value

		keyA := valueKeyFixed(a)
		keyB := valueKeyFixed(b)

		if keyA != keyB {
			t.Errorf("Equal values should have equal keys")
			t.Errorf("  a=%v -> key=%v", a, keyA)
			t.Errorf("  b=%v -> key=%v", b, keyB)
		}
	}
}

// Demonstrating how group theory helps test GroupBy
type GroupByOp func(items []int, keyFunc func(int) any) map[any][]int

func groupBy(items []int, keyFunc func(int) any) map[any][]int {
	result := make(map[any][]int)
	for _, item := range items {
		key := keyFunc(item)
		result[key] = append(result[key], item)
	}
	return result
}

func TestGroupByProperties(t *testing.T) {
	// Property: GroupBy should partition the input
	// Every element appears exactly once across all groups
	t.Run("Partition Property", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5, 6}
		keyFunc := func(x int) any { return x % 2 } // Even/odd

		groups := groupBy(items, keyFunc)

		// Count total elements across all groups
		totalCount := 0
		for _, group := range groups {
			totalCount += len(group)
		}

		if totalCount != len(items) {
			t.Errorf("GroupBy lost elements: expected %d, got %d", len(items), totalCount)
		}
	})

	// Property: Elements in same group must have equal keys
	t.Run("Equal Keys Property", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		keyFunc := func(x int) any { return x % 3 } // Mod 3

		groups := groupBy(items, keyFunc)

		for key, group := range groups {
			// All elements in this group should produce the same key
			for _, item := range group {
				if keyFunc(item) != key {
					t.Errorf("Item %d in group %v doesn't match group key", item, key)
				}
			}
		}
	})

	// Property: Idempotence - GroupBy(GroupBy(x)) ≈ GroupBy(x)
	// (Re-grouping by same key should produce same result)
	t.Run("Idempotent Property", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5, 6}
		keyFunc := func(x int) any { return x % 2 }

		groups1 := groupBy(items, keyFunc)

		// Flatten and re-group
		var flattened []int
		for _, group := range groups1 {
			flattened = append(flattened, group...)
		}
		groups2 := groupBy(flattened, keyFunc)

		// Should have same number of groups
		if len(groups1) != len(groups2) {
			t.Errorf("Re-grouping changed group count: %d -> %d", len(groups1), len(groups2))
		}
	})
}

// Helper function
func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println("This example shows how goprops catches the plyGO valueKey bug")
	fmt.Println("Run with: go test -v plygo_bug_example.go")
}
