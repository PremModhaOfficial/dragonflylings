# Explain It — strings3

## The Challenge

A colleague says: "I implemented a distributed lock like this — it works fine in testing but sometimes breaks in production."

```go
exists, _ := client.Exists(ctx, lockKey).Result()
if exists == 0 {
    client.Set(ctx, lockKey, ownerID, ttl)
    return true, nil
}
return false, nil
```

Write a code review explaining the bug in `feynman/explanations/strings3.md`.

## Your Explanation Should Cover

1. **The race condition between Exists and Set** — what can happen between those two lines? Draw a timeline with two concurrent goroutines.

2. **Why SETNX is atomic and this code is not** — at what level is SETNX atomic? (Hint: Redis processes commands one at a time.)

3. **What the production failure looks like** — when would this race condition actually trigger? Is it rare? What's the consequence when two processes both "acquire" the same lock?

## Quality Check

Your explanation should make your colleague say "Oh no, I see it now" within 30 seconds of reading it. Include the word "race condition" and explain it without jargon.
