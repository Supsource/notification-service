package main

import (
	"context"
	"encoding/json"
	"log"

	"notification-service/internal/db"
	"notification-service/internal/queue"
	"notification-service/internal/repository"
)

func main() {
	dbPool := db.NewPostgresPool()
	repo := repository.NewPostgresNotificationRepo(dbPool)

	rdb := queue.NewRedisClient()
	producer := queue.NewProducer(rdb, queue.NotificationQueue)
	dlqProducer := queue.NewProducer(rdb, queue.NotificationDLQ)

	for {
		result, err := rdb.BRPop(context.Background(), 0, queue.NotificationQueue).Result()
		if err != nil {
			log.Println("redis error:", err)
			continue
		}

		var job queue.NotificationJob
		if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
			log.Println("invalid job:", err)
			continue
		}

		processNotification(job.NotificationID, repo, producer, dlqProducer)
	}
}
