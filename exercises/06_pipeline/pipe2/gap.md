# Gap Notebook — pipe2

After completing this exercise, add to `feynman/gap_notebook.md`:

1. Pipeline partial failures -- if 50 of 100 commands fail, does Redis guarantee the 50 successes are durable?
2. What is the difference between a "queuing error" and a "runtime error" in a pipeline?
3. In a MULTI/EXEC transaction, if command #50 fails, do commands 1-49 get rolled back?
4. When would you actually want a pipeline to stop on the first error?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

- If a command in a pipeline errors due to wrong type (e.g., INCR on a list key), does Redis rollback the commands that succeeded before it?
- What's the distinction between a "protocol error" (malformed RESP) and a "command error" (wrong type) in terms of pipeline behavior?
- In go-redis, `pipe.Exec()` returns `([]Cmder, error)`. When is the top-level `error` non-nil vs. per-command `.Err()` non-nil? Can both be non-nil simultaneously?
