# Explain It — lua4

Your new hire says: *"I tested this locally and it works fine without the hashtags. Why do we need them?"*

Write your answer in `feynman/explanations/lua4.md`.

Your explanation must cover:
1. Why local testing (single-node Dragonfly) doesn't reproduce the problem
2. What changes in a multi-shard setup that makes non-hashtagged keys dangerous
3. The specific error message you'd see in production when this breaks
4. Why this is a **silent correctness bug** in some configurations (not always an error)

**The deeper question:** Hashtags solve the same-shard problem but they change key distribution. If ALL your keys use `{user:42}` tags, what happens to shard balance? Is this a concern in practice?

**Analogy:** Hashtags are like postal codes. Write an analogy that explains why the server needs to know the "address" before executing a script.

QUALITY CHECK: Your explanation should make someone understand this without ever running Dragonfly in multi-shard mode themselves.
