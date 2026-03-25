## Hint 1

Dragonfly assigns keys to shards by hashing the key name. Two keys only land on the same shard if they produce the same hash.

Hashtag syntax: if a key contains `{...}`, **only the content inside the braces** is hashed. So `{user:42}:balance` and `{user:42}:reserved` both hash to `"user:42"` → same shard.

Look at `MakeAccountKeys` — neither key has `{...}`. What's the stable part of the key that both keys share?

## Hint 2

The account ID is what both keys have in common. Wrap it in `{...}` in both key patterns:

```
account:{42}:balance   ← hashes to "42"
account:{42}:reserved  ← hashes to "42" ← same shard ✓
```

Or you can wrap more context:
```
{account:42}:balance   ← hashes to "account:42"
{account:42}:reserved  ← hashes to "account:42" ← same shard ✓
```

Either format works — just ensure both keys use the **same** `{tag}`.

## Hint 3

Fix in `MakeAccountKeys`:

```go
func MakeAccountKeys(accountID string) (balanceKey, reservedKey string) {
    balanceKey  = fmt.Sprintf("{account:%s}:balance", accountID)
    reservedKey = fmt.Sprintf("{account:%s}:reserved", accountID)
    return
}
```
