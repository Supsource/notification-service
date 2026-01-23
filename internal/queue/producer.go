package queue

import (
	"context"
	"encoding/json"
)

const NotificationQueue = "notification_queue"

type Producer struct {
	redis RedisClient
}

type RedisClient interface {
	LPush(ctx context.Context, key string, values ...interface{})
}

func NewProducer(redis RedisClient) *Producer {
	return &Producer{redis: redis}
}

func (p *Producer) Enqueue(notificationID string) error {
	job := NotificationJob{
		NotificationID: notificationID,
	}

	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return p.redis.LPush(context.Background(), NotificationQueue, data)
}
