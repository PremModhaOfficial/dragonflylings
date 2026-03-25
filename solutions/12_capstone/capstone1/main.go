package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type EventPipeline struct {
	client       *redis.Client
	streamKey    string
	leaderboard  string
	groupName    string
	consumerName string
}

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
// ARGV[1] = points (increment), ARGV[2] = player name
const updateLeaderboardScript = `
return redis.call('ZINCRBY', KEYS[1], ARGV[1], ARGV[2])
`

func (p *EventPipeline) PublishScore(ctx context.Context, player string, points int64) (string, error) {
	return p.client.XAdd(ctx, &redis.XAddArgs{
		Stream: p.streamKey,
		Values: map[string]interface{}{
			"player": player,
			"points": strconv.FormatInt(points, 10),
		},
	}).Result()
}

// CreateGroup sets up the consumer group starting from the beginning of the stream ("0").
// Uses MKSTREAM to create the stream if it doesn't exist.
func (p *EventPipeline) CreateGroup(ctx context.Context) error {
	err := p.client.XGroupCreateMkStream(ctx, p.streamKey, p.groupName, "0").Err()
	if err != nil && err.Error() == "BUSYGROUP Consumer Group name already exists" {
		return nil
	}
	return err
}

// ProcessPending reads and processes all undelivered events from the stream.
func (p *EventPipeline) ProcessPending(ctx context.Context) (int, error) {
	processed := 0
	for {
		streams, err := p.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    p.groupName,
			Consumer: p.consumerName,
			Streams:  []string{p.streamKey, ">"},
			Count:    10,
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

				_, err := p.client.Eval(ctx, updateLeaderboardScript,
					[]string{p.leaderboard},
					points, player, // ARGV[1]=points (increment), ARGV[2]=player (member)
				).Result()
				if err != nil {
					return processed, fmt.Errorf("lua error for msg %s: %w", msg.ID, err)
				}

				// Acknowledge the message so it's removed from PEL
				if err := p.client.XAck(ctx, p.streamKey, p.groupName, msg.ID).Err(); err != nil {
					return processed, fmt.Errorf("xack failed for %s: %w", msg.ID, err)
				}

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

func (p *EventPipeline) GetTopPlayers(ctx context.Context, n int) ([]redis.Z, error) {
	return p.client.ZRevRangeWithScores(ctx, p.leaderboard, 0, int64(n-1)).Result()
}

func (p *EventPipeline) GetPendingCount(ctx context.Context) (int64, error) {
	info, err := p.client.XPending(ctx, p.streamKey, p.groupName).Result()
	if err != nil {
		return 0, err
	}
	return info.Count, nil
}
