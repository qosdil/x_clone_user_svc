package x_clone_user_srv

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, user User) (UserResponse, error)
	GetList(ctx context.Context) (users []UserResponse, err error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetList(ctx context.Context) (users []UserResponse, err error) {
	return s.repo.Find(ctx)
}

func (s *service) Create(ctx context.Context, user User) (UserResponse, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	user, err := s.repo.Create(ctx, user)
	if err != nil {
		return UserResponse{}, err
	}
	return UserResponse{
		ID:        user.ID.Hex(),
		Username:  user.Username,
		CreatedAt: user.CreatedAt.T,
	}, nil
}
