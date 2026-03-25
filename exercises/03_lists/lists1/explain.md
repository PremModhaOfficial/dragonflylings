# Explain It — lists1

## The Challenge

A teammate asks: "I need a job queue where jobs are processed in submission order. I'm using Redis Lists. Should I use LPUSH+LPOP or RPUSH+LPOP?"

Write your answer in `feynman/explanations/lists1.md`.

## Your Explanation Should Cover

1. **Draw the list state** — after pushing job1, job2, job3 with LPUSH. Draw what the list looks like. Which end is the "front"? Which order does LPOP return them?

2. **FIFO vs LIFO** — explain both patterns using a real-world analogy (queue at a coffee shop vs. a stack of plates). Which one does a job queue need?

3. **The two correct implementations** — show both `LPUSH+RPOP` and `RPUSH+LPOP` and explain why both give FIFO. When might you choose one over the other?

## Quality Check

Your explanation should work as the answer to an interview question: "How do you implement a queue using a Redis List?"
