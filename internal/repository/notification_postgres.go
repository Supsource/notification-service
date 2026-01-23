package repository

import (
	"context"
	"notification-service/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresNotificationRepo struct {
	db *pgxpool.Pool
}

func NewPostgresNotificationRepo(db *pgxpool.Pool) *PostgresNotificationRepo {
	return &PostgresNotificationRepo{db: db}
}

func (r *PostgresNotificationRepo) Create(n *model.Notification) error {
	query := `
		INSERT INTO notifications 
		(id, user_id, type, title, body, status, retry_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(
		context.Background(),
		query,
		n.ID,
		n.UserID,
		n.Type,
		n.Title,
		n.Body,
		n.Status,
		n.RetryCount,
	)
	return err
}

func (r *PostgresNotificationRepo) GetByID(id string) (*model.Notification, error) {
	query := `SELECT id, user_id, type, title, body, status, retry_count, error_reason, created_at, updated_at FROM notifications WHERE id = $1`
	row := r.db.QueryRow(context.Background(), query, id)

	var n model.Notification
	err := row.Scan(
		&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.Status, &n.RetryCount, &n.ErrorReason, &n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // or error? usually specific error
		}
		return nil, err
	}
	return &n, nil
}

func (r *PostgresNotificationRepo) UpdateStatus(id string, status model.NotificationStatus, errorReason *string) error {
	query := `UPDATE notifications SET status = $1, error_reason = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(context.Background(), query, status, errorReason, id)
	return err
}

func (r *PostgresNotificationRepo) IncrementRetry(id string) error {
	query := `UPDATE notifications SET retry_count = retry_count + 1, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}
