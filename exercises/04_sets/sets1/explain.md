# Explain It — sets1

## The Challenge

You're designing a "filter by tag" feature for a blog. Posts have tags. Users can search for posts with tag "redis" or posts with BOTH "redis" AND "go".

Write your design in `feynman/explanations/sets1.md`.

## Your Explanation Should Cover

1. **Why Sets over Lists for tags** — what breaks if you store tags as a list and someone adds "redis" twice? What does the UI show? How do you deduplicate on read?

2. **SISMEMBER is O(1) — why does that matter?** — checking "does post 42 have tag 'redis'?" with a list requires scanning the whole list (O(N)). With a set, it's constant time. At 10,000 tags per post, the difference is dramatic.

3. **Set operations for filtering** — to find posts with BOTH "redis" AND "go" tags, you'd need the intersection of two sets. Describe how SINTER would power a tag filter.

## Quality Check

Your design should work as an engineering design doc for a small team. Include the Redis key naming scheme you'd use (e.g., `post:{id}:tags`).
