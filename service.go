package x_clone_user_svc

import (
	"context"
)

type Service interface {
	Create(ctx context.Context, user User) (UserSecureResponse, error)
	GetByUsernamePassword(ctx context.Context, username string, password string) (User, error)
	GetList(ctx context.Context) (users []UserSecureResponse, err error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetByUsernamePassword(ctx context.Context, username string, password string) (User, error) {
	user, err := s.repo.FirstByUsername(ctx, username)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *service) GetList(ctx context.Context) (users []UserSecureResponse, err error) {
	return s.repo.Find(ctx)
}

func (s *service) Create(ctx context.Context, user User) (UserSecureResponse, error) {
	user, err := s.repo.Create(ctx, user)
	if err != nil {
		return UserSecureResponse{}, err
	}
	return UserSecureResponse{
		ID:        user.ID.Hex(),
		Username:  user.Username,
		CreatedAt: user.CreatedAt.T,
	}, nil
}
