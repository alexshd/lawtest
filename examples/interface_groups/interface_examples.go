// Package main demonstrates interface-based group theory testing
package main

import (
	"math/rand"
	"strings"
)

// IntModGroup implements lawtest.Group[int] for integer addition modulo n
type IntModGroup struct {
	modulo int
}

func NewIntModGroup(n int) *IntModGroup {
	return &IntModGroup{modulo: n}
}

func (g *IntModGroup) Op(a, b int) int {
	return (a + b) % g.modulo
}

func (g *IntModGroup) Identity() int {
	return 0
}

func (g *IntModGroup) Inverse(a int) int {
	return (g.modulo - a) % g.modulo
}

func (g *IntModGroup) Gen() int {
	return rand.Intn(g.modulo)
}

// StringConcatMonoid implements lawtest.Monoid[string] for string concatenation
type StringConcatMonoid struct {
	maxLen int
}

func NewStringConcatMonoid(maxLen int) *StringConcatMonoid {
	return &StringConcatMonoid{maxLen: maxLen}
}

func (m *StringConcatMonoid) Op(a, b string) string {
	return a + b
}

func (m *StringConcatMonoid) Identity() string {
	return ""
}

func (m *StringConcatMonoid) Gen() string {
	const charset = "abc"
	n := rand.Intn(m.maxLen) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// TextCleanerOp implements lawtest.IdempotentOp[string] - proper idempotent cleaning
type TextCleanerOp struct{}

func NewTextCleanerOp() *TextCleanerOp {
	return &TextCleanerOp{}
}

func (o *TextCleanerOp) Apply(x string) string {
	x = strings.TrimSpace(x)
	x = strings.ToLower(x)
	for strings.Contains(x, "  ") {
		x = strings.ReplaceAll(x, "  ", " ")
	}
	return x
}

func (o *TextCleanerOp) Gen() string {
	examples := []string{
		"  Hello World  ",
		"UPPERCASE",
		"  multiple   spaces  ",
		"MiXeD CaSe",
		"",
		"   ",
	}
	return examples[rand.Intn(len(examples))]
}

// BrokenTextCleanerOp - NOT idempotent (demonstrates http_server.go bug)
type BrokenTextCleanerOp struct{}

func NewBrokenTextCleanerOp() *BrokenTextCleanerOp {
	return &BrokenTextCleanerOp{}
}

func (o *BrokenTextCleanerOp) Apply(x string) string {
	x = strings.ReplaceAll(x, "\"", "\\\"")
	return x
}

func (o *BrokenTextCleanerOp) Gen() string {
	examples := []string{
		`"test"`,
		`hello`,
		`\"already escaped\"`,
	}
	return examples[rand.Intn(len(examples))]
}

// BoolAndSemigroup implements lawtest.Semigroup[bool] for boolean AND
type BoolAndSemigroup struct{}

func NewBoolAndSemigroup() *BoolAndSemigroup {
	return &BoolAndSemigroup{}
}

func (s *BoolAndSemigroup) Op(a, b bool) bool {
	return a && b
}

func (s *BoolAndSemigroup) Gen() bool {
	return rand.Intn(2) == 1
}

// BrokenSubtractionGroup - NOT a valid group (subtraction not associative)
type BrokenSubtractionGroup struct {
	maxValue int
}

func NewBrokenSubtractionGroup(max int) *BrokenSubtractionGroup {
	return &BrokenSubtractionGroup{maxValue: max}
}

func (g *BrokenSubtractionGroup) Op(a, b int) int {
	return a - b
}

func (g *BrokenSubtractionGroup) Identity() int {
	return 0
}

func (g *BrokenSubtractionGroup) Inverse(a int) int {
	return -a
}

func (g *BrokenSubtractionGroup) Gen() int {
	return rand.Intn(g.maxValue*2) - g.maxValue
}

// Matrix2x2 type for matrix operations
type Matrix2x2 struct {
	A, B, C, D int
}

// Matrix2x2Group implements group operations for 2x2 matrices
type Matrix2x2Group struct {
	modulo int
}

func NewMatrix2x2Group(mod int) *Matrix2x2Group {
	return &Matrix2x2Group{modulo: mod}
}

func (g *Matrix2x2Group) Op(m1, m2 Matrix2x2) Matrix2x2 {
	return Matrix2x2{
		A: (m1.A*m2.A + m1.B*m2.C) % g.modulo,
		B: (m1.A*m2.B + m1.B*m2.D) % g.modulo,
		C: (m1.C*m2.A + m1.D*m2.C) % g.modulo,
		D: (m1.C*m2.B + m1.D*m2.D) % g.modulo,
	}
}

func (g *Matrix2x2Group) Identity() Matrix2x2 {
	return Matrix2x2{A: 1, B: 0, C: 0, D: 1}
}

func (g *Matrix2x2Group) Inverse(m Matrix2x2) Matrix2x2 {
	return g.Identity()
}

func (g *Matrix2x2Group) Gen() Matrix2x2 {
	return Matrix2x2{
		A: rand.Intn(g.modulo),
		B: rand.Intn(g.modulo),
		C: rand.Intn(g.modulo),
		D: rand.Intn(g.modulo),
	}
}
