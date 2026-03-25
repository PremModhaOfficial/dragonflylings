# Explain It — zsets1

## The Challenge

Your product manager asks: "Why are we using Redis Sorted Sets for the leaderboard? Can't we just sort the scores in Go after reading them from a hash?"

Write your technical response in `feynman/explanations/zsets1.md`.

## Your Explanation Should Cover

1. **The cost of sort-in-app vs sort-in-Redis** — if you store scores in a hash and sort in Go, you must read ALL scores every time a user views the leaderboard. With ZSets, `ZREVRANGE 0 9` always returns the top 10 in O(log N + 10) — no matter how many total players.

2. **What "sorted by insertion" really means** — ZSets maintain sort order as scores are updated. When alice's score changes, her position is automatically adjusted. No re-sort needed.

3. **When would you NOT use a ZSet for a leaderboard?** — ZSets work for global ranking. What if you need "rank among friends only"? Or rank by multiple criteria (score, then by time)?

## Quality Check

Your response should make your PM say "Oh, I see — it's not just about sorting, it's about keeping the sort efficient as scores change."
