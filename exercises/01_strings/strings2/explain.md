# Explain It — strings2

## The Challenge

Your teammate argues: "Why use Redis TTL for sessions? We can just store the creation time and check it in application code."

Write your rebuttal in `feynman/explanations/strings2.md`.

## Your Explanation Should Cover

1. **What TTL does automatically that your code can't** — expired keys are cleaned up by Redis without any application action. If the application crashes, sessions still expire. With application-managed expiry, you need a cleanup job that might not run.

2. **The "TTL=-1" production footgun** — what happens to your Redis memory if every session is stored with no TTL? Show the math: 1 million users × 1KB session data = ?

3. **When you'd use TTL vs. explicit cleanup** — TTL is great for: sessions, caches, rate limit windows, OTP codes. It's not great for: data you need to query historically, data with complex expiry logic.

## Quality Check

Your explanation should make someone understand why Redis key expiry is a feature, not a limitation. Include one real-world example where forgetting TTL caused a production incident (you can invent a realistic one).
