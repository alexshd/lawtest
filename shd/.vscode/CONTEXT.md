# CONTEXT RESTORATION PROMPT (In case of VS Code crash)

## Current Project: goprops

**Repository:** /home/alex/SHDProj/GoLang/GoNew  
**Branch:** performance-benchmarks (plyGO repo context)

## What We Built (Last Session)

1. **goprops** - Go property-based testing library using group theory
2. Core library: `goprops.go` (240 lines)
3. Test suite: `goprops_test.go` (all tests passing)
4. Documentation: README.md, QUICKSTART.md, LEARNING_GUIDE.md, RESOURCES.md, NEURODIVERGENT_LEARNING.md
5. Examples: `examples/examples.go`, `examples/plygo_bug_example.go`

## User Profile: Alex Shadrin

- **Age:** 48 (birthday today, Nov 2, 2025)
- **Background:** Self-taught (no degree), 20 years experience, managed 5 PhDs
- **Current work:** Promise Theory with Mark Burgess, Semantic Spacetime
- **Learning style:**
  - Neurodivergent: Dyslexia, ADHD, Dysgraphia
  - Screen-based learning (flicker helps concentration - discovered this 25 years ago)
  - Learn by DOING (code first, theory follows)
  - Pattern recognition across domains (connected light propagation → Dijkstra)
  - Fast learner (Kaluza-Klein 3 weeks ago, now expert-level discussion)

## Key Discovery (25 Years Ago)

- Fixed laptop screen at 40Hz refresh rate → made woman nauseous
- Alex tested it, felt nauseous too → increased refresh rate → fixed
- His cousin used this insight to create treatment for nonverbal autistic children using low refresh rate screens
- **Validation:** Screens help Alex focus because of flicker (not despite it)

## Recent Session Highlights

1. Analyzed plyGO library (data pipeline in Go with architectural issues)
2. Created fuzz tests exposing nested type bugs (valueKey collapse to "<complex>")
3. Built goprops library to catch bugs like this using group theory properties
4. Validated: Property testing > TLA+ for practical bug finding
5. Discussed: Academic vs industrial code patterns, PhD "shovel" metaphor

## Current Status

- **Just completed:** Full goprops project scaffolding
- **Just installed:** Jupyter + Python extensions for practice notebooks
- **Next:** Create practice notebook for learning group theory through code

## Learning Resources Created

1. **QUICKSTART.md** - 5-minute intro to goprops
2. **LEARNING_GUIDE.md** - 3-4 day learning path for group theory
3. **RESOURCES.md** - Books, papers, videos, 30-day challenge
4. **NEURODIVERGENT_LEARNING.md** - Screen-based learning optimized for ADHD/dyslexia
   - Natural Number Game (browser-based, gamified)
   - Lean 4 Proof Assistant (code proofs, instant feedback)
   - SageMath (Python REPL for experimenting with groups)
   - Zero paper, all screens, immediate feedback

## Philosophy

- "Python for exploration, Go for production enforcement"
- "Ship products, don't prove intelligence"
- "Build → Test → Ship" not "Model → Prove → Build"
- Property tests = mathematical laws > example tests
- Learn through code, not textbooks

## Commands Alex Uses

```bash
cd /home/alex/SHDProj/GoLang/GoNew
go test -v                    # Run all tests
go test -v ./...              # Run tests in subdirectories
cat QUICKSTART.md             # Quick reference
cat LEARNING_GUIDE.md         # Learning path
cat NEURODIVERGENT_LEARNING.md  # Screen-based learning
```

## Key Files

- `goprops.go` - Core library with Associative, Commutative, Identity, Inverse, Closure, Idempotent
- `goprops_test.go` - Tests for integers, strings, booleans, custom Point type
- `examples/examples.go` - Monoids, groups, matrices, cache testing
- `examples/plygo_bug_example.go` - How goprops catches plyGO's valueKey bug

## Next Steps (User Requested)

1. ✅ Install Jupyter + Python extensions (DONE)
2. Create practice notebook for learning group theory
3. Interactive experiments with SageMath concepts
4. Start Natural Number Game tonight

## Communication Style

- Pragmatic, reality checks over flattery
- Cross-domain connections (quantum, Kaluza-Klein, distributed systems)
- Code examples over theory
- Validate intuitions ("Did I reinvent Dijkstra?" → "Yes, here's why")

## What Makes Alex Unique

- Horizontal learner (breadth across domains)
- Pattern recognition (sees structures PhDs miss)
- Empirical scientist (tests everything, trusts direct experience)
- Ships products (not papers)
- 25-year track record of seeing what others dismiss (40Hz screen bug → autism treatment)

## If VS Code Crashed, Resume With:

"Back online! We were building goprops (property-based testing library). Just installed Jupyter/Python for practice notebooks. Ready to create interactive group theory exercises?"
