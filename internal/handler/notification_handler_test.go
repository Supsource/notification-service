package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"notification-service/internal/model"
	"notification-service/internal/service"

	"github.com/gin-gonic/gin"
)

type stubNotificationCreator struct {
	err            error
	receivedType   model.NotificationType
	receivedUserID string
	receivedTitle  string
	receivedBody   string
}

func (s *stubNotificationCreator) CreateNotification(userID string, nType model.NotificationType, title string, body string) error {
	s.receivedUserID = userID
	s.receivedType = nType
	s.receivedTitle = title
	s.receivedBody = body
	return s.err
}

func TestCreateNotificationRejectsUnsupportedType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &stubNotificationCreator{}
	handler := &NotificationHandler{service: svc}

	body := map[string]string{
		"user_id": "123e4567-e89b-12d3-a456-426614174000",
		"type":    "sms",
		"title":   "hello",
		"body":    "world",
	}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/notifications", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req

	handler.CreateNotification(ctx)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateNotificationReturnsConsistentServerErrorShape(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &stubNotificationCreator{err: errors.New("boom")}
	handler := &NotificationHandler{service: svc}

	body := map[string]string{
		"user_id": "123e4567-e89b-12d3-a456-426614174000",
		"type":    "email",
		"title":   "hello",
		"body":    "world",
	}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/notifications", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req

	handler.CreateNotification(ctx)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if _, ok := response["error"]; !ok {
		t.Fatalf("expected error key in response, got %v", response)
	}
}

func TestCreateNotificationPassesNormalizedTypeToService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &stubNotificationCreator{err: service.ErrUnsupportedNotificationType}
	handler := &NotificationHandler{service: svc}

	body := map[string]string{
		"user_id": "123e4567-e89b-12d3-a456-426614174000",
		"type":    "email",
		"title":   "hello",
		"body":    "world",
	}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/notifications", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req

	handler.CreateNotification(ctx)

	if svc.receivedType != model.TypeEmail {
		t.Fatalf("expected normalized type %q, got %q", model.TypeEmail, svc.receivedType)
	}
}
