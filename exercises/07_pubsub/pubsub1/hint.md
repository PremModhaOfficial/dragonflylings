## Hint 1

Look at these two lines in Chat:
```go
ctx, cancel := context.WithCancel(context.Background())
cancel() // immediately cancels
```
A cancelled context makes all operations using it fail immediately. What should the context be instead?

## Hint 2

Use `context.WithTimeout` to create a context with a deadline (e.g., 5 seconds). This gives the subscription time to receive messages without the risk of hanging forever.

## Hint 3

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel() // cancel when function returns, not immediately
```

Remove the immediate `cancel()` call and use `defer cancel()` instead.
