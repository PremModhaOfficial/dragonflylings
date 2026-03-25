# Explain It — prod3

Your team is designing a rate limiter for your API. Two engineers disagree:

Engineer A: *"Just use INCR with a minute key. Simple, fast, done."*
Engineer B: *"We need sliding window — fixed window lets users burst double the limit at minute boundaries."*

Write your answer in `feynman/explanations/prod3.md`.

Your explanation must cover:
1. A concrete example showing the fixed-window boundary burst (with numbers)
2. How the sliding window eliminates this with a diagram or timeline
3. The cost: sorted sets use more memory than a single counter. How much more?
4. When is fixed window "good enough" and when do you need sliding window?

**The sorted set trade-off:** Each request adds a member to the sorted set. For 1 million requests/minute across 10,000 users, how many sorted set entries exist? Is this feasible?

**Token bucket vs sliding window:** There's a third algorithm — token bucket. How does it compare to sliding window? What are the trade-offs?

QUALITY CHECK: Can you draw the fixed-window boundary burst attack on a whiteboard in 2 minutes?
