# Explain It — lua2

Your teammate says: *"Why does Redis care if I pass key names as ARGV vs KEYS? The script gets both either way. It's just a naming convention, right?"*

Write your answer in `feynman/explanations/lua2.md`.

Your explanation must cover:
1. What happens to key routing when Dragonfly sees a script — how does it decide which shard handles it?
2. Why the server **cannot** inspect ARGV to extract key names (hint: what if the key name is computed dynamically inside the script?)
3. What the real-world consequence is in a clustered or sharded environment
4. Why this rule exists even in single-node Redis

**Analogy:** Think of KEYS[] as a restaurant reservation — you must declare your party size before arriving. What goes wrong if you show up with 6 people after reserving for 2?

QUALITY CHECK: The answer should explain a *protocol-level constraint*, not just "Redis requires it." Why does it require it?
