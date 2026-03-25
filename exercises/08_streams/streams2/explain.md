# Explain It — streams2

Explain why stream IDs are useful for time-range queries.

Write your explanation in `feynman/explanations/streams2.md`.

Cover:
- How to query "all events from the last 10 minutes" using IDs (hint: the ID encodes time)
- What XREVRANGE does and when you'd prefer it over XRANGE
- How XRANGE COUNT limits results and what you'd do to paginate through millions of entries
- The difference between stream position ("where I'm reading from") and entry count

QUALITY CHECK: Include a concrete example: "to get all events after 2024-01-01 00:00:00 UTC, use ID..."
