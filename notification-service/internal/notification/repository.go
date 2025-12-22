package notification

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool}
}

func (r *Repository) Create(ctx context.Context, n *Notification) (*Notification, error) {
	q := `
		INSERT INTO notifications 
		    (sender, receiver, message) 
		VALUES 
		    ($1, $2, $3)
		RETURNING 
			id, sender, receiver, message, created_at, read_at
	`

	err := r.pool.QueryRow(ctx, q, n.Sender, n.Receiver, n.Message).
		Scan(
			&n.ID,
			&n.Sender,
			&n.Receiver,
			&n.Message,
			&n.CreatedAt,
			&n.ReadAt,
		)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (r *Repository) GetByReceiver(ctx context.Context, receiver uuid.UUID) ([]Notification, error) {
	q := `
		SELECT 
    		id, sender, receiver, message, created_at, read_at 
		FROM notifications 
		WHERE receiver = $1
	`

	rows, err := r.pool.Query(ctx, q, receiver)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	notifications := make([]Notification, 0)

	for rows.Next() {
		var n Notification

		err = rows.Scan(
			&n.ID,
			&n.Sender,
			&n.Receiver,
			&n.Message,
			&n.CreatedAt,
			&n.ReadAt,
		)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (r *Repository) GetUserByNotifID(ctx context.Context, notifID uuid.UUID) (uuid.UUID, error) {
	q := `SELECT receiver FROM notifications WHERE id = $1`

	var userID uuid.UUID
	err := r.pool.QueryRow(ctx, q, notifID).Scan(&userID)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (r *Repository) MarkRead(ctx context.Context, id uuid.UUID) error {
	q := `
		UPDATE notifications
		SET read_at = NOW()
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	q := `
		DELETE FROM notifications 
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteAllForUser(ctx context.Context, receiver uuid.UUID) error {
	q := `
		DELETE FROM notifications
		WHERE receiver = $1
	`

	_, err := r.pool.Exec(ctx, q, receiver)
	if err != nil {
		return err
	}

	return nil
}
