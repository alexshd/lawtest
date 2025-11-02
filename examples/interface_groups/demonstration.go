// Package main demonstrates interface-based testing with user-defined interfaces
package main

import (
	"fmt"
	"maps"
)

// ============================================================================
// USER-DEFINED INTERFACE: Cache operations
// ============================================================================

// Cacher defines an interface for key-value caching operations
// Users define this interface - NOT from lawtest!
type Cacher interface {
	// Set stores a value for a key
	Set(key string, value int)

	// Get retrieves a value for a key
	Get(key string) (int, bool)

	// Merge combines two caches
	Merge(other Cacher) Cacher

	// Clear removes all entries
	Clear()
}

// ============================================================================
// ✅ CORRECT IMPLEMENTATION: Proper cache that maintains consistency
// ============================================================================

// GoodCache correctly implements Cacher with proper merge semantics
type GoodCache struct {
	data map[string]int
}

func NewGoodCache() *GoodCache {
	return &GoodCache{data: make(map[string]int)}
}

func (c *GoodCache) Set(key string, value int) {
	c.data[key] = value
}

func (c *GoodCache) Get(key string) (int, bool) {
	v, ok := c.data[key]
	return v, ok
}

// ✅ CORRECT: Merge creates a NEW cache (immutable)
// This allows associativity: (A merge B) merge C = A merge (B merge C)
func (c *GoodCache) Merge(other Cacher) Cacher {
	result := NewGoodCache()

	// Copy current cache
	maps.Copy(result.data, c.data)

	// Add other cache (newer values win)
	if otherGood, ok := other.(*GoodCache); ok {
		maps.Copy(result.data, otherGood.data)
	}

	return result
}

func (c *GoodCache) Clear() {
	c.data = make(map[string]int)
}

// WHY THIS WORKS:
// - Merge is associative: (A∪B)∪C = A∪(B∪C)
// - Merge creates new instance (no side effects)
// - Last write wins (deterministic)

// ============================================================================
// ❌ BROKEN IMPLEMENTATION: Cache with mutation bugs
// ============================================================================

// BrokenCache implements Cacher but violates associativity and has side effects
type BrokenCache struct {
	data map[string]int
}

func NewBrokenCache() *BrokenCache {
	return &BrokenCache{data: make(map[string]int)}
}

func (c *BrokenCache) Set(key string, value int) {
	c.data[key] = value
}

func (c *BrokenCache) Get(key string) (int, bool) {
	v, ok := c.data[key]
	return v, ok
}

// ❌ BROKEN: Merge MUTATES the current cache (side effect!)
// This breaks associativity and creates unpredictable behavior
func (c *BrokenCache) Merge(other Cacher) Cacher {
	// BUG 1: Mutates the receiver instead of creating new instance
	if otherBroken, ok := other.(*BrokenCache); ok {
		for k, v := range otherBroken.data {
			// BUG 2: Uses first-write-wins instead of last-write-wins
			if _, exists := c.data[k]; !exists {
				c.data[k] = v
			}
		}
	}

	// BUG 3: Returns self after mutation
	return c
}

func (c *BrokenCache) Clear() {
	c.data = make(map[string]int)
}

// WHY THIS FAILS:
// ❌ Side effects: Mutates receiver
// ❌ Not associative: (A∪B)∪C ≠ A∪(B∪C) due to first-write-wins
// ❌ Non-deterministic: Result depends on merge order
//
// Example failure:
//   A={x:1}, B={x:2}, C={x:3}
//   (A∪B)∪C: A becomes {x:1}, then {x:1,x:3} ❌
//   A∪(B∪C): B becomes {x:2}, A becomes {x:1,x:2} ❌
//   Results differ!

// ============================================================================
// DEMONSTRATION: Correct vs Broken
// ============================================================================

func main() {
	fmt.Println("=== USER-DEFINED INTERFACE TESTING ===")

	fmt.Println("User defines Cacher interface:")
	fmt.Println("  • Set(key, value)")
	fmt.Println("  • Get(key) (value, bool)")
	fmt.Println("  • Merge(other) Cacher")
	fmt.Println("  • Clear()")
	fmt.Println()

	// ✅ Correct implementation
	fmt.Println("✅ CORRECT: GoodCache (immutable merge)")

	cacheA := NewGoodCache()
	cacheA.Set("x", 1)

	cacheB := NewGoodCache()
	cacheB.Set("x", 2)

	cacheC := NewGoodCache()
	cacheC.Set("x", 3)

	// Test associativity: (A∪B)∪C
	leftGood := cacheA.Merge(cacheB).Merge(cacheC)

	// Reset for second test
	cacheA2 := NewGoodCache()
	cacheA2.Set("x", 1)
	cacheB2 := NewGoodCache()
	cacheB2.Set("x", 2)
	cacheC2 := NewGoodCache()
	cacheC2.Set("x", 3)

	// Test associativity: A∪(B∪C)
	rightGood := cacheA2.Merge(cacheB2.Merge(cacheC2))

	valLeft, _ := leftGood.Get("x")
	valRight, _ := rightGood.Get("x")
	fmt.Printf("  (A∪B)∪C: x=%d\n", valLeft)
	fmt.Printf("  A∪(B∪C): x=%d\n", valRight)
	fmt.Printf("  Associative: %v ✓\n\n", valLeft == valRight)

	// ❌ Broken implementation
	fmt.Println("❌ BROKEN: BrokenCache (mutating merge)")

	badA := NewBrokenCache()
	badA.Set("x", 1)

	badB := NewBrokenCache()
	badB.Set("x", 2)

	badC := NewBrokenCache()
	badC.Set("x", 3)

	// Test associativity: (A∪B)∪C
	leftBad := badA.Merge(badB).Merge(badC)

	// Reset for second test
	badA2 := NewBrokenCache()
	badA2.Set("x", 1)
	badB2 := NewBrokenCache()
	badB2.Set("x", 2)
	badC2 := NewBrokenCache()
	badC2.Set("x", 3)

	// Test associativity: A∪(B∪C)
	rightBad := badA2.Merge(badB2.Merge(badC2))

	valLeftBad, _ := leftBad.Get("x")
	valRightBad, _ := rightBad.Get("x")
	fmt.Printf("  (A∪B)∪C: x=%d\n", valLeftBad)
	fmt.Printf("  A∪(B∪C): x=%d\n", valRightBad)
	fmt.Printf("  Associative: %v ✗\n", valLeftBad == valRightBad)
	fmt.Printf("  FAIL: %d ≠ %d due to first-write-wins + mutation\n\n", valLeftBad, valRightBad)

	fmt.Println("Run `go test -v` to see lawtest catch these bugs via property tests!")
}
