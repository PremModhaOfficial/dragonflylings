# Explain It — dragon1

Your teammate says: *"We're using Dragonfly for its multi-threaded performance, but our client code is completely sequential — one command at a time. Are we getting any benefit?"*

Write your answer in `feynman/explanations/dragon1.md`.

Your explanation must cover:
1. What "multi-threaded" means for Dragonfly specifically — which operations benefit?
2. Why a sequential client doesn't fully utilize Dragonfly's parallelism
3. How connection pooling (multiple connections) relates to throughput
4. The practical answer: when does Dragonfly's multi-threading help you even with single-threaded client code?

**The harder question:** If you have 10 goroutines each doing sequential Redis operations on different keys, and another service has 10 goroutines doing the same — are you fully utilizing Dragonfly's potential? What's the theoretical limit?

QUALITY CHECK: Can you draw a timeline diagram showing sequential vs concurrent commands and where Dragonfly's parallelism kicks in?
