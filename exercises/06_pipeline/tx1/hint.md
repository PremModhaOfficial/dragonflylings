## Hint 1

Read the two lines inside `TxPipelined` carefully. Which key should be decremented (the sender)? Which should be incremented (the receiver)? The variable names `from` and `to` are your guide.

## Hint 2

`DecrBy(ctx, key, amount)` subtracts amount from key. `IncrBy(ctx, key, amount)` adds amount to key. In a transfer, money leaves `from` and arrives at `to`.

## Hint 3

```go
_, err := client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
    pipe.DecrBy(ctx, from, amount) // sender loses money
    pipe.IncrBy(ctx, to, amount)   // receiver gains money
    return nil
})
```
