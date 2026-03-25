# Gap Notebook — prod4

After completing this exercise, add to `feynman/gap_notebook.md`:

1. Hash vs JSON string for sessions: what's the memory and performance difference?
2. If you need to list all sessions for a user (e.g., "show all active devices"), how would you add an index?
3. `HSCAN` vs `HGETALL` — when does `HSCAN` matter for session data?
4. Session fixation attack: what is it and how does storing sessions in Redis help or not help?
5. If your session contains sensitive data (auth tokens), should it be encrypted at rest? How would you implement this?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

1. Session rotation on privilege escalation: when a user logs in or elevates privileges, the old session must be atomically invalidated and a new one issued. Write the Redis operations that do this safely — specifically, how do you prevent a race where the old token is used between DELETE and SET?
2. If Redis goes down, all users are logged out. Design a graceful-degradation strategy using signed cookies: when Redis is unavailable, the app trusts a short-lived cookie for read-only session data. What are the security trade-offs vs requiring Redis?
3. GDPR requires deleting all data for a user on request. If sessions are stored as `session:{token}` with no user index, how do you find and purge all sessions for `user:42`? Design the index structure that makes this O(1) to enumerate, and keep it consistent with session expiry.
