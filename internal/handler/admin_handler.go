package handler

import (
	"net/http"
	"strconv"

	"notification-service/internal/repository"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	outboxRepo repository.OutboxRepository
}

func NewAdminHandler(outboxRepo repository.OutboxRepository) *AdminHandler {
	return &AdminHandler{outboxRepo: outboxRepo}
}

type RetryFailedRequest struct {
	IDs []string `json:"ids"`
}

func (h *AdminHandler) ListFailedNotifications(c *gin.Context) {
	limit := parseIntQuery(c, "limit", 50)
	offset := parseIntQuery(c, "offset", 0)

	items, err := h.outboxRepo.ListFailed(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load failed notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *AdminHandler) RetryFailedNotifications(c *gin.Context) {
	var req RetryFailedRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := h.outboxRepo.RetryFailed(c.Request.Context(), req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retry notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"retried": count})
}

func parseIntQuery(c *gin.Context, key string, def int) int {
	raw := c.Query(key)
	if raw == "" {
		return def
	}
	val, err := strconv.Atoi(raw)
	if err != nil || val < 0 {
		return def
	}
	return val
}
