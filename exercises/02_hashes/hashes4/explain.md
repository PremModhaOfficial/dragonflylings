# Explain It — hashes4

## The Challenge

You're designing a user profile store for 1 million users, each with 8 fields. Compare these two approaches and recommend one in `feynman/explanations/hashes4.md`.

**Option A:** `user:1:name`, `user:1:email`, `user:1:age`, ... (8M string keys)
**Option B:** `user:1` hash with 8 fields (1M hash keys)

## Your Explanation Should Cover

1. **Memory overhead math** — each Redis key has ~50-60 bytes of overhead (pointer, expiry, LRU clock, type). For 8M keys vs 1M keys, what's the overhead difference? Multiply by your user count.

2. **Operational complexity** — "Delete user 42 and all their data." Write the code for Option A and Option B. Which is simpler? Which is safer?

3. **When Option A might be better** — is there a scenario where separate string keys are the right choice over a hash? (Hint: think about individual field TTLs, or cross-user MGET patterns.)

## Quality Check

Your recommendation should include a concrete number: "Option B uses approximately X MB less memory for 1M users."
