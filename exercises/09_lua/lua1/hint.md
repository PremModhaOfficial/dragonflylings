## Hint 1

Lua arrays are **1-indexed** (not 0-indexed like most languages). `ARGV[1]` is the first argument you pass, `ARGV[2]` is the second.

Look at how `client.Eval(...)` is called: `client.Eval(ctx, casScript, []string{key}, expected, newVal)`

What does `expected` map to in Lua? What does `newVal` map to?

## Hint 2

The call passes arguments in this order: `..., expected, newVal`.

So in the Lua script:
- `ARGV[1]` = `expected` (the value we're checking against)
- `ARGV[2]` = `newVal` (the value we want to set)

Look at the two uses of ARGV in `casScript`. Which line should use `ARGV[1]`? Which should use `ARGV[2]`?

## Hint 3

The fix is swapping `ARGV[1]` and `ARGV[2]` in the Lua script:

```lua
if current == ARGV[1] then          -- ARGV[1] = expected (compare)
  redis.call('SET', KEYS[1], ARGV[2])  -- ARGV[2] = newVal (set)
  return 1
end
```
