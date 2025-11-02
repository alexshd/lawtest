# Let's Learn

## Group Definition

A **group** $(G, *)$ consists of:

1. **Set of elements**: $G$ (non-empty set)

2. **Binary operation**: $* : G \times G \to G$ (takes two elements, returns one)

3. **Closure**: $\forall x, y \in G \implies x * y \in G$  
   _(Result stays in the group)_

4. **Identity element**: $\exists e \in G : \forall x \in G, \; e * x = x * e = x$  
   _(There's a "do nothing" element)_

5. **Inverse elements**: $\forall x \in G, \; \exists x^{-1} \in G : x * x^{-1} = x^{-1} * x = e$  
   _(Every element has an "undo")_

6. **Associativity**: $\forall a, b, c \in G, \; (a * b) * c = a * (b * c)$  
   _(Order of operations doesn't matter)_

---

### Examples

**Is a group:**

- $(\mathbb{Z}, +)$ - Integers under addition ✓
- $(\mathbb{Q} \setminus \{0\}, \times)$ - Non-zero rationals under multiplication ✓

**Not a group:**

- $(\mathbb{N}, +)$ - Natural numbers: no inverses (5 needs -5) ✗
- $(\mathbb{Z}, -)$ - Integers under subtraction: not associative ✗

---

### Alternative Notation (Block Format)

$$
\begin{align}
\text{Closure:} \quad & x, y \in G \implies x \circ y \in G \\
\text{Identity:} \quad & \exists e : x \circ e = e \circ x = x \\
\text{Inverse:} \quad & \forall x, \exists x^{-1} : x \circ x^{-1} = e \\
\text{Associativity:} \quad & (a \circ b) \circ c = a \circ (b \circ c)
\end{align}
$$
