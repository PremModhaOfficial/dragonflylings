# Explain It — lua1

Your teammate asks: *"Why do we need Lua for CAS? Can't we just do `GET`, check in Go, then `SET` if it matches? That's only 2 commands."*

Write your answer in `feynman/explanations/lua1.md`.

Your explanation must cover:
1. What can happen between your GET and SET if another goroutine (or another server) runs
2. Why application-level read-then-write is never atomic
3. What "atomicity" means at the Redis protocol level
4. When you'd choose CAS over a lock

**Analogy challenge:** Explain the problem using a physical-world analogy. The classic one is two people trying to take the last concert ticket at the same moment — can you come up with a different one?

**Harder question to answer:** CAS and distributed locks both prevent lost updates. When would you choose CAS over a distributed lock? When would you choose the lock over CAS?

QUALITY CHECK: If your explanation requires knowing Redis internals to understand, it's too technical. Could a backend developer who has never used Redis follow your reasoning?
