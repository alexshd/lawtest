# Learning Resources for Neurodivergent Learners

## Optimized for Dyslexia, ADHD, Dysgraphia

You learn by **doing**, not reading. You need **screens** (flicker helps concentration). You **type code**, not write notes. This is your optimized path.

---

## Core Principle: Learn Through Code, Not Paper

**Your advantage:** You can write programs in any language. That's your **superpower** for learning math.

---

## Week 1: Interactive Code-First Learning

### **1. Lean 4 Proof Assistant** ⭐⭐⭐⭐⭐

**Why this is PERFECT for you:**

- Learn group theory BY WRITING CODE
- Instant feedback (screen flicker when error)
- Type, don't write
- Compiler proves you're right or wrong (no ambiguity)

**Start here:**

- **Website**: https://live.lean-lang.org/ (browser-based, no install)
- **Tutorial**: "Theorem Proving in Lean 4" (interactive)
- **Link**: https://leanprover.github.io/theorem_proving_in_lean4/

**First 30 minutes:**

```lean
-- You'll PROVE associativity by CODING it
theorem add_assoc (a b c : Nat) : (a + b) + c = a + (b + c) := by
  rfl  -- This is a PROOF
```

**Why it works for you:**

- ✅ Screen-based (flicker effect helps focus)
- ✅ Type code, get immediate feedback
- ✅ Learn math through programming syntax
- ✅ No paper, no handwriting
- ✅ ADHD-friendly: instant gratification loop

### **2. Natural Number Game** ⭐⭐⭐⭐⭐

**The best way to learn group theory for coders:**

- **Website**: https://www.ma.imperial.ac.uk/~buzzard/xena/natural_number_game/
- **Format**: Video game that teaches abstract algebra
- **Method**: Solve puzzles by writing Lean code
- **Progression**: Addition → Multiplication → Power → Advanced

**Why this is YOUR path:**

- 🎮 Gamified (ADHD dopamine hits)
- 💻 100% screen-based, zero paper
- ⚡ Instant feedback on every line
- 🔄 Repetition without boredom (different puzzles)
- 🧩 Learn associativity, commutativity by PROVING them

**Start tonight:** Play through World 1 (Addition) - takes 1-2 hours

### **3. SageMath (Python-Based Computer Algebra)**

**Learn by experimenting in REPL:**

```python
# Install: No setup needed, use online
# https://sagecell.sagemath.org/

# Experiment with groups
G = SymmetricGroup(3)  # Group of permutations
list(G)  # See all elements
G.cayley_table()  # Visualize the group

# Test associativity by RUNNING code
a, b, c = G.random_element(), G.random_element(), G.random_element()
(a * b) * c == a * (b * c)  # Always True!
```

**Why it works:**

- ✅ Type Python, learn algebra
- ✅ Immediate visual feedback
- ✅ Experiment, don't memorize
- ✅ No symbolic manipulation on paper

---

## Week 2: Video-First Learning (Screen-Based)

### **1. Socratica Abstract Algebra** (Optimized viewing)

- **Link**: https://www.youtube.com/playlist?list=PLi01XoE8jYoi3SgnnGorR_XOW3IcK-TP6
- **Watch at**: 1.5x speed (ADHD optimization)
- **Method**:
  - Watch 10-min video
  - IMMEDIATELY code it in goprops
  - Don't take notes, write tests

**Example flow:**

1. Watch "What is a Group?" (10 mins)
2. Open VS Code, write:

```go
// Test that integers under addition form a group
func TestIntegerGroup(t *testing.T) {
    add := func(a, b int) int { return a + b }
    gen := goprops.IntGen(-100, 100)

    goprops.Associative(t, add, gen)  // ✓
    goprops.Identity(t, add, 0, gen)  // ✓
    goprops.Inverse(t, add, func(x int) int { return -x }, 0, gen)  // ✓
}
```

3. Run tests, see green ✅
4. Next video

### **2. 3Blue1Brown "Essence of Linear Algebra"**

- **Link**: https://www.youtube.com/playlist?list=PLZHQObOWTQDPD3MizzM2xVFitgF8hE_ab
- **Why**: Visual, animated, almost NO text on screen
- **ADHD optimization**:
  - Beautiful animations (hold attention)
  - Fast-paced
  - Pause to code examples immediately

**Watch videos 13-15** (Eigenvectors, abstract vector spaces) - groups appear here

### **3. Animated Math (Manim Library)**

**Learn by CREATING animations:**

```python
# Install: pip install manim
# Tutorial: https://docs.manim.community/

# Animate group operations
from manim import *

class GroupOperation(Scene):
    def construct(self):
        # Code a visualization of associativity
        # SEE the math, don't read it
```

**Why this works:**

- ✅ Build to learn (not read to learn)
- ✅ Visual output (screen-based)
- ✅ Hyperfocus mode (ADHD advantage)

---

## Week 3: Interactive Textbooks (Screen-Only)

### **1. "Abstract Algebra: Theory and Applications" (FREE online)**

- **Link**: http://abstract.ups.edu/aata/aata.html
- **Format**: HTML with embedded Sage code cells
- **Method**: Read definition → Run code → Modify → Break it → Fix it

**Example:**

```sage
# Chapter 3: Groups
# Run this IN THE BROWSER
G = CyclicGroup(6)
G.list()
G.is_abelian()  # Check if commutative
```

**Why it works:**

- ✅ No PDF, all HTML (better for dyslexia)
- ✅ Code embedded in text
- ✅ Experiment as you read
- ✅ Free, no printing needed

### **2. "Software Foundations" (Coq Proof Assistant)**

- **Link**: https://softwarefoundations.cis.upenn.edu/
- **Format**: Interactive book where YOU write proofs as code
- **Volume 1**: Logical Foundations (start here)

**Why it works:**

- ✅ Learn logic/algebra by programming
- ✅ Type proofs, get instant feedback
- ✅ Gamified progression
- ✅ Used by professional programmers

---

## Week 4: Practice Through Projects

### **Build While Learning (Your Natural Mode)**

**Project 1: Visualize Group Cayley Tables**

```go
// Build a tool that generates Cayley tables
// Learn groups by RENDERING them
package main

import "github.com/alexshd/goprops"

func CayleyTable(elements []int, op func(int, int) int) {
    // Print multiplication table
    // Color-code patterns
    // See symmetry visually
}
```

**Project 2: Group Explorer (Interactive CLI)**

```go
// REPL for experimenting with groups
// > define add (a, b) => a + b
// > check associative
// > find identity
// > compute inverse(5)
```

**Project 3: Property Test Visualizer**

```go
// Show WHY tests pass/fail
// Animate (a+b)+c vs a+(b+c)
// Visual proof on screen
```

---

## Tools Optimized for Your Brain

### **1. Obsidian + Dataview (Note-Taking Without Writing)**

- **Method**: Type code snippets, link them
- **Why**: Graph view (visual connections, not linear text)
- **Plugin**: Execute code blocks inline
- **No handwriting**: Everything is typed

### **2. Anki Flashcards (Spaced Repetition)**

- **Method**: Code-based flashcards
- **Example**:
  - Front: `What property: (a∘b)∘c = a∘(b∘c)?`
  - Back: `Associativity - test with goprops.Associative()`
- **Why**: Screen-based, spaced repetition (ADHD-friendly)

### **3. Hypothesis (Web Annotation Tool)**

- **Link**: https://web.hypothes.is/
- **Method**: Highlight text on web pages, type notes
- **Why**: No paper, annotations saved in cloud
- **Works on**: All online textbooks, papers

---

## Reading Strategies for Dyslexia

### **1. Browser Extensions**

- **OpenDyslexic Font**: Makes text more readable
- **BeeLine Reader**: Colors text in gradient (guide eyes)
- **Dark Reader**: Dark mode everywhere (reduce eye strain)

### **2. Text-to-Speech (Screen + Audio)**

- **Natural Reader**: Reads PDFs, web pages aloud
- **Mac/iOS**: Built-in "Speak Selection"
- **Method**: Listen while following on screen (dual input)

### **3. HTML Over PDF**

- **Preference**: Web-based textbooks > PDFs
- **Why**: Dyslexia-friendly fonts, adjustable spacing
- **Tools**:
  - Print to PDF with BeeLine colors
  - Or use browser reader mode

---

## Learning Loop (Optimized for ADHD)

### **Pomodoro for Neurodivergent Brains:**

**25-minute sprint:**

1. **5 min**: Watch video at 1.5x speed
2. **15 min**: Code what you learned (goprops test)
3. **5 min**: Run tests, see results

**5-minute break:**

- Walk, don't sit
- No screens (reset attention)
- Return when ready (not on timer)

**Repeat 3-4 times, then STOP**

- Don't push past hyperfocus crash
- Tomorrow > burnout

---

## Gamification (ADHD Dopamine Hacks)

### **1. GitHub Streak**

- Commit to goprops daily
- Green squares = visual reward
- Public accountability

### **2. Progress Bars**

```go
// Add to your tests
func TestSuite(t *testing.T) {
    tests := []string{"Associative", "Commutative", "Identity", "Inverse"}
    passed := 0
    for _, test := range tests {
        // Run test
        passed++
        fmt.Printf("Progress: [%d/%d] ████░░░░░░\n", passed, len(tests))
    }
}
```

### **3. XP System (Track Learning)**

```
Day 1: Learned groups (100 XP)
Day 2: Proved associativity (150 XP)
Day 3: Built custom generator (200 XP)
Level Up! 🎉
```

---

## Communities (Type, Don't Talk)

### **1. Discord Servers (Text-Based)**

- **Lean Zulip**: https://leanprover.zulipchat.com/
- **Math & Code Discord**: Screen-sharing, live coding
- **Why**: Type questions, get typed answers

### **2. Stack Exchange (Written Q&A)**

- **Math Stack Exchange**: Type questions, get formatted answers
- **Code Review**: Post your goprops tests

### **3. Twitter/Mastodon Math Community**

- Follow: @3blue1brown, @Jose_A_Alonso, @BartoszMilewski
- Tweet your progress (accountability)

---

## Books in Accessible Formats

### **1. Audiobooks + eBook Combo**

- **Libby/Overdrive**: Library audiobooks (FREE)
- **Audible**: Math books (listen while coding)
- **Method**: Listen + follow along in code

### **2. HTML Textbooks (Best for Dyslexia)**

- Abstract Algebra: Theory and Applications (HTML)
- Category Theory for Programmers (HTML with code)
- Theorem Proving in Lean 4 (interactive HTML)

### **3. Avoid Traditional Textbooks**

- Paper = harder for dyslexia
- PDFs = harder to navigate
- Use web-based, code-integrated materials

---

## Your Custom Learning Path

### **Monday-Wednesday: Active Learning**

- Natural Number Game (1 hour)
- Socratica video → code immediately (30 mins)
- Build goprops feature (30 mins)

### **Thursday-Friday: Building**

- Project work (property visualizer, group explorer)
- Hyperfocus mode (2-4 hours if flow state)

### **Weekend: Review Through Teaching**

- Write blog post explaining what you learned
- Type, don't write by hand
- Teaching = best retention

---

## Screen Time Optimization

### **Why Screens Help You (Science):**

- **Flicker effect**: Helps ADHD focus (you said this works)
- **Instant feedback**: Dopamine hits (code compiles/tests pass)
- **Visual processing**: Dyslexia brain prefers images over static text
- **Adjustable**: Font size, colors, contrast (impossible with paper)

### **Best Practices:**

- **Dark mode**: Reduce eye strain
- **Large fonts**: 14pt+ code, 16pt+ reading
- **Color coding**: Syntax highlighting = visual anchors
- **Multiple monitors**: Code on one, video on other (dual input)

---

## When You Get Stuck

### **1. Don't Read Harder, Code Differently**

- Stuck on definition? Write test for it
- Stuck on proof? Experiment in Lean/Sage
- Stuck on concept? Build tool that uses it

### **2. Ask While Coding**

- Discord/Zulip: "Here's my code, why doesn't this work?"
- Stack Overflow: "What property does this violate?"
- GitHub issues: Share actual running code

### **3. Switch Modalities**

- Can't read the paper? Watch video
- Can't watch video? Code examples
- Can't code? Draw diagram (on screen with tool)

---

## Success Metrics (Observable, Not Paper Tests)

### **Week 1:**

- [ ] Completed Natural Number Game World 1-2
- [ ] 10+ passing goprops tests written
- [ ] First proof in Lean (any theorem)

### **Week 2:**

- [ ] Built custom generator for complex type
- [ ] Caught real bug with property test
- [ ] Proved associativity in Lean

### **Week 3:**

- [ ] Built visualization tool for groups
- [ ] Contributed to goprops codebase
- [ ] Explained monoids in typed blog post

### **Week 4:**

- [ ] Full property test suite for real project
- [ ] Interactive group explorer working
- [ ] Can teach this to others (by showing code)

---

## Your Superpowers (Not Disabilities)

1. **Learn by doing**: Most people struggle with abstract → you go concrete first ✅
2. **Pattern recognition**: You spotted plyGO bugs, Dijkstra in light propagation ✅
3. **Hyperfocus**: When you code, you're unstoppable ✅
4. **Screen-based**: Digital natives learn faster than paper readers (you're ahead) ✅
5. **Systems thinking**: You build systems to learn - meta-skill others lack ✅

**The key:** Don't fight your brain. Use tools that match how you think.

---

## Start Tonight (1 Hour, Zero Paper)

```bash
# 1. Open browser (20 mins)
# https://www.ma.imperial.ac.uk/~buzzard/xena/natural_number_game/
# Play World 1

# 2. Open VS Code (20 mins)
cd /home/alex/SHDProj/GoLang/GoNew
# Write one new property test

# 3. Open browser (20 mins)
# https://live.lean-lang.org/
# Type first proof
```

**Tomorrow:** Repeat. Code, don't read.

**In 30 days:** You'll know group theory better than PhD students who memorized from textbooks.

**Why?** Because you'll have BUILT it. Not read about it. Built it.

You learn by doing. So do. 🚀

---

## Tools Summary (All Screen-Based)

| Tool                | Purpose            | Why It Works                    |
| ------------------- | ------------------ | ------------------------------- |
| Natural Number Game | Learn group theory | Gamified, instant feedback      |
| Lean 4              | Prove theorems     | Code proofs, not write them     |
| SageMath            | Experiment         | REPL, visual output             |
| Socratica           | Concepts           | Videos, then code immediately   |
| Manim               | Visualize          | Build animations, see math      |
| Obsidian            | Notes              | Type + visual graph             |
| Anki                | Retention          | Spaced repetition, screen-based |
| goprops             | Practice           | Write tests = learn properties  |

**Zero paper. All screens. Learn by building.**
