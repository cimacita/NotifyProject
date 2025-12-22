package user

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

func (r *Repository) Create(ctx context.Context, u *User) (*User, error) {
	q := `
		INSERT INTO users 
    		(email, password) 
		VALUES 
		    ($1, $2)
		RETURNING id, email, password, created_at, deleted_at
	`

	err := r.pool.QueryRow(ctx, q, u.Email, u.Password).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.DeletedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User

	q := `
		SELECT 
		    id, email, password, created_at, deleted_at
		FROM users 
		WHERE email = $1
	`

	err := r.pool.QueryRow(ctx, q, email).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.DeletedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *Repository) Delete(ctx context.Context, userID uuid.UUID) error {
	q := `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, q, userID)
	if err != nil {
		return err
	}

	return nil
}
