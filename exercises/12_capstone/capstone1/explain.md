# Explain It — capstone1

You've built a real-time event processing pipeline. Now explain it.

Write your answer in `feynman/explanations/capstone1.md`.

**Part 1: The Architecture**
Draw the data flow from score event to leaderboard. Include:
- Which data structure is used at each step
- What "atomicity" means at each step
- Where things could fail and what happens if they do

**Part 2: The Guarantees**
- What happens if a consumer crashes mid-processing (between XREADGROUP and XACK)?
- What does the PEL (Pending Entries List) guarantee?
- What happens if the Lua script fails — is the leaderboard corrupted?

**Part 3: The Scale**
- Your game has 1 million score events per minute. How does this architecture scale?
- What's the bottleneck? (Stream throughput? Lua execution? Leaderboard updates?)
- How would you add more consumers to parallelize processing?

**Part 4: The Production Concern**
The stream grows forever. How do you bound it? (Hint: XTRIM or XADD MAXLEN)

QUALITY CHECK: Could you whiteboard this architecture in a system design interview in 10 minutes?
