# Gap Notebook — streams3

After completing this exercise, add to `feynman/gap_notebook.md`:

1. Two consumers in the same group read from the same stream -- can they get the same message?
2. What happens if a consumer crashes after receiving messages but before calling XACK?
3. Can a consumer in group A and a consumer in group B both receive the same message?
4. What is XAUTOCLAIM and how does it simplify the XPENDING + XCLAIM pattern?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

- XGROUP CREATE with `MKSTREAM` creates the stream if it doesn't exist. Why is this important in a microservice that might start before the producer?
- Two processes use the same consumer name in the same group. What happens to the PEL — do they share it or conflict?
- XREADGROUP returns an empty array when no new messages exist. How would you implement a blocking wait for new group messages? Which argument controls this?
