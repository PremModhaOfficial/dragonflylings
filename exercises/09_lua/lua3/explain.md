# Explain It — lua3

Your tech lead asks during code review: *"Why bother with EVALSHA? We already use EVAL and it works fine."*

Write your answer in `feynman/explanations/lua3.md`.

Your explanation must cover:
1. The network overhead difference between EVAL and EVALSHA at scale (quantify it)
2. What happens to the script cache on server restart — and what your code must do about it
3. The `NOSCRIPT` error: when does it happen and how should production code handle it?
4. When would you choose EVAL over EVALSHA?

**Production question:** Your deployment process restarts Dragonfly. Your service still has the old SHA cached. What happens on the first request after restart? Write a code sketch of the `NOSCRIPT` fallback pattern.

QUALITY CHECK: Can you explain the NOSCRIPT fallback without showing code? Could you draw it as a flowchart?
