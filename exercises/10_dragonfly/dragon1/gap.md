# Gap Notebook — dragon1

After completing this exercise, add to `feynman/gap_notebook.md`:

1. go-redis has a connection pool. How does pool size interact with Dragonfly's thread count for maximum throughput?
2. Dragonfly processes commands on the same key sequentially (within a shard). What happens when 100 goroutines all write to the same hot key?
3. How does Dragonfly's throughput scale with CPU core count? Is it linear?
4. What's the difference between latency and throughput? Can you have high throughput with high latency?
5. Does Dragonfly pipeline commands internally, or does each command go to network immediately?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

1. Pipelining batches commands to reduce round-trip overhead. Implement a pipelined version of the concurrent SET benchmark and compare throughput. Why does pipelining help more at high network latency than at low latency?
2. go-redis `PoolSize` defaults may not match your workload. Design a benchmark that measures throughput vs pool size to find the saturation point for your hardware — what metric signals you've found it?
3. At extreme concurrency (1000+ goroutines), Go's scheduler and Dragonfly's thread pool both become bottlenecks. Which saturates first on a 4-core machine? How would you instrument each independently to find out?
