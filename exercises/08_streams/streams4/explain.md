# Explain It — streams4

Explain the full crash-recovery flow using Streams to someone designing a distributed job queue.

Write your explanation in `feynman/explanations/streams4.md`.

Walk through the exact steps:
1. Worker A reads job 42 from the stream → enters whose PEL?
2. Worker A crashes → what happens to job 42?
3. Worker B calls XPENDING → what does it see?
4. Worker B calls XCLAIM on job 42 → what changes?
5. Worker B processes job 42, calls XACK → what changes?

Cover: what is "min-idle-time" and why does it exist? (Hint: what if Worker A is just slow?)

QUALITY CHECK: What is the right idle timeout for your specific use case? How would you decide?
