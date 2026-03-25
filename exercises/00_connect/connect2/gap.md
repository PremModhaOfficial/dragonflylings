# Gap Notebook — connect2

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **Does go-redis reconnect automatically if Dragonfly restarts?**
   Format: `- [ ] [your question] -- [why it matters]`
   Think about: if your app runs for days and Dragonfly is restarted once, does your app recover?

2. **What's the difference between DialTimeout and ReadTimeout?**
   The exercise uses DialTimeout. But what if Dragonfly connects but then hangs while responding? Is there a separate timeout for that?

3. **What happens to in-flight commands when a connection drops?**
   If your app sends a SET command and the TCP connection drops mid-flight, does the SET execute? How do you know?

## Push Further

- Look up the go-redis `Options` struct — how many timeout fields are there? What does each protect against?
- What is "connection retry" and does go-redis do it automatically?
- In Kubernetes, pods restart frequently. How would you build a Redis client that handles this gracefully?
