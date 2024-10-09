package repository

import (
	"context"

	"github.com/qosdil/x_clone_user_svc/model"
)

// Repository follows gORM convention for the method namings
type Repository interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	Find(ctx context.Context) (users []model.SecureUser, err error)
	FirstByUsername(ctx context.Context, username string) (user model.User, err error)
}
