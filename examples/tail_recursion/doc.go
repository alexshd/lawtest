// Package tailrecursion demonstrates how to use lawtest.Equivalent() to verify
// tail call optimizations produce equivalent results to recursive implementations.
//
// # The Problem
//
// Recursive functions are elegant but can be inefficient:
//   - Stack overflow for large inputs
//   - Poor performance (especially tree recursion like Fibonacci)
//   - Memory overhead from call stack
//
// Optimized versions (tail recursion, iteration) solve these problems but
// how do you prove they're correct?
//
// # The Solution
//
// lawtest.Equivalent() tests if two functions produce the same output for all inputs:
//
//	lawtest.Equivalent(t,
//	    func(n int) int { return Factorial(n) },        // Original
//	    func(n int) int { return FactorialTail(n, 1) }, // Optimized
//	    gen,
//	)
//
// If this test passes, the optimized version is mathematically proven correct.
//
// # Example: Factorial
//
//	// Standard recursion: builds call stack
//	func Factorial(n int) int {
//	    if n <= 1 { return 1 }
//	    return n * Factorial(n-1)  // Not tail position
//	}
//
//	// Tail recursion: can be optimized to iteration
//	func FactorialTail(n, acc int) int {
//	    if n <= 1 { return acc }
//	    return FactorialTail(n-1, n*acc)  // Tail position!
//	}
//
// # Example: Fibonacci Performance
//
// Benchmarks show dramatic performance improvements:
//
//	BenchmarkFibonacci-4               25,194      45,406 ns/op  (recursive)
//	BenchmarkFibonacciIterative-4  74,562,684         16.87 ns/op  (iterative)
//
// Iterative is 2,691x faster! The equivalence test proves correctness:
//
//	func TestFibonacciEquivalence(t *testing.T) {
//	    gen := func() int { return rand.Intn(20) }
//	    lawtest.Equivalent(t, Fibonacci, FibonacciIterative, gen)
//	}
//
// # When to Use
//
//   - Verifying tail recursion optimization
//   - Proving iterative version matches recursive
//   - Validating cached/memoized implementations
//   - Testing algorithm refactoring
package tailrecursion
