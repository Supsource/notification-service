package main

import (
	"context"
	"log"

	"notification-service/internal/db"
	"notification-service/internal/delivery"
	"notification-service/internal/repository"
	"notification-service/internal/workers"
)

func main() {
	dbPool := db.NewPostgresPool()
	outboxRepo := repository.NewPostgresOutboxRepo(dbPool)
	notificationRepo := repository.NewPostgresNotificationRepo(dbPool)

	factory := delivery.NewFactory(
		delivery.NewEmailSender("", "", "", ""),
		delivery.NewPushSender(),
	)

	worker := workers.NewOutboxWorker(notificationRepo, outboxRepo, factory, log.Default())
	worker.Run(context.Background())
}
