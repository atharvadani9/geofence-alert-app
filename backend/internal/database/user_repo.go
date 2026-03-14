package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/atharvadani9/geofence-alert-app/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, email, passwordHash, name, role string) (*models.User, error) {
	query := `
		INSERT INTO users (email, password_hash, name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, password_hash, name, role, created_at, updated_at
	`

	user := &models.User{}
	err := r.db.Pool().QueryRow(ctx, query, email, passwordHash, name, role).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("email already exists: %w", err)
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r *UserRepo) getUser(ctx context.Context, email, userId string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, created_at, updated_at
		FROM users
		WHERE (email = $1 OR ($2 != '' AND id = $2::uuid))
	`
	user := &models.User{}
	err := r.db.Pool().QueryRow(ctx, query, email, userId).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return r.getUser(ctx, email, "")
}

func (r *UserRepo) GetUserByID(ctx context.Context, userId string) (*models.User, error) {
	return r.getUser(ctx, "", userId)
}

func (r *UserRepo) UpdateUser(ctx context.Context, userId, role string) error {
	query := `
		UPDATE users
		SET role = $1, updated_at = now()
		WHERE id = $2
	`
	_, err := r.db.Pool().Exec(ctx, query, role, userId)
	return err
}

func (r *UserRepo) DeleteUser(ctx context.Context, userId string) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := r.db.Pool().Exec(ctx, query, userId)
	return err
}
