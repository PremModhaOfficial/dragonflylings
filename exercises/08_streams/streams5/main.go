package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// EXERCISE: streams5 - Bounded Streams
//
// PREDICT: Before writing any code, answer in your head:
//   What happens to a stream if you never trim it?
//   What is the difference between exact MAXLEN and approximate MAXLEN (~)?
//   Why is approximate trimming preferred in production?
//
// TODO: Fix AddAndTrim -- it trims with maxLen=0, deleting everything.

// AddAndTrim appends an entry and trims the stream to at most maxLen entries.
// BUG: Passes 0 as maxLen to XTrimMaxLen -- this deletes ALL entries from the stream.
// The function should use the maxLen parameter passed in.
func AddAndTrim(client *redis.Client, ctx context.Context, stream string, fields map[string]interface{}, maxLen int64) error {
	if err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		ID:     "*",
		Values: fields,
	}).Err(); err != nil {
		return err
	}
	// BUG: hardcoded 0 instead of maxLen -- trims everything
	return client.XTrimMaxLen(ctx, stream, 0).Err()
}

func main() {}
