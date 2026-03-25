# Explain It — streams5

Explain the tradeoff between exact and approximate MAXLEN trimming.

Write your explanation in `feynman/explanations/streams5.md`.

Cover:
- What does `XTRIM stream MAXLEN 1000` guarantee?
- What does `XTRIM stream MAXLEN ~ 1000` guarantee? What does "~" mean?
- Why is approximate trimming O(1) amortized while exact is O(n)?
- In a high-throughput system writing 100k events/sec, which would you use and why?

QUALITY CHECK: Include the actual Redis command syntax, not just the go-redis API. What is the exact RESP protocol command for approximate trimming?
