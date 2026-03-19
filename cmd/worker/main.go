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

	factory := delivery.NewFactory(
		&delivery.EmailSender{},
		&delivery.PushSender{},
	)

	worker := workers.NewOutboxWorker(outboxRepo, factory, log.Default())
	worker.Run(context.Background())
}
