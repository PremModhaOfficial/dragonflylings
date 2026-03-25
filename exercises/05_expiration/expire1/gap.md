# Gap Notebook — expire1

After completing this exercise, add to `feynman/gap_notebook.md`:

1. Redis uses lazy deletion + periodic sampling. Could a key still consume memory 60 seconds after its TTL fires?
2. What happens if you call EXPIRE on a key that doesn't exist?
3. What's the difference between `SET key value EX 10` and `SET key value` + `EXPIRE key 10`? Is one safer?
4. If your server clock drifts, does TTL-based expiry become inaccurate?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

- Can you read a key's TTL with millisecond precision? Which command does this? What's the return unit?
- What happens to a key's TTL when you do `SET key newvalue` on a key that already has an expiry? Does the TTL survive?
- What's the difference between `EXPIREAT` (absolute Unix timestamp) and `EXPIRE` (relative seconds)? When does the absolute form matter?
