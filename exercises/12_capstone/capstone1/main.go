package main

// CAPSTONE 1: Real-Time Event Processing Pipeline
//
// Time target: 45-60 minutes
// No hints available. You are the expert now.
//
// CHALLENGE: Build a real-time game score processing pipeline using:
//   - Redis Streams for event ingestion
//   - Consumer Groups for distributed processing
//   - Lua scripting for atomic leaderboard updates
//   - Sorted Sets for the leaderboard
//
// ARCHITECTURE:
//
//   Publishers → XADD → Stream → XREADGROUP → Processors → Lua → Leaderboard
//                                                          ↓
//                                                        XACK
//
// FLOW:
//   1. PublishScore: XADD score events to the stream
//   2. CreateGroup: XGROUP CREATE with MKSTREAM (creates stream if not exists)
//   3. ProcessPending: XREADGROUP, process each event with Lua, XACK when done
//   4. GetTopPlayers: ZREVRANGE leaderboard sorted set
//
// BUGS TO FIX (there are 4):
//
// Bug 1 (CreateGroup): Uses "$" as the start ID, which means the group only sees
//   messages published AFTER group creation. But the test publishes events BEFORE
//   creating the group. Fix: use "0" so the group processes from the beginning of
//   the stream, including messages that already exist.
//
// Bug 2 (ProcessPending): Missing XACK after processing. Events stay in the
//   Pending Entries List (PEL) forever — they'll be redelivered on consumer restart.
//
// Bug 3 (updateLeaderboardScript): The Lua script has ARGV indices swapped.
//   ARGV[1] should be the player name, ARGV[2] should be the points value.
//
// Bug 4 (ProcessPending): Uses XREAD instead of XREADGROUP. XREAD doesn't track
//   per-consumer delivery or support acknowledgment.
//
// Good luck! Read the tests carefully — they tell you exactly what must pass.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// EventPipeline processes game score events from a Redis Stream into a leaderboard.
type EventPipeline struct {
	client       *redis.Client
	streamKey    string
	leaderboard  string
	groupName    string
	consumerName string
}

// NewEventPipeline creates a pipeline. Keys use hashtags to ensure same shard.
func NewEventPipeline(client *redis.Client, gameID string) *EventPipeline {
	tag := fmt.Sprintf("{game:%s}", gameID)
	return &EventPipeline{
		client:       client,
		streamKey:    tag + ":events",
		leaderboard:  tag + ":leaderboard",
		groupName:    "score-processors",
		consumerName: "worker-1",
	}
}

// updateLeaderboardScript atomically adds points to a player's score.
// KEYS[1] = leaderboard key
// ARGV[1] = player name
// ARGV[2] = points to add
// Returns the new score.
//
// BUG 3: ARGV indices are swapped — treats points as player name and vice versa.
const updateLeaderboardScript = `
return redis.call('ZINCRBY', KEYS[1], ARGV[1], ARGV[2])
`

// PublishScore adds a score event to the stream.
func (p *EventPipeline) PublishScore(ctx context.Context, player string, points int64) (string, error) {
	return p.client.XAdd(ctx, &redis.XAddArgs{
		Stream: p.streamKey,
		Values: map[string]interface{}{
			"player": player,
			"points": strconv.FormatInt(points, 10),
		},
	}).Result()
}

// CreateGroup sets up the consumer group.
// BUG 1: Uses "$" — only processes messages after the group was created.
// The test publishes messages BEFORE creating the group, so you need "0"
// to process all messages from the beginning of the stream.
func (p *EventPipeline) CreateGroup(ctx context.Context) error {
	err := p.client.XGroupCreateMkStream(ctx, p.streamKey, p.groupName, "$").Err()
	if err != nil && err.Error() == "BUSYGROUP Consumer Group name already exists" {
		return nil // group already exists — OK
	}
	return err
}

// ProcessPending reads and processes all pending events from the stream.
// Returns the number of events processed.
//
// BUG 2: Missing XACK — processed events stay in PEL and will be redelivered.
// BUG 4: Uses XREAD instead of XREADGROUP — doesn't track delivery per consumer.
func (p *EventPipeline) ProcessPending(ctx context.Context) (int, error) {
	processed := 0
	for {
		// BUG 4: XREAD doesn't support consumer groups or delivery tracking
		streams, err := p.client.XRead(ctx, &redis.XReadArgs{
			Streams: []string{p.streamKey, "0"},
			Count:   10,
			Block:   0,
		}).Result()
		if err == redis.Nil || (err == nil && len(streams) == 0) {
			break
		}
		if err != nil {
			return processed, err
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				player, _ := msg.Values["player"].(string)
				pointsStr, _ := msg.Values["points"].(string)
				points, _ := strconv.ParseInt(pointsStr, 10, 64)

				// Update leaderboard using Lua script
				_, err := p.client.Eval(ctx, updateLeaderboardScript,
					[]string{p.leaderboard},
					player, points, // BUG 3: Lua script has ARGV[1]=player and ARGV[2]=points swapped
				).Result()
				if err != nil {
					return processed, fmt.Errorf("lua error for msg %s: %w", msg.ID, err)
				}

				// BUG 2: Missing XACK — should call:
				// p.client.XAck(ctx, p.streamKey, p.groupName, msg.ID)

				processed++
			}
			if len(stream.Messages) < 10 {
				break
			}
		}
		break
	}
	return processed, nil
}

// GetTopPlayers returns the top N players and their scores from the leaderboard.
func (p *EventPipeline) GetTopPlayers(ctx context.Context, n int) ([]redis.Z, error) {
	return p.client.ZRevRangeWithScores(ctx, p.leaderboard, 0, int64(n-1)).Result()
}

// GetPendingCount returns how many messages are in the consumer's PEL (unacknowledged).
func (p *EventPipeline) GetPendingCount(ctx context.Context) (int64, error) {
	info, err := p.client.XPending(ctx, p.streamKey, p.groupName).Result()
	if err != nil {
		return 0, err
	}
	return info.Count, nil
}
