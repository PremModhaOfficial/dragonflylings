# Gap Notebook — prod1

After completing this exercise, add to `feynman/gap_notebook.md`:

1. `singleflight` deduplicates within a single process. For multi-instance deployments, what pattern deduplicates across instances?
2. What happens to waiting goroutines if the singleflight call panics? Is there panic recovery built in?
3. Cache stampede can also be prevented with probabilistic early expiration (PER). What is it and when is it better than singleflight?
4. If `fetch()` is slow (500ms), all 50 goroutines wait 500ms. Is there a way to serve stale data while refreshing in background?
5. What's the difference between `singleflight.Do` and `singleflight.DoChan`? When would you use each?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

1. Singleflight deduplicates within one process. For 10 service replicas, a cache miss still triggers 10 origin fetches. Design a Redis-based distributed singleflight using a lock + cached-result pattern — handle the case where the lock holder crashes before writing the result.
2. Cache-aside is a pull pattern (app fetches on miss). Write-through is a push pattern (app writes to cache and DB together on every write). Extend the current implementation to support write-through: what's the atomicity challenge, and how does it affect cache consistency?
3. All callers waiting on singleflight get the same result. If the origin fetch takes 2 seconds and you have mixed-priority callers (user-facing vs background), how do you prevent low-priority work from blocking high-priority callers? Is singleflight still the right tool?
