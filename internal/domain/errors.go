package domain

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUsernameTaken    = errors.New("username already taken")
	ErrAlreadyFollowing = errors.New("already following this user")
	ErrCannotFollowSelf = errors.New("cannot follow yourself")
	ErrNotFollowing     = errors.New("not following this user")
)
