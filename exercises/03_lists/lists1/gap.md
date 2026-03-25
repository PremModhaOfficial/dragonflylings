# Gap Notebook — lists1

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **What's the maximum length of a Redis list?**
   Format: `- [ ] [your question] -- [why it matters]`
   Is there a hard limit? What happens to memory if you push endlessly?

2. **What does LPUSH return?**
   It's not just "ok" — it returns a number. What is it? When is it useful?

3. **Can LPUSH push multiple values at once?**
   What's the order when you do `LPUSH mylist a b c`? Is it `[a, b, c]` or `[c, b, a]`?

## Push Further

- What's `LINSERT`? When would you need to insert in the middle of a list?
- What's `LSET`? Can you update a value at a specific index?
- How does a Redis List compare to a Go channel for inter-goroutine communication? What can Redis do that a channel can't?
