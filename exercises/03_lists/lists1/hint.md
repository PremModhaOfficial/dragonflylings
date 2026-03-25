# Hints for lists1

## Hint 1 — The Concept

A Redis List is a doubly-linked list. You can push/pop from either end:
- `LPUSH` = push to the **L**eft (front)
- `RPUSH` = push to the **R**ight (back)
- `LPOP` = pop from the **L**eft (front)
- `RPOP` = pop from the **R**ight (back)

For a **FIFO queue** (first in, first out): push to one end, pop from the other.
- `LPUSH` + `RPOP`: push to left, pop from right → FIFO ✓
- `RPUSH` + `LPOP`: push to right, pop from left → FIFO ✓

For a **LIFO stack** (last in, first out): push and pop from the same end.
- `LPUSH` + `LPOP` → LIFO (the bug in this exercise)

## Hint 2 — The Specific Issue

`EnqueueTask` uses `LPush` (prepends to front). `DequeueTask` uses `LPop` (removes from front).

- Push task1: `[task1]`
- Push task2: `[task2, task1]`
- Push task3: `[task3, task2, task1]`
- LPop returns: `task3` (last in, first out — wrong!)

Fix: change `EnqueueTask` to use `RPush` (append to back) instead of `LPush`. Then `LPop` removes from the front and you get FIFO.

## Hint 3 — Near Solution

```go
func EnqueueTask(client *redis.Client, queueKey, task string) error {
    ctx := context.Background()
    return client.RPush(ctx, queueKey, task).Err()
}
```
