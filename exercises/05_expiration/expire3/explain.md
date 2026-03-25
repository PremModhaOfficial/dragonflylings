# Explain It — expire3

A teammate asks: "Instead of polling Redis every second to check if a user's session expired, can we just get notified? How?"

Write your answer in `feynman/explanations/expire3.md`.

Cover:
- What configuration is needed and why it's off by default
- The difference between keyspace and keyevent notification channels
- What payload the subscriber receives (key name? event? both?)
- At least one real use case beyond session expiry

QUALITY CHECK: Include a concrete example showing the channel name and what arrives in the message.
