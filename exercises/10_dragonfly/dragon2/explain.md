# Explain It — dragon2

Your infrastructure team says: *"We need to double our Redis server memory to handle BGSAVE spikes. Why does Redis use so much memory during a snapshot?"*

Write your answer in `feynman/explanations/dragon2.md`.

Your explanation must cover:
1. What `fork()` does at the OS level and why copy-on-write matters here
2. The worst-case memory scenario during a Redis BGSAVE
3. How Dragonfly's forkless approach avoids the spike (what does it do instead?)
4. The trade-off: what does Dragonfly's approach cost vs Redis's fork approach?

**Operational question:** Your Redis server is at 8GB of used memory. You need to guarantee BGSAVE can run without OOM-killing the process. How much total RAM do you need? Show your math.

**Dragonfly question:** How would you verify Dragonfly's snapshot truly has no memory spike? What metric would you monitor and how?

QUALITY CHECK: Could a sysadmin who doesn't know Redis internals follow your fork() explanation?
