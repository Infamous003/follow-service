package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Infamous003/follow-service/internal/domain"
	"github.com/lib/pq"
)

type FollowRepository struct {
	db *sql.DB
}

func NewFollowRepository(db *sql.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

func (r *FollowRepository) FollowUser(followerID, followeeID int64) error {
	query := `
		INSERT INTO follows (follower_id, followee_id, created_at)
		VALUES ($1, $2, NOW())
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, followerID, followeeID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "follows_follower_id_fkey", "follows_followee_id_fkey":
				return domain.ErrUserNotFound
			case "follows_pkey":
				return domain.ErrAlreadyFollowing
			}
		}
	}
	return nil
}

func (r *FollowRepository) UnfollowUser(followerID, followeeID int64) error {
	query := `
		DELETE FROM follows
		WHERE follower_id = $1 AND followee_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(ctx, query, followerID, followeeID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *FollowRepository) ListFollowers(userID int64) ([]*domain.User, error) {
	query := `
		SELECT u.id, u.username, u.created_at
		FROM users u
		JOIN follows f ON u.id = f.follower_id
		WHERE f.followee_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt); err != nil {
			return nil, err
		}
		followers = append(followers, &user)
	}

	return followers, nil
}
