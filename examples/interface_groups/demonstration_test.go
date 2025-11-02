package main

import (
	"testing"

	"github.com/alexshd/lawtest"
)

// ============================================================================
// CONCURRENCY SAFETY TESTS - User-Defined Interface Testing
// ============================================================================
//
// This demonstrates testing concurrency properties of user-defined interfaces.
// The key insight: Immutability enables parallelism.

// TestGoodCacheParallelSafety tests if GoodCache.Merge is safe for concurrent use
func TestGoodCacheParallelSafety(t *testing.T) {
	t.Log("=== Testing GoodCache Parallel Safety ===")

	op := func(a, b *GoodCache) *GoodCache {
		return a.Merge(b).(*GoodCache)
	}

	gen := func() *GoodCache {
		cache := NewGoodCache()
		cache.Set("x", 1)
		cache.Set("y", 2)
		return cache
	}

	// Test with 20 concurrent goroutines
	isSafe := lawtest.ParallelSafe(t, op, gen, 20)

	if isSafe {
		t.Log("✅ GoodCache.Merge is parallel-safe (immutable operations)")
	} else {
		t.Error("❌ GoodCache.Merge should be parallel-safe but detected issues")
	}
}

// TestBrokenCacheParallelSafety tests if BrokenCache.Merge has race conditions
func TestBrokenCacheParallelSafety(t *testing.T) {
	t.Log("=== Testing BrokenCache Parallel Safety ===")
	t.Log("⚠️  Note: This test may exhibit non-deterministic behavior due to race conditions")

	op := func(a, b *BrokenCache) *BrokenCache {
		return a.Merge(b).(*BrokenCache)
	}

	gen := func() *BrokenCache {
		cache := NewBrokenCache()
		cache.Set("x", 1)
		cache.Set("y", 2)
		return cache
	}

	// Test with 20 concurrent goroutines
	// We expect this might fail due to mutation
	t.Run("RaceDetection", func(t *testing.T) {
		isSafe := lawtest.ParallelSafe(t, op, gen, 20)

		if !isSafe {
			t.Log("✅ Correctly detected that BrokenCache is NOT parallel-safe")
		} else {
			t.Log("⚠️  BrokenCache appears safe but has mutation bugs (may need -race flag)")
		}
	})
}

// TestGoodCacheImmutability verifies that GoodCache operations don't mutate inputs
func TestGoodCacheImmutability(t *testing.T) {
	t.Log("=== Testing GoodCache Immutability ===")

	// For pointer types, we need a different approach
	// Let's manually verify immutability
	cache1 := NewGoodCache()
	cache1.Set("x", 1)

	cache2 := NewGoodCache()
	cache2.Set("y", 2)

	// Store original values
	val1Before, _ := cache1.Get("x")
	val2Before, _ := cache2.Get("y")
	len1Before := len(cache1.data)
	len2Before := len(cache2.data)

	// Perform merge
	result := cache1.Merge(cache2)

	// Check that originals weren't modified
	val1After, _ := cache1.Get("x")
	val2After, _ := cache2.Get("y")
	len1After := len(cache1.data)
	len2After := len(cache2.data)

	if val1Before != val1After || len1Before != len1After {
		t.Errorf("❌ cache1 was mutated! Before: x=%d len=%d, After: x=%d len=%d",
			val1Before, len1Before, val1After, len1After)
	} else {
		t.Log("✅ cache1 not mutated (immutable)")
	}

	if val2Before != val2After || len2Before != len2After {
		t.Errorf("❌ cache2 was mutated! Before: y=%d len=%d, After: y=%d len=%d",
			val2Before, len2Before, val2After, len2After)
	} else {
		t.Log("✅ cache2 not mutated (immutable)")
	}

	// Verify result is a new instance
	if result == cache1 || result == cache2 {
		t.Error("❌ Result is same instance as input (not immutable)")
	} else {
		t.Log("✅ Result is new instance (immutable merge)")
	}
}

// TestBrokenCacheImmutability verifies that BrokenCache DOES mutate (demonstrates bug)
func TestBrokenCacheImmutability(t *testing.T) {
	t.Log("=== Testing BrokenCache Immutability (Expect FAILURE) ===")

	cache1 := NewBrokenCache()
	cache1.Set("x", 1)

	cache2 := NewBrokenCache()
	cache2.Set("x", 2)

	// Store original values
	len1Before := len(cache1.data)

	// Perform merge
	result := cache1.Merge(cache2)

	// Check if cache1 was modified
	len1After := len(cache1.data)

	if len1Before != len1After {
		t.Logf("✅ Correctly detected mutation! cache1 length changed: %d → %d", len1Before, len1After)
		t.Log("   BrokenCache violates immutability (mutates receiver)")
	} else {
		t.Log("⚠️  Mutation not detected in this run (may depend on data)")
	}

	// Verify result is same instance (bug!)
	if result == cache1 {
		t.Log("✅ Correctly detected: result is SAME instance as input (violates immutability)")
	} else {
		t.Error("❌ Expected result to be same instance (demonstrating bug)")
	}
}

// ============================================================================
// KEY INSIGHTS
// ============================================================================
//
// 1. USER-DEFINED INTERFACE: Cacher is NOT from lawtest - you defined it!
//
// 2. TWO IMPLEMENTATIONS:
//    ✅ GoodCache: Immutable merge, last-write-wins → passes all tests
//    ❌ BrokenCache: Mutating merge, first-write-wins → fails associativity
//
// 3. LAWTEST CATCHES THE BUG:
//    - No need to write specific test cases
//    - Property tests automatically find violations
//    - Mathematical laws reveal implementation bugs
//
// 4. REAL BENEFIT:
//    - Define your interface
//    - Implement it
//    - Let lawtest verify it satisfies algebraic properties
//    - Catch subtle bugs (mutation, ordering, side effects)
//
// Run: go test -v
// See: GoodCache ✓ passes, BrokenCache ✗ fails
//
// This is how you use lawtest with YOUR custom interfaces!
