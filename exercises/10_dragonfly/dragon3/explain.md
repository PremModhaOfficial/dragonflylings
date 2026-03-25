# Explain It — dragon3

Your team is debating: *"Should we use Dragonfly's native JSON type, or just serialize to a string with `json.Marshal` ourselves and use regular SET/GET? What's the actual difference?"*

Write your answer in `feynman/explanations/dragon3.md`.

Your explanation must cover:
1. What you can do with native JSON that you can't do with a serialized string (specific operations)
2. The operational difference: what does Dragonfly store internally for each approach?
3. JSONPath queries (`$.field`, `$.array[0]`, `$..field`) — what problems do they solve?
4. When would you choose plain string serialization over native JSON?

**Practical scenario:** You have 1 million user profiles. You need to find all users with `age > 30`. How would you do this with native JSON? With plain strings? Which is more efficient and why?

**Dragonfly-specific:** Does Dragonfly's JSON support require any modules or configuration? How does this differ from Redis?

QUALITY CHECK: Could you explain JSONPath to someone who knows SQL but not Redis? What's the SQL equivalent of `JSON.GET key $.name`?
