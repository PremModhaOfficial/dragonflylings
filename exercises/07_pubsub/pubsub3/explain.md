# Explain It — pubsub3

Your team is choosing between Pub/Sub and Streams for a notification system.

Write your recommendation in `feynman/explanations/pubsub3.md`.

The system requirements:
- Millions of events per day
- Each event must be processed by exactly one consumer
- Consumer pods restart frequently (Kubernetes)
- Some events are critical and must not be lost

Answer:
- Which would you choose and why?
- What specific Pub/Sub property makes it wrong for this use case?
- What Stream features address the requirements?
- Is there any use case where Pub/Sub is actually BETTER than Streams?

QUALITY CHECK: Be specific -- name the exact Stream commands (XREADGROUP, XACK, XPENDING) that solve each requirement.
