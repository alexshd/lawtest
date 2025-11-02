# Learning Guide for goprops

## What You're Learning (3-4 Day Journey)

### Day 1: Group Theory Basics

**Concepts to explore:**

1. **Binary Operation**: A function that takes two values and returns one

   - Example: `add(a, b) = a + b`
   - Example: `concat(s1, s2) = s1 + s2`

2. **Associativity**: Order of operations doesn't matter

   - `(a ∘ b) ∘ c = a ∘ (b ∘ c)`
   - Example: `(2 + 3) + 4 = 2 + (3 + 4)`
   - Why it matters: Allows parallel computation, reordering optimizations

3. **Identity Element**: A "do nothing" value
   - `a ∘ e = a` and `e ∘ a = a`
   - Example: `0` is identity for addition: `x + 0 = x`
   - Example: `1` is identity for multiplication: `x * 1 = x`
   - Example: `""` is identity for string concatenation: `s + "" = s`

**Exercises:**

- [ ] Test if subtraction is associative (it's not!)
- [ ] Find identity elements for boolean AND and OR
- [ ] Try to find identity for division (there isn't one)

### Day 2: Beyond Monoids

4. **Commutativity**: Order of arguments doesn't matter

   - `a ∘ b = b ∘ a`
   - Example: `2 + 3 = 3 + 2` (addition IS commutative)
   - Example: `"abc" + "def" ≠ "def" + "abc"` (concatenation is NOT)
   - Example: Matrix multiplication is NOT commutative

5. **Inverse Elements**: Every element has an "undo" operation

   - `a ∘ a⁻¹ = e` and `a⁻¹ ∘ a = e`
   - Example: For addition, inverse of `5` is `-5`: `5 + (-5) = 0`
   - Example: For multiplication, inverse of `2` is `0.5`: `2 * 0.5 = 1`
   - Example: For concatenation, there's no inverse (can't "unappend")

6. **Closure**: Result stays in same type
   - If `a` and `b` are integers, `a + b` is an integer
   - Go's generics enforce this at compile time

**Exercises:**

- [ ] Test if division has inverses (careful with division by zero!)
- [ ] Create a custom Point type and test vector addition properties
- [ ] Try matrix multiplication - associative but not commutative

### Day 3: Real-World Applications

**How this catches bugs:**

1. **Cache Consistency**

   - Property: Last write wins
   - Test: `cache.set(k, a); cache.set(k, b); cache.get(k) == b`
   - Bug caught: Cache returning stale values

2. **Data Deduplication** (like plyGO's GroupBy)

   - Property: Distinct values → distinct keys
   - Test: `if a != b, then key(a) != key(b)`
   - Bug caught: All slices mapping to `"<complex>"` key

3. **Distributed Systems**

   - Property: Merge operations are associative + commutative
   - Test: `merge(merge(a, b), c) = merge(a, merge(b, c))`
   - Bug caught: Different merge orders giving different results

4. **Numeric Precision**
   - Property: Addition is "approximately" associative for floats
   - Test: `|(a + b) + c - (a + (b + c))| < epsilon`
   - Bug caught: Precision loss from bad operation ordering

**Exercises:**

- [ ] Test a map merge function (right-side wins on conflicts)
- [ ] Find which operations on your codebase should be associative
- [ ] Test floating point operations with epsilon tolerance

### Day 4: Advanced Patterns

**Abstract Algebra Structures:**

```
Magma: Just a binary operation (a ∘ b)
  ↓ + Associativity
Semigroup: Associative operation
  ↓ + Identity
Monoid: Semigroup + identity element
  ↓ + Inverse
Group: Monoid + inverse for every element
  ↓ + Commutativity
Abelian Group: Commutative group
```

**When to use each:**

- **Semigroup**: String concatenation, list concatenation
- **Monoid**: Addition, multiplication, set union, map merge
- **Group**: Integer addition (has negatives), vector addition
- **Abelian Group**: Integer addition, XOR operation

**Power moves:**

1. **Homomorphism**: Function that preserves structure

   - `f(a ∘ b) = f(a) ⊕ f(b)`
   - Example: `log(a * b) = log(a) + log(b)`
   - Use case: Testing that serialization preserves operations

2. **Isomorphism**: Two structures are "the same" mathematically

   - Example: Complex numbers ≅ 2D vectors
   - Use case: Testing different implementations are equivalent

3. **Property composition**: Chain property tests
   - Test monoid laws, then add inverse test = group test
   - Reuse property testers for more complex structures

**Exercises:**

- [ ] Build a `TestMonoid()` helper that runs associativity + identity
- [ ] Build a `TestGroup()` helper that includes inverse tests
- [ ] Create custom generators for your domain types

## Resources to Learn More

- **Book**: "Abstract Algebra" by Pinter (very readable intro)
- **Video**: 3Blue1Brown's "Essence of Linear Algebra" (groups in geometry)
- **Paper**: "QuickCheck" by Koen Claessen (original property-based testing)
- **Practice**: Implement vector spaces, rings, fields using goprops

## Common Pitfalls

1. **Floating Point**: Not truly associative due to precision

   - Use approximate equality with epsilon
   - Test on known problematic cases (very large + very small numbers)

2. **Non-comparable Types**: Maps and slices can't be used with `comparable`

   - Use custom equality functions
   - Or wrap in comparable structs

3. **Side Effects**: Properties assume pure functions

   - Cache tests need fresh instances
   - Database operations might need transactions/rollbacks

4. **Randomness**: Sometimes you need specific edge cases
   - Combine property tests with example-based tests
   - Use custom generators for domain-specific values

## Success Metrics

After 3-4 days, you should be able to:

- [ ] Identify which operations in your code should be associative
- [ ] Write property tests that catch structural bugs (like plyGO's valueKey)
- [ ] Explain why matrix multiplication is associative but not commutative
- [ ] Build custom generators for your domain types
- [ ] Understand when to use monoids vs groups
- [ ] Catch bugs that unit tests miss (missing cases, wrong assumptions)

## What Makes This Different from Unit Tests

**Unit test:**

```go
assert.Equal(t, 5, add(2, 3))  // Tests ONE case
```

**Property test:**

```go
goprops.Associative(t, add, gen)  // Tests 100+ random cases
                                   // Checks MATHEMATICAL LAW
```

**The insight:** Laws are more general than examples. One property test replaces hundreds of unit tests.

---

**Next Steps:**

1. Run the examples: `go test -v examples/*.go`
2. Pick one operation from your codebase
3. Write property tests for it
4. Watch it catch bugs you didn't know existed

Welcome to algebraic thinking! 🎓
