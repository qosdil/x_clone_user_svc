package x_clone_user_srv

import "context"

type Service interface {
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
