package service

import (
	"context"
	"strings"

	"github.com/qosdil/x_clone_user_svc/model"
	"github.com/qosdil/x_clone_user_svc/repository"
)

type Service interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)
	GetList(ctx context.Context) (users []model.SecureUser, err error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetByUsername(ctx context.Context, username string) (model.User, error) {
	return s.repo.FirstByUsername(ctx, username)
}

func (s *service) GetList(ctx context.Context) (users []model.SecureUser, err error) {
	return s.repo.Find(ctx)
}

func (s *service) Create(ctx context.Context, user model.User) (model.User, error) {
	user.Username = strings.ToLower(user.Username)
	return s.repo.Create(ctx, user)
}
