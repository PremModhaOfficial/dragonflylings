# Module 03: Lists — "The Queue"

## Mental Model

A Redis List is a double-ended queue (deque). Think of a conveyor belt where you can place items on either end and take from either end. This makes it perfect for job queues, activity feeds, and message buffers.

```
LPUSH pushes to LEFT:    [new | old2 | old1]
RPUSH pushes to RIGHT:   [old1 | old2 | new]

LPOP removes from LEFT:  old1 ← [old2 | old3]
RPOP removes from RIGHT: [old1 | old2] → old3

FIFO queue:  RPUSH (enqueue right) + LPOP (dequeue left)
LIFO stack:  LPUSH (push left)    + LPOP (pop left)
```

## Predict Before Starting

Before writing any code, answer in your head:
1. Which combination of push/pop gives FIFO? Which gives LIFO?
2. What does LRANGE return for indices 0 to -1?
3. What happens when you LPOP from an empty list?
4. What's the difference between LPOP and BLPOP?

Write your predictions in `feynman/gap_notebook.md`.

## Key Concepts

| Command | Description |
|---------|-------------|
| LPUSH / RPUSH | Push to left/right end |
| LPOP / RPOP | Pop from left/right end (non-blocking) |
| BLPOP / BRPOP | Blocking pop — wait until item available or timeout |
| LRANGE key start stop | Get elements in range (0-based, -1 = last) |
| LLEN key | Number of elements |
| LTRIM key start stop | Keep only elements in range, discard rest |

## Exercises

- **lists1**: LPUSH vs RPUSH — build a FIFO queue (not a stack)
- **lists2**: LRANGE pagination — understand 0-based indexing and -1
- **lists3**: BLPOP blocking worker — wait for jobs without polling
- **lists4**: LLEN + LTRIM — bounded "last N items" pattern for activity feeds

## Before You Start

```bash
# Experiment with lists
redis-cli -p 6380 RPUSH myqueue job1 job2 job3
redis-cli -p 6380 LRANGE myqueue 0 -1   # [job1, job2, job3]
redis-cli -p 6380 LPOP myqueue           # job1 (FIFO)
redis-cli -p 6380 LLEN myqueue           # 2
```
