package repository

import (
	"context"
	"time"

	"notification-service/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOutboxRepo struct {
	db *pgxpool.Pool
}

func NewPostgresOutboxRepo(db *pgxpool.Pool) *PostgresOutboxRepo {
	return &PostgresOutboxRepo{db: db}
}

func (r *PostgresOutboxRepo) Enqueue(ctx context.Context, n *model.OutboxNotification) error {
	query := `
		INSERT INTO notifications_outbox
		(id, user_id, payload, status, retry_count, next_retry_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, n.ID, n.UserID, n.Payload, n.Status, n.RetryCount, n.NextRetryAt)
	return err
}

func (r *PostgresOutboxRepo) ClaimPending(ctx context.Context, batchSize int, processingTimeout time.Duration) ([]model.OutboxNotification, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	selectQuery := `
		SELECT id, user_id, payload, status, retry_count, next_retry_at, created_at, updated_at
		FROM notifications_outbox
		WHERE (status = $1 OR status = $2)
		  AND next_retry_at <= NOW()
		ORDER BY next_retry_at ASC
		LIMIT $3
		FOR UPDATE SKIP LOCKED
	`

	rows, err := tx.Query(ctx, selectQuery, model.OutboxStatusPending, model.OutboxStatusProcessing, batchSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.OutboxNotification
	var ids []string
	for rows.Next() {
		var n model.OutboxNotification
		if err := rows.Scan(
			&n.ID, &n.UserID, &n.Payload, &n.Status, &n.RetryCount, &n.NextRetryAt, &n.CreatedAt, &n.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, n)
		ids = append(ids, n.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		if err := tx.Commit(ctx); err != nil {
			return nil, err
		}
		return items, nil
	}

	updateQuery := `
		UPDATE notifications_outbox
		SET status = $1, next_retry_at = $2, updated_at = NOW()
		WHERE id = ANY($3)
	`
	processingUntil := time.Now().Add(processingTimeout)
	if _, err := tx.Exec(ctx, updateQuery, model.OutboxStatusProcessing, processingUntil, ids); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PostgresOutboxRepo) MarkSent(ctx context.Context, id string) error {
	query := `UPDATE notifications_outbox SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, model.OutboxStatusSent, id)
	return err
}

func (r *PostgresOutboxRepo) ScheduleRetry(ctx context.Context, id string, retryCount int, nextRetryAt time.Time) error {
	query := `
		UPDATE notifications_outbox
		SET status = $1, retry_count = $2, next_retry_at = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.Exec(ctx, query, model.OutboxStatusPending, retryCount, nextRetryAt, id)
	return err
}

func (r *PostgresOutboxRepo) MarkFailed(ctx context.Context, id string) error {
	query := `UPDATE notifications_outbox SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, model.OutboxStatusFailed, id)
	return err
}

func (r *PostgresOutboxRepo) ListFailed(ctx context.Context, limit, offset int) ([]model.OutboxNotification, error) {
	query := `
		SELECT id, user_id, payload, status, retry_count, next_retry_at, created_at, updated_at
		FROM notifications_outbox
		WHERE status = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, model.OutboxStatusFailed, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.OutboxNotification
	for rows.Next() {
		var n model.OutboxNotification
		if err := rows.Scan(
			&n.ID, &n.UserID, &n.Payload, &n.Status, &n.RetryCount, &n.NextRetryAt, &n.CreatedAt, &n.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PostgresOutboxRepo) RetryFailed(ctx context.Context, ids []string) (int64, error) {
	if len(ids) == 0 {
		query := `
			UPDATE notifications_outbox
			SET status = $1, retry_count = 0, next_retry_at = NOW(), updated_at = NOW()
			WHERE status = $2
		`
		ct, err := r.db.Exec(ctx, query, model.OutboxStatusPending, model.OutboxStatusFailed)
		return ct.RowsAffected(), err
	}

	query := `
		UPDATE notifications_outbox
		SET status = $1, retry_count = 0, next_retry_at = NOW(), updated_at = NOW()
		WHERE status = $2 AND id = ANY($3)
	`
	ct, err := r.db.Exec(ctx, query, model.OutboxStatusPending, model.OutboxStatusFailed, ids)
	return ct.RowsAffected(), err
}
