package tailrecursion

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/alexshd/lawtest"
)

// TestFactorialEquivalence proves that tail recursive factorial
// produces the same results as standard recursive factorial.
func TestFactorialEquivalence(t *testing.T) {
	gen := func() int { return rand.Intn(15) + 1 } // Keep small to avoid overflow

	lawtest.Equivalent(t,
		func(n int) int { return Factorial(n) },
		func(n int) int { return FactorialTail(n, 1) },
		gen,
	)
}

// TestSumEquivalence proves tail recursive sum is equivalent to standard recursion.
func TestSumEquivalence(t *testing.T) {
	gen := func() int { return rand.Intn(100) }

	lawtest.Equivalent(t,
		func(n int) int { return Sum(n) },
		func(n int) int { return SumTail(n, 0) },
		gen,
	)
}

// TestReverseListEquivalence proves iterative reverse is equivalent to recursive.
func TestReverseListEquivalence(t *testing.T) {
	gen := func() []int {
		n := rand.Intn(20) + 1
		list := make([]int, n)
		for i := range list {
			list[i] = rand.Intn(100)
		}
		return list
	}

	eq := func(a, b []int) bool {
		return reflect.DeepEqual(a, b)
	}

	lawtest.EquivalentCustom(t, ReverseList, ReverseListIterative, gen, eq)
}

// TestFibonacciEquivalence proves iterative Fibonacci is equivalent to recursive.
func TestFibonacciEquivalence(t *testing.T) {
	gen := func() int { return rand.Intn(20) } // Keep small for naive recursive version

	lawtest.Equivalent(t, Fibonacci, FibonacciIterative, gen)
}

// TestFibonacciIterEquivalence proves iterator-based Fibonacci is equivalent.
func TestFibonacciIterEquivalence(t *testing.T) {
	gen := func() int { return rand.Intn(20) }

	lawtest.Equivalent(t, Fibonacci, FibonacciIter, gen)
}

// TestPowerEquivalence proves tail recursive power is equivalent to standard recursion.
func TestPowerEquivalence(t *testing.T) {
	gen := func() struct{ base, exp int } {
		return struct{ base, exp int }{
			base: rand.Intn(5) + 1,
			exp:  rand.Intn(10),
		}
	}

	lawtest.Equivalent(t,
		func(p struct{ base, exp int }) int { return Power(p.base, p.exp) },
		func(p struct{ base, exp int }) int { return PowerTail(p.base, p.exp, 1) },
		gen,
	)
}

// BenchmarkFactorial benchmarks standard recursive factorial.
func BenchmarkFactorial(b *testing.B) {
	for b.Loop() {
		Factorial(15)
	}
}

// BenchmarkFactorialTail benchmarks tail recursive factorial.
func BenchmarkFactorialTail(b *testing.B) {
	for b.Loop() {
		FactorialTail(15, 1)
	}
}

// BenchmarkFibonacci benchmarks naive recursive Fibonacci.
func BenchmarkFibonacci(b *testing.B) {
	for b.Loop() {
		Fibonacci(20)
	}
}

// BenchmarkFibonacciIterative benchmarks iterative Fibonacci.
func BenchmarkFibonacciIterative(b *testing.B) {
	for b.Loop() {
		FibonacciIterative(50)
	}
}

// BenchmarkFibonacciIter benchmarks iterator-based Fibonacci.
func BenchmarkFibonacciIter(b *testing.B) {
	for b.Loop() {
		FibonacciIter(50)
	}
}
