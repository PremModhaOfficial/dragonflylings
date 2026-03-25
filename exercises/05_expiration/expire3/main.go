package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// EXERCISE: expire3 - Eavesdrop on the Janitor
//
// PREDICT: Before writing any code, answer in your head:
//   What is the difference between keyspace and keyevent notification channels?
//   What does the payload of an expiry notification contain -- the key name or the command?
//   Why is the channel named with "@0" in it?
//
// NOTE: Keyspace notifications are pre-configured in this Dragonfly setup.
// The docker-compose starts Dragonfly with notify-keyspace-events already set to "KEx".
// Your job is to subscribe to the RIGHT channel.
//
// TODO: Fix the two bugs below.

// ExpiryChannel returns the pub/sub channel that broadcasts key names as they expire in db 0.
// BUG: "__keyspace@0__:expired" is the keyspace channel (publishes the command "expired" when
// a command runs on a key). "__keyevent@0__:expired" is the keyevent channel (publishes the
// KEY NAME whenever the "expired" event occurs). For expiry monitoring, you want keyevent.
func ExpiryChannel() string {
	return "__keyspace@0__:expired" // BUG: should be __keyevent@0__:expired
}

// WatchExpiredKeys subscribes to expiry notifications and returns the subscription.
// BUG: Uses wrong database index -- db 1 instead of db 0.
func WatchExpiredKeys(subClient *redis.Client, ctx context.Context) (*redis.PubSub, error) {
	// BUG: channel uses @1 (db 1), but keys live in db 0
	channel := "__keyevent@1__:expired"
	sub := subClient.Subscribe(ctx, channel)
	if _, err := sub.Receive(ctx); err != nil {
		sub.Close()
		return nil, err
	}
	return sub, nil
}

func main() {}
