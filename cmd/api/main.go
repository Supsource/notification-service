package main

import (
	"notification-service/internal/db"
	"notification-service/internal/handler"
	"notification-service/internal/queue"
	"notification-service/internal/repository"
	"notification-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	dbPool := db.NewPostgresPool()

	repo := repository.NewPostgresNotificationRepo(dbPool)
	rdb := queue.NewRedisClient()
	producer := queue.NewProducer(rdb, queue.NotificationQueue)
	notificationService := service.NewNotificationService(repo, producer)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	r := gin.Default()

	r.POST("/notifications", notificationHandler.CreateNotification)
	// route will be added later
	_ = notificationService

	r.Run(":8080")
}
