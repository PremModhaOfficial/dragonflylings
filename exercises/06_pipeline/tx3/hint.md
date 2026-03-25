## Hint 1

Think about what each scenario actually needs:
- "batch-100-independent-writes": you want speed. Do you need all 100 to succeed together? No -- each is independent.
- "atomic-debit-credit": you need both operations to happen as a single unit. No partial state.
- "server-side-conditional-set": you need a conditional (if X then set Y) that runs without round trips.

## Hint 2

Decision tree:
1. Do you need **speed without atomicity**? → Pipeline
2. Do you need **atomicity across multiple keys**? → Transaction (MULTI/EXEC)
3. Do you need **server-side conditional logic** or **true atomicity with reads**? → Lua script

## Hint 3

```go
case "batch-100-independent-writes":
    return "pipeline"      // speed, no atomicity
case "atomic-debit-credit":
    return "transaction"   // MULTI/EXEC guarantees both happen
case "server-side-conditional-set":
    return "lua"           // if/else logic runs server-side, atomically
```
