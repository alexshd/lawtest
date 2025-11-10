// Package tailrecursion demonstrates how to use lawtest to verify
// tail call optimizations produce equivalent results to recursive implementations.
package tailrecursion

// Factorial computes n! using standard recursion.
// This builds up the call stack: 5 * (4 * (3 * (2 * 1)))
func Factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * Factorial(n-1)
}

// FactorialTail computes n! using tail recursion with an accumulator.
// The recursive call is the last operation (tail position).
// Compilers can optimize this to use constant stack space.
func FactorialTail(n, acc int) int {
	if n <= 1 {
		return acc
	}
	return FactorialTail(n-1, n*acc)
}

// Sum computes the sum of integers from 1 to n using recursion.
func Sum(n int) int {
	if n <= 0 {
		return 0
	}
	return n + Sum(n-1)
}

// SumTail computes the sum using tail recursion with an accumulator.
func SumTail(n, acc int) int {
	if n <= 0 {
		return acc
	}
	return SumTail(n-1, acc+n)
}

// ReverseList reverses a slice using recursion.
// Creates new slices at each level.
func ReverseList(list []int) []int {
	if len(list) == 0 {
		return []int{}
	}
	return append(ReverseList(list[1:]), list[0])
}

// ReverseListIterative reverses a slice using iteration.
// Single allocation, constant stack space.
func ReverseListIterative(list []int) []int {
	result := make([]int, len(list))
	for i := range list {
		result[len(list)-1-i] = list[i]
	}
	return result
}

// Fibonacci computes the nth Fibonacci number using naive recursion.
// Exponential time complexity O(2^n) - very slow for large n.
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// FibonacciIterative computes Fibonacci using iteration.
// Linear time complexity O(n) - much faster.
func FibonacciIterative(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// Power computes base^exp using recursion.
func Power(base, exp int) int {
	if exp == 0 {
		return 1
	}
	return base * Power(base, exp-1)
}

// PowerTail computes base^exp using tail recursion.
func PowerTail(base, exp, acc int) int {
	if exp == 0 {
		return acc
	}
	return PowerTail(base, exp-1, acc*base)
}
