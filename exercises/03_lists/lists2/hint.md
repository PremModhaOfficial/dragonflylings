# Hints for lists2

## Hint 1 — The Concept

`LRANGE key start stop` returns elements from index `start` to `stop` inclusive. Redis uses **0-based indexing**:
- Index 0 = first element
- Index 1 = second element
- Index -1 = last element
- Index -2 = second to last

For pagination with `pageSize` items per page:
- Page 1: indices `0` to `pageSize-1`
- Page 2: indices `pageSize` to `2*pageSize-1`
- Page N: indices `(N-1)*pageSize` to `N*pageSize-1`

## Hint 2 — The Specific Issue

The broken code calculates `start = page * pageSize`:
- Page 1: start = 1 × 2 = 2 (skips the first 2 items!)
- Page 2: start = 2 × 2 = 4 (skips the first 4 items!)

The fix: `start = (page-1) * pageSize`
- Page 1: start = 0 × 2 = 0 ✓
- Page 2: start = 1 × 2 = 2 ✓

## Hint 3 — Near Solution

```go
func GetPage(client *redis.Client, listKey string, page, pageSize int64) ([]string, error) {
    ctx := context.Background()
    start := (page - 1) * pageSize
    stop := start + pageSize - 1
    return client.LRange(ctx, listKey, start, stop).Result()
}
```
