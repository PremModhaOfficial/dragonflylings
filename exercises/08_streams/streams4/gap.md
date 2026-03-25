# Gap Notebook — streams4

After completing this exercise, add to `feynman/gap_notebook.md`:

1. What is a "delivery count" in the PEL? How do you use it to detect poison messages?
2. XAUTOCLAIM combines XPENDING + XCLAIM in one call -- when was it added to Redis?
3. What happens if the recovery consumer also crashes before ACKing?
4. In a consumer group with 10 consumers, how do you efficiently check for stale messages across all of them?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

- The PEL entry has a "delivery count" that increments on each XCLAIM. How would you use this to detect and quarantine "poison" messages that always fail processing?
- XAUTOCLAIM combines XPENDING + XCLAIM in a single round trip. What Redis version introduced it, and what does it return that XCLAIM doesn't?
- If you set `MinIdle` for claiming to 30 seconds but your normal processing takes 45 seconds, what will happen in your cluster? How do you choose the right idle threshold?
