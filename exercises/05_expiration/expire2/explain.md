# Explain It — expire2

A production incident: your Redis/Dragonfly instance OOM-killed a process. Explain what happened and how to prevent it.

Write your explanation in `feynman/explanations/expire2.md`.

Cover:
- What `maxmemory=0` means (hint: it's not "no memory")
- What happens when Dragonfly reaches its memory ceiling with different eviction policies
- Why `noeviction` is dangerous for a primary data store but safe for a pure cache
- What monitoring you would set up to detect approaching maxmemory

QUALITY CHECK: Include the `CONFIG GET maxmemory` command output format and what each value means.
