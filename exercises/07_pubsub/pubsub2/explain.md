# Explain It — pubsub2

When would you use PSUBSCRIBE over SUBSCRIBE in a real system?

Write your answer in `feynman/explanations/pubsub2.md`.

Describe a concrete architecture where pattern subscriptions are the right tool:
- Example: a multi-tenant system where each tenant has their own channel namespace
- What would you have to do instead if you didn't have PSUBSCRIBE?
- What are the downsides of pattern subscriptions at scale (many channels, many patterns)?
- Does matching "a.*" match "a.b.c"? Write your answer before testing it.

QUALITY CHECK: Include the glob syntax rules (*, ?, [abc]) with one example each.
