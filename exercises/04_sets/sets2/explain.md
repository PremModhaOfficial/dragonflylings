# Explain It — sets2

## The Challenge

You're implementing "People You Might Know" for a social network. The algorithm is: people that your friends follow, that you don't already follow.

Write the Redis implementation in `feynman/explanations/sets2.md`.

## Your Explanation Should Cover

1. **Map the algorithm to set operations** — "friends of friends" = SUNION of all your friends' follow sets. "Exclude people you already follow" = SDIFF. Show the full pipeline.

2. **SINTERSTORE and SUNIONSTORE** — instead of returning results to the app, these commands store the result in a new key. When would you use store variants? (Hint: caching intermediate results)

3. **Time complexity** — SINTER with two sets of sizes M and N is O(M×N) in the worst case. What's the practical concern for large social graphs? What's the mitigation?

## Quality Check

Your explanation should include a worked example with real users. Draw the sets and show which operation produces each result.
