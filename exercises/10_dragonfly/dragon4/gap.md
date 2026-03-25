# Gap Notebook — dragon4

After completing this exercise, add to `feynman/gap_notebook.md`:

1. What HA (high availability) options does Dragonfly support if not Sentinel? (DNS failover? K8s operator?)
2. Are there Redis commands that Dragonfly doesn't support at all? How do you find out?
3. OBJECT FREQ and OBJECT IDLETIME — do these work on Dragonfly? What are they used for?
4. Does Dragonfly support Redis replication at all? If yes, when does WAIT return non-zero?
5. What's the Dragonfly equivalent of Redis's `DEBUG SLEEP` for testing timeout behavior?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

1. `OBJECT ENCODING` reveals internal representation (listpack, hashtable, skiplist, etc.). How does encoding affect memory per key? At what threshold does a listpack promote to a hashtable, and how does that affect your memory capacity planning?
2. WAIT returns 0 in standalone Dragonfly. If you deploy Dragonfly with a replica using the Kubernetes operator, does WAIT return non-zero? How would you write a CI test that validates replication lag stays under 100ms?
3. Find one Redis command that behaves differently in Dragonfly (different output, different error, or missing). Explain the implication for a team migrating an existing Redis codebase to Dragonfly without a full regression suite.
