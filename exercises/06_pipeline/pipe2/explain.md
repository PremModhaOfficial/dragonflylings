# Explain It — pipe2

Explain to a teammate why pipeline error handling is counter-intuitive:

Write your explanation in `feynman/explanations/pipe2.md`.

Cover:
- Why `Exec` returns a top-level error even if only ONE command fails
- Why you should almost always ignore the top-level error and check individual commands
- How this behavior is different from a MULTI/EXEC transaction (does tx fail if one cmd fails?)
- A real scenario where this matters: loading user data where some fields might already be the wrong type

QUALITY CHECK: Walk through your explanation with the exact three commands from this exercise.
