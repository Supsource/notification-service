package repository

import (
	"context"
	"notification-service/internal/model"

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
