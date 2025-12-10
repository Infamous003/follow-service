package service

import (
	"github.com/Infamous003/follow-service/internal/domain"
	"github.com/Infamous003/follow-service/internal/repository"
)

type FollowService struct {
	followRepo *repository.FollowRepository
	userRepo   *repository.UserRepository
}

func NewFollowService(followRepo *repository.FollowRepository, userRepo *repository.UserRepository) *FollowService {
	return &FollowService{
		followRepo: followRepo,
		userRepo:   userRepo,
	}
}

func (s *FollowService) FollowUser(followerID, followeeID int64) error {
	if followerID == followeeID {
		return domain.ErrCannotFollowSelf
	}

	return s.followRepo.FollowUser(followerID, followeeID)
}

func (s *FollowService) UnfollowUser(followerID, followeeID int64) error {
	return s.followRepo.UnfollowUser(followerID, followeeID)
}

func (s *FollowService) ListFollowers(userID int64) ([]*domain.User, error) {
	if _, err := s.userRepo.GetUserByID(userID); err != nil {
		return nil, err
	}

	followers, err := s.followRepo.ListFollowers(userID)
	if err != nil {
		return nil, err
	}

	return followers, nil
}

func (s *FollowService) ListFollowing(userID int64) ([]*domain.User, error) {
	if _, err := s.userRepo.GetUserByID(userID); err != nil {
		return nil, err
	}

	following, err := s.followRepo.ListFollowing(userID)
	if err != nil {
		return nil, err
	}

	return following, nil
}
