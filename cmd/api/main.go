package main

import (
	"notification-service/internal/db"
	"notification-service/internal/handler"
	"notification-service/internal/repository"
	"notification-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	dbPool := db.NewPostgresPool()

	repo := repository.NewPostgresNotificationRepo(dbPool)
	outboxRepo := repository.NewPostgresOutboxRepo(dbPool)
	notificationService := service.NewNotificationService(repo, outboxRepo)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	adminHandler := handler.NewAdminHandler(outboxRepo)

	r := gin.Default()

	r.POST("/notifications", notificationHandler.CreateNotification)
	r.GET("/admin/failed-notifications", adminHandler.ListFailedNotifications)
	r.POST("/admin/retry-failed", adminHandler.RetryFailedNotifications)

	r.Run(":8080")
}
