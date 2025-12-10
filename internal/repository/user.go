package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Infamous003/follow-service/internal/domain"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(username string) (*domain.User, error) {
	query := `
		INSERT INTO users (username, created_at)
		VALUES ($1, NOW())
		RETURNING id, username, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				return nil, domain.ErrUsernameTaken
			}
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(userID int64) (*domain.User, error) {
	query := `
		SELECT id, username, created_at
		FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrUserNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (r *UserRepository) ListUsers() ([]*domain.User, error) {
	query := `
		SELECT id, username, created_at
		FROM users
		ORDER BY id ASC
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
