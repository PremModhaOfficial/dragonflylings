## Hint 1

`client.Subscribe` subscribes to exact channel names. `client.PSubscribe` subscribes to glob patterns. The "P" stands for Pattern. That's the fix for `SubscribePattern`.

## Hint 2

For `ReceivePatternMessage`: pattern subscription messages arrive as `*redis.PMessage`, not `*redis.Message`. `sub.ReceiveMessage(ctx)` only handles `*redis.Message` and will block forever if only PMessages arrive. Use `sub.Receive(ctx)` instead, which returns `interface{}`.

## Hint 3

```go
func ReceivePatternMessage(sub *redis.PubSub, ctx context.Context) (channel, pattern, payload string, err error) {
    iface, err := sub.Receive(ctx)
    if err != nil {
        return "", "", "", err
    }
    switch msg := iface.(type) {
    case *redis.PMessage:
        return msg.Channel, msg.Pattern, msg.Payload, nil
    default:
        return "", "", "", fmt.Errorf("unexpected message type: %T", iface)
    }
}
```
