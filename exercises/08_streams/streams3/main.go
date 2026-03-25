package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// EXERCISE: streams3 - Consumer Groups
//
// PREDICT: Before writing any code, answer in your head:
//   What does ">" mean in XREADGROUP?
//   What is the difference between creating a group with "0" vs "$"?
//   What happens to a message after XACK is called?
//
// TODO: Fix the three bugs below.

// CreateGroup creates a consumer group on the stream, reading from the very beginning.
// BUG: Uses "$" as the starting ID -- this means the group will only see NEW messages
// added after group creation. Existing messages in the stream are skipped.
// Use "0" to process all existing messages from the start.
func CreateGroup(client *redis.Client, ctx context.Context, stream, group string) error {
	return client.XGroupCreateMkStream(ctx, stream, group, "$").Err() // BUG: should be "0"
}

// ReadGroup reads up to count undelivered messages for the given consumer.
// BUG: Uses "0" as the stream position -- "0" means "re-deliver my pending messages".
// Use ">" to get NEW, undelivered messages.
func ReadGroup(client *redis.Client, ctx context.Context, stream, group, consumer string, count int64) ([]redis.XMessage, error) {
	results, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{stream, "0"}, // BUG: should be ">"
		Count:    count,
		Block:    0,
	}).Result()
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0].Messages, nil
}

// AckMessages acknowledges all messages so they are removed from the PEL.
// BUG: Calls XAck on the wrong key -- uses stream+":ack" instead of stream.
func AckMessages(client *redis.Client, ctx context.Context, stream, group string, ids ...string) error {
	return client.XAck(ctx, stream+":ack", group, ids...).Err() // BUG: should be stream
}

func main() {}
