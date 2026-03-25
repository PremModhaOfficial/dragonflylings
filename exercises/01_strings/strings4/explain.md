# Explain It — strings4

## The Challenge

Your teammate says: "I don't get why INCR exists. I can just do GET, add 1, then SET. It's the same thing and more flexible."

Write your explanation in `feynman/explanations/strings4.md`.

## Your Explanation Should Cover

1. **What "atomic" means at the Redis level** — Redis processes commands sequentially (one at a time). INCR is one command. GET+SET is two commands. Between those two commands, other clients can run other commands.

2. **A concrete race condition scenario** — draw out two goroutines running GET+SET at the same time. Show the exact sequence of operations that causes a lost update. Use a simple counter starting at 5.

3. **Why this matters in practice** — give a real example: page view counter, API request counter, inventory decrement. What goes wrong in production if you lose 1 in 10,000 increments?

## Quality Check

Your explanation should work as a "pull request review" comment on code that uses GET+SET for a counter. Make it convincing enough that the reviewer changes their code.

Bonus: explain what `INCRBY` and `INCRBYFLOAT` do and when you'd use them.
