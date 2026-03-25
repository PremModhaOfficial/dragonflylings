# Gap Notebook — strings1

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **Can Redis keys contain spaces or special characters?**
   Format: `- [ ] [your question] -- [why it matters]`
   Think about: what characters are safe? What about unicode? Does it affect performance?

2. **What's the maximum size of a Redis key? Of a Redis string value?**
   Keys can be very long, but should they be? What's the tradeoff between descriptive keys and memory usage?

3. **Is there a way to list all keys matching a pattern in Redis?**
   And why is `KEYS *` dangerous in production?

## Push Further

- What does `SET key value XX` do? What about `SET key value NX`?
- What does `GETSET` do? When would you use it?
- What's `STRLEN`? Can you think of a use case where knowing the byte length of a string is useful?
