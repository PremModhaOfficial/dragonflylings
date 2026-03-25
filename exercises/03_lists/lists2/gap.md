# Gap Notebook — lists2

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **Is LRANGE O(N)? What does that mean for large lists?**
   Format: `- [ ] [your question] -- [why it matters]`
   Requesting 10 items from index 0 vs from index 10,000 — is one faster?

2. **Does LRANGE include the stop index?**
   In most programming languages, ranges are exclusive of the stop. Redis is inclusive. When have you been burned by off-by-one errors from this?

3. **How does this compare to SQL LIMIT/OFFSET pagination?**
   SQL OFFSET pagination has a "deep pagination" performance problem. Does LRANGE have the same issue?

## Push Further

- What's `LPOS`? When would you search for a value's position in a list?
- Implement "infinite scroll" using LRANGE. What happens when the user reaches the end?
- Try: what does `LRANGE mylist 0 0` return? And `LRANGE mylist -1 -1`?
