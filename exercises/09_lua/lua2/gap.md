# Gap Notebook — lua2

After completing this exercise, add to `feynman/gap_notebook.md`:

1. What happens if you declare 2 keys in `numkeys` but the script actually accesses a 3rd key hardcoded in Lua — does it error, or just work (but unsafely)?
2. In Redis Cluster, cross-slot key access in Lua raises a `CROSSSLOT` error. Does Dragonfly raise the same error? When?
3. Can you access keys NOT listed in KEYS[] in Redis (not Cluster/Dragonfly)? What's the cost of doing so?
4. Is there a Dragonfly or Redis config that enforces strict KEYS[] compliance as a policy?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

1. Redis Cluster enforces KEYS[] for slot routing. Design a Lua script that atomically operates on keys in two different slots — and explain why that's architecturally wrong even if it works on standalone.
2. `redis.call('KEYS', '*')` inside a Lua script returns all keys. What are the production risks, and why is it banned by some ops teams? (Hint: Lua blocks the server while running.)
3. Write down the exact production scenario where using a hardcoded key inside a Lua script silently breaks without any error — even though it works perfectly in your dev environment.
