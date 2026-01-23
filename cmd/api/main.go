package main

import (
	"notification-service/internal/db"
	"notification-service/internal/repository"
	"notification-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	dbPool := db.NewPostgresPool()

	repo := repository.NewPostgresNotificationRepo(dbPool)
	notificationService := service.NewNotificationService(repo)

	r := gin.Default()

	// route will be added later
	_ = notificationService

	r.Run(":8080")
}
