# Explain It — expire1

A junior developer asks: "Why do we need PERSIST? If I don't want a key to expire, I just won't call EXPIRE on it."

Write your answer in `feynman/explanations/expire1.md`.

Your explanation should cover:
- When would a key get an expiry set that you then want to remove?
- What real-world scenario requires PERSIST? (Think: a temporary session that the user keeps extending)
- What is the TTL of a key that has never had EXPIRE called on it?
- Is there a difference between TTL=-1 and TTL=-2?

QUALITY CHECK: Could a 12-year-old understand your explanation?
If you used jargon, replace it with an analogy. Think of a parking meter.
