# Explain It — strings1

## The Challenge

A junior developer asks: "Why does Redis return `redis.Nil` for a missing key instead of just returning an empty string? Wouldn't that be simpler?"

Write your answer in `feynman/explanations/strings1.md`.

## Your Explanation Should Cover

1. **Why nil and empty string are different things** — if a user explicitly SET a key to `""` (empty string), that's a valid value. If the key was never set, that's a different state. How would you distinguish them if missing returned `""`?

2. **How go-redis represents this** — `redis.Nil` is a sentinel error value. This is a design choice: the function signature `Get() (string, error)` uses the error return to signal "key doesn't exist" rather than returning a special string value.

3. **What this means for your application code** — when you call `client.Get()`, you should always check if `err == redis.Nil` before treating it as a real error. Show the pattern.

## Quality Check

Your answer should make someone say "Oh, that's why my code panicked when it called `.Result()` on a missing key — I wasn't checking for redis.Nil."
