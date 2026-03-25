# Explain It — tx1

A teammate asks: "Why use TxPipelined? Won't two separate commands work just as well?"

Write your answer in `feynman/explanations/tx1.md`.

Cover:
- What can happen between the DecrBy and IncrBy if they're not atomic
- Why MULTI/EXEC prevents other clients from seeing partial state
- Whether MULTI/EXEC prevents a negative balance (spoiler: it does NOT -- explain why)
- What additional mechanism you'd need to prevent overdrafts atomically

QUALITY CHECK: Describe the exact failure scenario: two goroutines, one bank account, both read balance=100, both try to withdraw 80.
