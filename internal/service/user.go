package service

import (
	"github.com/Infamous003/follow-service/internal/domain"
	"github.com/Infamous003/follow-service/internal/repository"
	"github.com/Infamous003/follow-service/internal/validator"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(username string) (*domain.User, error) {
	user := &domain.User{
		Username: username,
	}

	v := validator.New()
	validator.ValidateUser(v, user)
	if !v.Valid() {
		return nil, v
	}

	user, err := s.repo.CreateUser(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByID(userID int64) (*domain.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) ListUsers() ([]*domain.User, error) {
	users, err := s.repo.ListUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}
