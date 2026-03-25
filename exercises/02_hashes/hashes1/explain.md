# Explain It — hashes1

## The Challenge

You're reviewing code that stores user profiles as:
```
user:42:name    = "alice"
user:42:email   = "alice@example.com"
user:42:age     = "29"
user:42:role    = "admin"
```

Your colleague says "that works fine." Write your code review in `feynman/explanations/hashes1.md`.

## Your Explanation Should Cover

1. **The keyspace pollution problem** — 10,000 users × 4 fields = 40,000 keys. A hash stores the same data as 10,000 keys. How does this affect `KEYS *`, `SCAN`, memory overhead?

2. **The atomicity problem** — if you want to get all user data, you need 4 GETs. Between the first and fourth GET, another process might update the user. With `HGETALL`, you get a consistent snapshot.

3. **The deletion problem** — deleting a user with string keys requires knowing every field name and calling `DEL` for each. With a hash, `DEL user:42` removes everything.

## Quality Check

Your review should be constructive, not critical. End with: "Here's how I'd change it: [show the hash version]."
