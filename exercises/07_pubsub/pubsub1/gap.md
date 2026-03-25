# Gap Notebook — pubsub1

After completing this exercise, add to `feynman/gap_notebook.md`:

1. go-redis auto-reconnects subscriptions after network failures -- but are missed messages replayed?
2. What happens if ReceiveMessage is called on a closed PubSub?
3. A subscribed connection can only do PubSub commands. What error do you get if you call GET on it?
4. Is there a way to know how many active subscribers a channel has before publishing?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

- `PUBSUB NUMSUB channel` returns the subscriber count for a channel. Is this count exact or eventually consistent? Can it be stale?
- If you publish to a channel with zero subscribers, PUBLISH returns 0. Is there any Redis-native way to buffer messages for late subscribers without switching to Streams?
- A new subscriber on a connection receives a confirmation message before real messages. What type does go-redis return for that confirmation — how do you detect it?
