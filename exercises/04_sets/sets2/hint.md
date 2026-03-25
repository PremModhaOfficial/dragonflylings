# Hints for sets2

## Hint 1 — The Concept

Redis set operations mirror mathematical set theory:

- `SINTER key1 key2` — **intersection**: members present in BOTH sets
- `SUNION key1 key2` — **union**: members present in EITHER set
- `SDIFF key1 key2` — **difference**: members in key1 but NOT in key2

For "common friends" (who do both Alice and Bob follow?): that's intersection — people in BOTH follow lists.

For "people you might know" (who does Bob follow that you don't?): that's difference — `SDIFF bob_follows alice_follows`.

## Hint 2 — The Specific Issue

`CommonFollows` uses `SUnion` — this returns EVERYONE that either user follows, not just the ones in common. Change it to `SInter`.

The other functions are correct:
- `AllFollows` correctly uses `SUnion`
- `UniqueFollows` correctly uses `SDiff`

Only `CommonFollows` has the wrong operation.

## Hint 3 — Near Solution

```go
func CommonFollows(client *redis.Client, key1, key2 string) ([]string, error) {
    ctx := context.Background()
    return client.SInter(ctx, key1, key2).Result()
}
```
