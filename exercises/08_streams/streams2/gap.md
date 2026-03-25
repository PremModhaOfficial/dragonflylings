# Gap Notebook — streams2

After completing this exercise, add to `feynman/gap_notebook.md`:

1. XRANGE returns entries in ascending ID order. What if you want newest-first?
2. Can you use partial timestamps in XRANGE? E.g., "1711234567" without the sequence part?
3. What happens to the stream's XLEN counter when entries are deleted with XDEL?
4. Is XRANGE O(n)? What determines its time complexity?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

- XRANGE is O(n). Can you paginate through a large stream using XRANGE without loading all entries at once? Write the idiomatic cursor-style loop.
- If you delete 500 entries with XDEL from a 1000-entry stream, what does XLEN return? (Try it — the answer may surprise you)
- Can you use partial timestamps in XRANGE (e.g. `"1711234567"` without the sequence part)? What does Redis infer for the sequence?
