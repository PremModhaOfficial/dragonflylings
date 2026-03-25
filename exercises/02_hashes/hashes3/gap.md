# Gap Notebook — hashes3

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **Is there HINCRBYFLOAT? When would you use it?**
   Format: `- [ ] [your question] -- [why it matters]`
   Float counters have precision issues. How does Redis store them internally?

2. **What happens when you HINCRBY a field that holds a non-numeric string?**
   What error do you get? How does this differ from INCR on a non-numeric string key?

3. **Can a hash field hold a negative number? Can HINCRBY decrement?**
   What happens if you HINCRBY with -1?

## Push Further

- Design a rate limiter using a hash with `HINCRBY` + `EXPIRE`. What are its limitations compared to sorted sets?
- How would you reset all page counters for a day? Is `HDEL` field by field better or worse than `DEL key`?
- Compare `OBJECT ENCODING analytics:key` before and after the hash grows large. What changes?
