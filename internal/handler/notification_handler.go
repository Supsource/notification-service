package handler

import (
	"net/http"
	"notification-service/internal/model"
	"notification-service/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *service.NotificationService
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

	err := h.service.CreateNotification(
		req.UserID,
		model.NotificationType(req.Type),
		req.Title,
		req.Body,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "failed to create notification"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "notification created",
	})
}
