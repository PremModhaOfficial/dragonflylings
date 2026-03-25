# Explain It — lists2

## The Challenge

You're building an activity feed API. The endpoint `GET /feed?page=1&size=10` should return items 1-10. Your Redis list stores activities with index 0 = most recent.

Write a design doc for your implementation in `feynman/explanations/lists2.md`.

## Your Explanation Should Cover

1. **Why -1 is a valid index** — explain Redis's negative indexing with an analogy. When is `LRANGE 0 -1` useful? When would `LRANGE 0 -11` be useful?

2. **The pagination math** — show the formula for converting page number + page size to start/stop indices. Test your formula with page=1, page=2, page=3 and pageSize=10.

3. **What LRANGE returns when the range is out of bounds** — if a list has 5 items and you request indices 5-9, what happens? Does it error? Return empty? Return partial?

## Quality Check

Your explanation should double as the implementation notes in your team's tech spec. Include the edge case: what happens on the last page when there aren't enough items to fill it?
