# Explain It — prod2

Your senior engineer does a code review and says: *"This distributed lock implementation is missing three things that would bite you in production."*

Write your answer in `feynman/explanations/prod2.md`.

Identify and explain at least three production concerns with distributed locking:

1. **Lock expiry vs lock extension:** Your lock TTL is 5 seconds but the operation takes 8 seconds. The lock expires while you're still working. What do you do?

2. **Fencing tokens:** Process A holds a lock, the lock expires, Process B acquires it. Process A's network was slow and it resumes. Now both A and B think they hold the lock. How do fencing tokens help?

3. **Redlock:** Martin Kleppmann wrote a famous critique of Redlock (Redis's multi-node locking algorithm). What's the core argument? Is single-node locking safer or less safe than Redlock?

**Practical question:** For most web apps, is a distributed lock overkill? What's the simplest correct locking strategy for a PostgreSQL-backed app?

QUALITY CHECK: Can you explain the fencing token concept without mentioning Redis?
