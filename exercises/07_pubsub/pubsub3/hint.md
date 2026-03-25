## Hint 1

This exercise has an unusual goal: the "correct" code intentionally loses messages. The test passes when ALL messages are missed. The broken code causes a timeout instead of a clean miss. Study the current code: what race condition causes the timeout?

## Hint 2

The bug is that `sub.Receive(ctx)` blocks waiting for the subscription confirmation, but the context may expire before any real messages arrive (because the messages were already published). The subscription structure isn't quite right. Look at the order: publish → subscribe → receive. This should work but there's still an issue with the confirmation timing.

## Hint 3

The code is almost right. The issue is that `sub.Receive(ctx)` might block if the subscription confirmation hasn't arrived yet and the context has already timed out (the messages were published BEFORE subscribe). Try adding a small sleep after publish before subscribing, OR ensure the Receive call handles the timeout gracefully. The key insight: when context expires, `ReceiveMessage` returns an error, which correctly breaks the loop -- and received remains empty.
