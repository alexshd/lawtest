// Package main tests interface-based group theory implementations
package main

import (
	"fmt"
	"testing"

	"github.com/alexshd/lawtest"
)

// Test IntModGroup satisfies all group properties
func TestIntModGroup_IsValidGroup(t *testing.T) {
	g := NewIntModGroup(12)
	lawtest.TestGroup[int](t, g)
}

// Test with different modulos
func TestIntModGroup_DifferentModulos(t *testing.T) {
	modulos := []int{5, 7, 10, 13, 100}

	for _, m := range modulos {
		t.Run(fmt.Sprintf("Modulo_%d", m), func(t *testing.T) {
			g := NewIntModGroup(m)
			lawtest.TestGroup[int](t, g)
		})
	}
}

// Test StringConcatMonoid is a valid monoid
func TestStringConcatMonoid_IsValidMonoid(t *testing.T) {
	m := NewStringConcatMonoid(5)
	lawtest.TestMonoid[string](t, m)
}

// Test TextCleanerOp is idempotent
func TestTextCleanerOp_IsIdempotent(t *testing.T) {
	op := NewTextCleanerOp()
	lawtest.TestIdempotentOp[string](t, op)
}

// Test BrokenTextCleanerOp demonstrates idempotence violation
func TestBrokenTextCleanerOp_DemonstratesBug(t *testing.T) {
	op := NewBrokenTextCleanerOp()

	t.Run("ManualDemo", func(t *testing.T) {
		input := `"test"`
		r1 := op.Apply(input)
		r2 := op.Apply(r1)
		r3 := op.Apply(r2)

		t.Logf("Idempotence violation (http_server.go bug):")
		t.Logf("  Input:     %q", input)
		t.Logf("  Apply(x):  %q", r1)
		t.Logf("  Apply²(x): %q", r2)
		t.Logf("  Apply³(x): %q", r3)

		if r1 == r2 {
			t.Error("BUG: Should NOT be idempotent, but it is!")
		}
	})
}

// Test BoolAndSemigroup is valid
func TestBoolAndSemigroup_IsValidSemigroup(t *testing.T) {
	s := NewBoolAndSemigroup()
	lawtest.TestSemigroup[bool](t, s)
}

// Test BrokenSubtractionGroup fails associativity
func TestBrokenSubtractionGroup_FailsAssociativity(t *testing.T) {
	g := NewBrokenSubtractionGroup(10)

	t.Run("ManualDemo", func(t *testing.T) {
		a, b, c := 10, 5, 3

		left := g.Op(g.Op(a, b), c)
		right := g.Op(a, g.Op(b, c))

		t.Logf("(10 - 5) - 3 = %d", left)
		t.Logf("10 - (5 - 3) = %d", right)

		if left == right {
			t.Error("BUG: Subtraction should NOT be associative!")
		}
	})
}

// Test Matrix2x2Group properties
func TestMatrix2x2Group_Properties(t *testing.T) {
	g := NewMatrix2x2Group(7)

	t.Run("Associativity", func(t *testing.T) {
		lawtest.Associative(t, g.Op, g.Gen)
	})

	t.Run("Identity", func(t *testing.T) {
		lawtest.Identity(t, g.Op, g.Identity(), g.Gen)
	})
}

// Comprehensive test running all valid implementations
func TestAllValidImplementations(t *testing.T) {
	t.Run("IntModGroup", func(t *testing.T) {
		g := NewIntModGroup(12)
		lawtest.TestGroup[int](t, g)
	})

	t.Run("StringConcatMonoid", func(t *testing.T) {
		m := NewStringConcatMonoid(5)
		lawtest.TestMonoid[string](t, m)
	})

	t.Run("TextCleanerOp", func(t *testing.T) {
		op := NewTextCleanerOp()
		lawtest.TestIdempotentOp[string](t, op)
	})

	t.Run("BoolAndSemigroup", func(t *testing.T) {
		s := NewBoolAndSemigroup()
		lawtest.TestSemigroup[bool](t, s)
	})
}
