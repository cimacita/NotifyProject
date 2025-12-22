package userShadow

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool}
}

func (r *Repository) InsertIfNotExists(ctx context.Context, id uuid.UUID) error {
	q := `
		INSERT INTO user_shadow (id, is_deleted)
		VALUES ($1, FALSE)
		ON CONFLICT (id) DO NOTHING
	`
	_, err := r.pool.Exec(ctx, q, id)
	return err
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	q := `
		UPDATE user_shadow 
		SET is_deleted = TRUE 
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, q, id)
	return err
}

func (r *Repository) IsDeleted(ctx context.Context, id uuid.UUID) (bool, error) {
	var isDeleted bool

	q := `
		SELECT is_deleted
		FROM user_shadow 
		WHERE id = $1
	`

	err := r.pool.QueryRow(ctx, q, id).Scan(&isDeleted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return true, nil
		}
		return false, err
	}

	return isDeleted, err
}
