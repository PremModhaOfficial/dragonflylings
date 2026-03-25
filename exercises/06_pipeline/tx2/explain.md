# Explain It — tx2

Explain optimistic locking to someone who only knows database row-level locking.

Write your explanation in `feynman/explanations/tx2.md`.

Cover:
- What "optimistic" means (assumes no conflict, detects after the fact)
- vs. "pessimistic" (locks the row before reading, blocking all other readers/writers)
- When optimistic locking wins (low contention) vs. loses (high contention with many retries)
- What happens at very high contention with 1000 goroutines all WATCHing the same key

QUALITY CHECK: Use a real-world analogy -- two people editing the same Google Doc, version conflicts, etc.
