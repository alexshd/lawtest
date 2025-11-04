# Learning Resources for Group Theory & Property-Based Testing

## You Said: "I Learn Fast" - Here's Your Path

You didn't know Kaluza-Klein existed 3 weeks ago. Now you do. Same energy for this.

---

## Week 1: Abstract Algebra Foundations

### **Book: "A Book of Abstract Algebra" by Charles C. Pinter**

- **Why**: Readable, practical, minimal prerequisites
- **What**: Groups, rings, fields from first principles
- **Time**: 2-3 hours per chapter, chapters 1-9
- **Get it**: Dover edition, $15, or PDF available
- **Focus**: Chapters 3-5 (Groups), Chapter 9 (Homomorphisms)

**How to read it:**

1. Do EVERY exercise (this is where learning happens)
2. When you see a theorem, try to prove it before reading the proof
3. Look for "where would this catch a bug in code?"

### **Video Series: "Abstract Algebra" by Socratica (YouTube)**

- **Link**: https://www.youtube.com/playlist?list=PLi01XoE8jYoi3SgnnGorR_XOW3IcK-TP6
- **Why**: Visual, 10-minute chunks, great for intuition
- **Watch**: Videos 1-15 (Groups fundamentals)
- **Speed**: 1.5x, pause to think, rewind when needed

### **Interactive: "Visual Group Theory" by Nathan Carter**

- **Book + Software**: Visualizes group operations geometrically
- **Why**: Makes abstract concepts concrete
- **Focus**: Symmetry groups, Cayley diagrams
- **Connection**: See how groups appear in graphics, games, physics

---

## Week 2: Property-Based Testing Origins

### **Paper: "QuickCheck: A Lightweight Tool for Random Testing"**

- **Authors**: Koen Claessen & John Hughes (2000)
- **Link**: https://www.cs.tufts.edu/~nr/cs257/archive/john-hughes/quick.pdf
- **Why**: The original property-based testing paper
- **Length**: 13 pages, readable
- **Key insight**: "Specification = Properties, not examples"

### **Blog: "Choosing properties for property-based testing"**

- **Author**: Scott Wlaschin (F# for Fun and Profit)
- **Link**: https://fsharpforfunandprofit.com/posts/property-based-testing-2/
- **Why**: Practical guide to finding properties in real code
- **Pattern library**: 7 types of properties with examples

### **Video: "Property-Based Testing with PropEr, Erlang, and Elixir"**

- **Speaker**: Fred Hebert
- **Book**: "Property-Based Testing with PropEr, Erlang and Elixir"
- **Why**: Shows property testing catching REAL distributed systems bugs
- **Example**: Found race conditions in Riak (production database)

---

## Week 3: Category Theory (The Next Level)

### **Book: "Category Theory for Programmers" by Bartosz Milewski**

- **Format**: Free online, blog series + PDF
- **Link**: https://bartoszmilewski.com/2014/10/28/category-theory-for-programmers-the-preface/
- **Why**: Explains functors, monads, morphisms with code examples
- **Language**: Haskell + C++ (but concepts transfer to Go)
- **Read**: Chapters 1-10 first, then skip around

**Key chapters for you:**

- Chapter 3: "Categories Great and Small" (monoids appear!)
- Chapter 5: "Products and Coproducts"
- Chapter 10: "Natural Transformations"

### **Video: "Category Theory" by Bartosz Milewski (YouTube)**

- **Link**: https://www.youtube.com/playlist?list=PLbgaMIhjbmEnaH_LTkxLI7FMa2HsnawM_
- **Format**: 30 lectures, ~30 mins each
- **Watch**: First 10 lectures
- **Speed**: 1.25x, pause for diagrams

### **Interactive: "Awodey's Category Theory" (if you want rigor)**

- **Book**: "Category Theory" by Steve Awodey
- **Why**: Standard textbook, more mathematical
- **Use case**: Reference when you need formal definitions

---

## Week 4: Algebraic Structures in Code

### **Paper: "Algebraic Effects for the Rest of Us"**

- **Why**: Shows how algebra structures real systems
- **Connection**: Effect handlers = free monads
- **Example**: Error handling, logging, state management

### **Blog Series: "Algebra Driven Design"**

- **Author**: Sandy Maguire
- **Book**: "Algebra-Driven Design"
- **Why**: Uses algebra to design better APIs
- **Example**: Designing a music synthesizer with algebraic laws

### **Repository: "Laws" in Haskell/Cats/etc.**

- **Haskell Prelude laws**: https://hackage.haskell.org/package/base/docs/Prelude.html
- **Cats laws**: https://github.com/typelevel/cats/tree/main/laws
- **Why**: See REAL property tests in production libraries

---

## Bonus: Topics You'll Connect To

### **1. Quantum Computing (You mentioned this)**

- **Connection**: Unitary groups, SU(2), rotation groups
- **Resource**: "Quantum Computing Since Democritus" by Scott Aaronson
- **Group theory appears**: Quantum gates form groups (composable, reversible)

### **2. Kaluza-Klein Theory (You mentioned this)**

- **Connection**: Gauge groups, Lie groups, symmetry
- **Resource**: "Gauge Fields, Knots and Gravity" by Baez & Muniain
- **Group theory appears**: Compact Lie groups in extra dimensions

### **3. Game Development / Graphics**

- **Connection**: Rotation groups SO(3), quaternions
- **Resource**: "Geometric Algebra for Computer Science" by Dorst et al.
- **Use case**: Your light propagation pathfinding insight

### **4. Distributed Systems (Your work with Burgess)**

- **Connection**: CRDTs = commutative + idempotent operations
- **Paper**: "A comprehensive study of CRDTs" by Shapiro et al.
- **Group theory appears**: Merge operations form semilattices

---

## Tools to Learn By Doing

### **1. Proof Assistants**

- **Lean 4**: https://lean-lang.org/ (proof assistant with theorem proving)
- **Why**: Actually PROVE properties, not just test them
- **Example**: Prove associativity of your operations formally

### **2. Computer Algebra Systems**

- **SageMath**: https://www.sagemath.org/ (Python-based)
- **GAP**: https://www.gap-system.org/ (Group theory specific)
- **Why**: Experiment with groups, test conjectures

### **3. Visualization Tools**

- **Manim** (3Blue1Brown's library): Animate group operations
- **Graphviz**: Draw Cayley diagrams
- **Why**: See structure, not just symbols

---

## Communities to Join

### **Math Stack Exchange** - "abstract-algebra" tag

- Ask questions, read answers
- See how others think about problems

### **r/math and r/abstractalgebra** (Reddit)

- Weekly "What are you reading" threads
- Lots of self-learners like you

### **Category Theory Zulip** - https://categorytheory.zulipchat.com/

- Active community, beginners welcome
- Fast responses to questions

### **Papers We Love** - https://paperswelove.org/

- Reading group for CS papers
- Property testing / formal methods papers

---

## The 30-Day Challenge (Since You Learn Fast)

### **Days 1-7: Basics**

- [ ] Read Pinter chapters 3-5
- [ ] Watch Socratica videos 1-15
- [ ] Do exercises in `LEARNING_GUIDE.md`
- [ ] Build custom generator for complex domain type

### **Days 8-14: Property Testing**

- [ ] Read QuickCheck paper
- [ ] Read Scott Wlaschin's property series
- [ ] Apply goprops to 3 of your existing projects
- [ ] Find 1 real bug

### **Days 15-21: Category Theory**

- [ ] Read Milewski chapters 1-5
- [ ] Watch Milewski videos 1-5
- [ ] Understand: Functor, Applicative, Monad as patterns
- [ ] Connect to Go interfaces + generics

### **Days 22-30: Advanced**

- [ ] Read Fred Hebert on distributed systems properties
- [ ] Study CRDTs and commutative operations
- [ ] Build property tests for Promise Theory concepts
- [ ] Write blog post explaining what you learned

---

## How to Know You've "Got It"

### **Week 1 Success Metrics:**

- [ ] Can explain associativity to a 10-year-old
- [ ] Can identify which operations have inverses
- [ ] Can spot non-commutative operations in code

### **Week 2 Success Metrics:**

- [ ] Can write property tests for any function
- [ ] Can distinguish "example" from "property"
- [ ] Can generate counterexamples when properties fail

### **Week 3 Success Metrics:**

- [ ] Can explain what a functor is (in code, not theory)
- [ ] Can identify monoids in your codebase
- [ ] Can use composition to simplify complex operations

### **Week 4 Success Metrics:**

- [ ] Can design APIs using algebraic laws
- [ ] Can prove properties hold (not just test)
- [ ] Can teach this to others

---

## Books by Difficulty (Ranked)

### **Beginner (Start Here):**

1. Pinter - "A Book of Abstract Algebra" ⭐⭐⭐⭐⭐
2. Milewski - "Category Theory for Programmers" (first 10 chapters)
3. Wlaschin - "Domain Modeling Made Functional"

### **Intermediate (After 2-3 weeks):**

1. Mac Lane - "Categories for the Working Mathematician" (reference)
2. Awodey - "Category Theory" (formal)
3. Baez & Stay - "Physics, Topology, Logic and Computation" (connections)

### **Advanced (When ready):**

1. Aluffi - "Algebra: Chapter 0" (comprehensive)
2. Jacobson - "Basic Algebra I & II" (encyclopedic)
3. Leinster - "Basic Category Theory" (concise)

---

## Papers You'll Want to Read

### **Property Testing:**

- QuickCheck (2000) - Original paper
- "Testing Monadic Code with QuickCheck" (2002) - Claessen
- "SmallCheck and Lazy SmallCheck" (2008) - Runciman
- "CRDTs: Consistency without concurrency control" (2009) - Shapiro

### **Algebra in CS:**

- "Data types à la carte" (2008) - Swierstra
- "Functional Pearl: Implicit Configurations" (2000) - Kiselyov
- "Free Monads and Free Monoids" (2008) - Various

### **Formal Methods (That actually work):**

- "How Amazon Web Services Uses Formal Methods" (2015)
- "Jepsen: On the perils of network partitions" (Aphyr's blog series)

---

## YouTube Channels Beyond Socratica

1. **3Blue1Brown** - "Essence of Linear Algebra" (groups via geometry)
2. **Mathologer** - Abstract algebra visualized
3. **Mu Prime Math** - Group theory basics
4. **Richard Southwell** - Category theory for scientists
5. **The Math Sorcerer** - Study tips, book reviews

---

## The "Aha!" Moments to Look For

1. **Moment 1**: "Wait, cache merge is just a monoid!"
2. **Moment 2**: "Error handling is just the Maybe monad"
3. **Moment 3**: "Promises are applicative functors"
4. **Moment 4**: "CRDTs are commutative monoids"
5. **Moment 5**: "My pathfinding IS Dijkstra because the structure is identical"

---

## Final Advice (From Someone Who Knows You Learn Fast)

**Don't read linearly.** Jump between:

- Theory (Pinter)
- Practice (goprops)
- Connections (Kaluza-Klein, quantum, distributed systems)
- Code (build, test, ship)

**Do exercises.** Theory without practice is useless. Practice without theory is blind.

**Teach it.** After week 2, explain property testing to someone. After week 4, write a blog post.

**Connect everything.** You already do this (light propagation → pathfinding). Now: groups → CRDTs → distributed consensus → Promise Theory.

**Ship something.** Don't just learn—build. goprops is your first. What's next? A CRDT library? A proof assistant? A Promise Theory debugger?

You've got the pattern recognition. Now add the vocabulary. The math gives you **precision** where you already have **intuition**.

Go learn. Fast. 🚀

---

**Start tonight:**

1. Watch Socratica videos 1-5 (50 minutes)
2. Do first 5 exercises from `LEARNING_GUIDE.md`
3. Order Pinter's book ($15, arrives in 2 days)
4. Read QuickCheck paper before bed (13 pages)

**Tomorrow morning:**

- Apply goprops to one of your existing projects
- Find the first bug it catches
- Tell me about it

You're ready.
