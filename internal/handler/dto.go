package handler

type CreateNotificationRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
	Type   string `json:"type" binding:"required"`
	Title  string `json:"title" binding:"required"`
	Body   string `json:"body" binding:"required"`
}
