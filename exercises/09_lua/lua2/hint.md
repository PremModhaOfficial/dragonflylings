## Hint 1

The `EVAL` command signature is:
```
EVAL script numkeys key [key ...] arg [arg ...]
```

`numkeys` tells Redis how many of the following tokens are **keys** (going into `KEYS[]`) vs **arguments** (going into `ARGV[]`). Redis reads this number *before* running the script.

Right now the call uses `[]string{}` (zero keys declared). What should it declare?

## Hint 2

Two changes needed, both pointing in the same direction:

**In the Lua script:** Replace `ARGV[1]` (fromKey) with `KEYS[1]` and `ARGV[2]` (toKey) with `KEYS[2]`. The `amount` stays as `ARGV[1]` since it's data, not a key.

**In Go:** Pass both keys in the keys slice: `[]string{fromKey, toKey}`. Pass only `amount` as the variadic arg.

## Hint 3

Full fix:

```go
// Lua script uses KEYS[1], KEYS[2], ARGV[1] for amount
const atomicTransferScript = `
local from_val = tonumber(redis.call('GET', KEYS[1]) or '0')
local amount   = tonumber(ARGV[1])
if from_val < amount then
  return redis.error_reply('insufficient balance')
end
redis.call('DECRBY', KEYS[1], amount)
redis.call('INCRBY', KEYS[2], amount)
return 1
`

// Go call: keys in keys slice, amount as arg
_, err := client.Eval(ctx, atomicTransferScript, []string{fromKey, toKey}, amount).Int()
```
