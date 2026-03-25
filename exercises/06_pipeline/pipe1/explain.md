# Explain It — pipe1

Your colleague says: "Pipelining is just batching. Why not just use MSET? It's one command."

Write your response in `feynman/explanations/pipe1.md`.

Cover:
- What MSET can do that a pipeline cannot (and vice versa)
- Whether pipeline guarantees atomicity
- What "round trip" means and why it dominates Redis latency
- A scenario where pipeline saves 10x time vs. individual commands

QUALITY CHECK: Include a concrete number. If RTT is 1ms and you send 100 commands, what is the total time with vs. without pipelining?
