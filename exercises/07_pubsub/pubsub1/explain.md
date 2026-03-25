# Explain It — pubsub1

Explain to a teammate why Pub/Sub requires careful goroutine design.

Write your explanation in `feynman/explanations/pubsub1.md`.

Cover:
- Why the subscriber and publisher must use different Redis connections
- Why subscription must be confirmed before publishing (race condition)
- What happens to messages published before the subscription is active
- Why context cancellation matters for long-lived subscriptions

QUALITY CHECK: Include a concrete sequence diagram (even in text form) showing the correct order of operations.
