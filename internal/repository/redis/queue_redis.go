package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type QueueRedisRepo struct {
	client *redis.Client
}

func NewQueueRedisRepo(client *redis.Client) *QueueRedisRepo {
	return &QueueRedisRepo{client: client}
}

const (
	QueueKey       = "radiohead:waiting_room"
	ActiveUsersKey = "radiohead:active_users"
)

// Enqueue pushes a user into the waiting room queue with the current time as the score
func (r *QueueRedisRepo) Enqueue(ctx context.Context, userID string) error {
	score := float64(time.Now().UnixNano())
	return r.client.ZAdd(ctx, QueueKey, redis.Z{
		Score:  score,
		Member: userID,
	}).Err()
}

// GetRank returns the current 0-indexed position of the user in the queue
func (r *QueueRedisRepo) GetRank(ctx context.Context, userID string) (int64, error) {
	rank, err := r.client.ZRank(ctx, QueueKey, userID).Result()
	if err != nil {
		return 0, err
	}
	return rank, nil
}

// AdmitUsers acts as a Cronjob/Worker running every second to transfer batchSize users from Queue to Active state
func (r *QueueRedisRepo) AdmitUsers(ctx context.Context, batchSize int64) ([]string, error) {
	// Fetch the first batchSize users from the queue (those who arrived earliest)
	users, err := r.client.ZRange(ctx, QueueKey, 0, batchSize-1).Result()
	if err != nil || len(users) == 0 {
		return nil, err
	}

	pipeline := r.client.Pipeline()
	for _, user := range users {
		// Remove from queue
		pipeline.ZRem(ctx, QueueKey, user)
		// Grant access (Session Token) by adding to Active Users with a 15-minute TTL
		pipeline.Set(ctx, ActiveUsersKey+":"+user, "granted", 15*time.Minute)
	}
	
	_, err = pipeline.Exec(ctx)
	return users, err
}