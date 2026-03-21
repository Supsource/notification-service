package handler

import (
	"errors"
	"net/http"
	"notification-service/internal/model"
	"notification-service/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service notificationCreator
}

type notificationCreator interface {
	CreateNotification(userID string, nType model.NotificationType, title string, body string) error
}

func NewNotificationHandler(s *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: s}
}

func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var req CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notificationType, err := model.ParseNotificationType(req.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.CreateNotification(
		req.UserID,
		notificationType,
		req.Title,
		req.Body,
	)

	if err != nil {
		if errors.Is(err, service.ErrUnsupportedNotificationType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create notification"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "notification created",
	})
}
