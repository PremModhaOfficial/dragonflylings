## Hint 1

There are two bugs: one in `ExpiryChannel()` and one in `WatchExpiredKeys()`.

In `ExpiryChannel`: the channel type is wrong. Redis has two notification namespaces:
- `__keyspace@<db>__:<command>` — receives the command name when it runs on a specific key
- `__keyevent@<db>__:<event>` — receives the KEY NAME when a specific event type occurs

For expiry monitoring (which key just expired?), you want **keyevent**, not keyspace.

## Hint 2

In `WatchExpiredKeys`: the database index is wrong. Dragonfly uses `@0` for the default database (db 0). The channel name contains `@1` which watches database 1, but keys live in database 0.

## Hint 3

```go
func ExpiryChannel() string {
    return "__keyevent@0__:expired"  // keyevent, db 0
}

func WatchExpiredKeys(...) (*redis.PubSub, error) {
    channel := "__keyevent@0__:expired"  // same fix: @1 -> @0
    ...
}
```
