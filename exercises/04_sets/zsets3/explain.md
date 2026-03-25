# Explain It — zsets3

## The Challenge

A player in your game completes a round and earns 500 points. Your code calls `ZADD leaderboard 500 "alice"`. Alice's score on the leaderboard is now 500, no matter what she had before.

Your support inbox fills with "my score was reset!" complaints.

Write your post-mortem in `feynman/explanations/zsets3.md`.

## Your Explanation Should Cover

1. **ZADD vs ZINCRBY semantics** — ZADD sets the score to an absolute value. ZINCRBY adds to the existing score. One word difference in code; completely different behavior.

2. **When would you use ZADD to update?** — if you're syncing a score from an authoritative source (not incrementing), ZADD is correct. If you're adding "points earned this round," ZINCRBY is correct.

3. **The ZRANK vs ZSCORE distinction** — explain to a non-technical PM: "The player's score is 1500. Their rank is 3rd." Why are these two different numbers? Why do you need both?

## Quality Check

Your post-mortem should include a "lessons learned" section with the specific API confusion that caused the bug, and a PR description for the fix.
