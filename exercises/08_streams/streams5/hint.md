## Hint 1

Look at this line: `client.XTrimMaxLen(ctx, stream, 0)`. The third argument is `0`. What does trimming to 0 entries mean? Check what the `maxLen` parameter to `AddAndTrim` holds.

## Hint 2

`XTrimMaxLen(ctx, key, maxLen)` keeps the NEWEST `maxLen` entries and removes older ones. Passing `0` removes everything. You want to pass the `maxLen` variable from the function parameter.

## Hint 3

```go
// Fix: use the maxLen parameter, not the literal 0
return client.XTrimMaxLen(ctx, stream, maxLen).Err()

// Bonus: for production, prefer approximate trim (more efficient):
return client.XTrimMaxLenApprox(ctx, stream, maxLen, 0).Err()
```
