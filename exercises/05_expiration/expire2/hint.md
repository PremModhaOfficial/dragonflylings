## Hint 1

Both functions use the wrong config parameter name. Dragonfly's memory limit is controlled by the config key `maxmemory` — no hyphen, no underscore, one word.

## Hint 2

`GetMemoryLimit` uses `"max-memory"` (hyphen) and `SetMemoryLimit` uses `"max_memory"` (underscore). Both should be `"maxmemory"`.

Try in a Redis CLI: `CONFIG GET maxmemory`

## Hint 3

```go
// Both functions use the same key name:
client.ConfigGet(ctx, "maxmemory")
client.ConfigSet(ctx, "maxmemory", fmt.Sprintf("%d", bytes))
```
