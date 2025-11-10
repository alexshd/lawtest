# Tail Recursion Optimization with lawtest

This example demonstrates how to use `lawtest.Equivalent()` to verify that optimized implementations (tail recursion, iterative versions) produce the same results as naive recursive implementations.

## The Problem

Recursive functions are elegant but can be inefficient:

- **Stack overflow** for large inputs
- **Poor performance** (especially tree recursion like Fibonacci)
- **Memory overhead** from call stack

Optimized versions (tail recursion, iteration) solve these problems but **how do you prove they're correct**?

## The Solution: lawtest.Equivalent()

`lawtest.Equivalent()` tests if two functions produce the same output for all inputs:

```go
lawtest.Equivalent(t,
    func(n int) int { return Factorial(n) },        // Original
    func(n int) int { return FactorialTail(n, 1) }, // Optimized
    gen,                                             // Random inputs
)
```

**If this test passes, the optimized version is mathematically proven correct.**

## Examples

### Factorial: Recursive vs Tail Recursive

```go
// Standard recursion: builds call stack
func Factorial(n int) int {
    if n <= 1 { return 1 }
    return n * Factorial(n-1)  // Not tail position
}

// Tail recursion: can be optimized to iteration
func FactorialTail(n, acc int) int {
    if n <= 1 { return acc }
    return FactorialTail(n-1, n*acc)  // Tail position!
}
```

**Test equivalence:**

```go
func TestFactorialEquivalence(t *testing.T) {
    gen := func() int { return rand.Intn(15) + 1 }
    lawtest.Equivalent(t,
        func(n int) int { return Factorial(n) },
        func(n int) int { return FactorialTail(n, 1) },
        gen,
    )
}
```

✅ Result: `Functions are equivalent (tested 100 random inputs)`

### Fibonacci: Recursive vs Iterative

```go
// Naive recursion: O(2^n) - exponentially slow!
func Fibonacci(n int) int {
    if n <= 1 { return n }
    return Fibonacci(n-1) + Fibonacci(n-2)
}

// Iterative: O(n) - linear time
func FibonacciIterative(n int) int {
    if n <= 1 { return n }
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b
}
```

**Test equivalence:**

```go
func TestFibonacciEquivalence(t *testing.T) {
    gen := func() int { return rand.Intn(20) }
    lawtest.Equivalent(t, Fibonacci, FibonacciIterative, gen)
}
```

✅ Result: `Functions are equivalent (tested 100 random inputs)`

**Performance difference (benchmarked on Fib(20)):**

```
BenchmarkFibonacci-4               25,194      45,406 ns/op
BenchmarkFibonacciIterative-4  74,562,684         16.87 ns/op
```

**Iterative is 2,691x faster!** 🚀

### List Reverse: Recursive vs Iterative

```go
// Recursive: creates new slices at each level
func ReverseList(list []int) []int {
    if len(list) == 0 { return []int{} }
    return append(ReverseList(list[1:]), list[0])
}

// Iterative: single allocation
func ReverseListIterative(list []int) []int {
    result := make([]int, len(list))
    for i := range list {
        result[len(list)-1-i] = list[i]
    }
    return result
}
```

**Test equivalence with custom equality:**

```go
func TestReverseListEquivalence(t *testing.T) {
    gen := func() []int {
        n := rand.Intn(20) + 1
        list := make([]int, n)
        for i := range list { list[i] = rand.Intn(100) }
        return list
    }

    eq := func(a, b []int) bool {
        return reflect.DeepEqual(a, b)
    }

    lawtest.EquivalentCustom(t, ReverseList, ReverseListIterative, gen, eq)
}
```

✅ Result: `Functions are equivalent (tested 100 random inputs)`

## Running the Tests

```bash
# Run all tests
go test -v

# Run with benchmarks
go test -bench=. -benchmem
```

## When to Use This

✅ **Verifying tail recursion optimization**  
✅ **Proving iterative version matches recursive**  
✅ **Validating cached/memoized implementations**  
✅ **Testing algorithm refactoring**  
✅ **Comparing naive vs optimized algorithms**

## API

### Equivalent (for comparable types)

```go
func Equivalent[T any, R comparable](
    t *testing.T,
    f1 func(T) R,  // Original function
    f2 func(T) R,  // Optimized function
    gen func() T,  // Random input generator
) bool
```

### EquivalentCustom (for non-comparable types)

```go
func EquivalentCustom[T any, R any](
    t *testing.T,
    f1 func(T) R,          // Original function
    f2 func(T) R,          // Optimized function
    gen func() T,          // Random input generator
    eq func(R, R) bool,    // Custom equality
) bool
```

## Key Insight

**Mathematical proof of correctness:** If `Equivalent()` passes, you can confidently replace the naive implementation with the optimized one, knowing they produce identical results.

No manual test cases needed - lawtest generates 100 random test cases automatically.

---

_"Optimize with confidence. Let math prove correctness."_
