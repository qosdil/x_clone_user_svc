package x_clone_user_svc

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, user User) (UserSecureResponse, error)
	GetByUsernamePassword(ctx context.Context, username string, password string) (UserSecureResponse, error)
	GetList(ctx context.Context) (users []UserSecureResponse, err error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetByUsernamePassword(ctx context.Context, username string, password string) (UserSecureResponse, error) {
	user, err := s.repo.FirstByUsername(ctx, username)
	if err != nil {
		return UserSecureResponse{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return UserSecureResponse{}, errors.New("invalid password")
	}
	return UserSecureResponse{
		ID:        user.ID.Hex(),
		Username:  user.Username,
		CreatedAt: user.CreatedAt.T,
	}, nil
}

func (s *service) GetList(ctx context.Context) (users []UserSecureResponse, err error) {
	return s.repo.Find(ctx)
}

func (s *service) Create(ctx context.Context, user User) (UserSecureResponse, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
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
