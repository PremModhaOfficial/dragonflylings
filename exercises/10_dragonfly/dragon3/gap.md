# Gap Notebook — dragon3

After completing this exercise, add to `feynman/gap_notebook.md`:

1. Can you use TTL (EXPIRE) on a JSON key just like any other key?
2. What's `JSON.ARRAPPEND`? When would you use it instead of GET + modify + SET?
3. JSONPath `$..name` (recursive descent) — how does Dragonfly handle deeply nested documents?
4. What's the maximum depth/size of a JSON document Dragonfly supports?
5. Does `JSON.SET key $.newfield value` add a new field to an existing document, or does it fail?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

1. Dragonfly stores JSON natively rather than as a serialized string. What are the memory trade-offs? When would storing JSON as a plain string (with Go-side marshal/unmarshal) actually be more efficient?
2. `JSON.SET key $.field value` is a partial update. Compare it to a GET-modify-SET cycle: what race condition does native partial update avoid, and when does it not eliminate the race?
3. JSONPath filter expressions like `$.items[?(@.price > 10)]` can query arrays server-side. Design a product catalog use case where a server-side filter is dramatically more efficient than loading the whole document — and one where it isn't.
