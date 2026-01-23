package queue

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

const NotificationQueue = "notification_queue"

type Producer struct {
	redis     RedisClient
	queueName string
}

type RedisClient interface {
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
}

func NewProducer(redis RedisClient, queueName string) *Producer {
	return &Producer{redis: redis, queueName: queueName}
}

func (p *Producer) Enqueue(notificationID string) error {
	job := NotificationJob{
		NotificationID: notificationID,
	}

	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return p.redis.LPush(context.Background(), p.queueName, data).Err()
}
