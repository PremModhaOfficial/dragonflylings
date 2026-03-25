# Explain It — lists4

## The Challenge

You're designing a "recent activity" feature: show each user their last 20 actions. You have 1 million users, and each performs ~100 actions per day.

Write your design in `feynman/explanations/lists4.md`.

## Your Explanation Should Cover

1. **Why you need LTRIM** — without it, each user's list grows by 100 items per day. After a year: 36,500 items per user × 1M users = 36 billion stored actions. With LTRIM after each push: exactly 20 items per user, always.

2. **Is LPUSH + LTRIM atomic?** — two separate commands. What could go wrong if another process reads between the push and the trim? Is there a way to make it atomic? (Hint: pipeline or Lua script.)

3. **LTRIM vs explicit delete** — you could also delete old items with `DEL` + rebuild. Why is LTRIM better? When might rebuild be necessary?

## Quality Check

Include the memory math: 20 items × ~50 bytes each × 1M users = ? MB. Compare to the unbounded version after 1 year.
