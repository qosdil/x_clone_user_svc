package x_clone_user_svc

import (
	"context"
	"time"
	"x_clone_user_svc/model"
	"x_clone_user_svc/service"

	"github.com/go-kit/log"
)

type loggingMiddleware struct {
	next   service.Service
	logger log.Logger
}

func (mw loggingMiddleware) Create(ctx context.Context, user model.User) (model.User, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Create", "took", time.Since(begin), "err", nil)
	}(time.Now())
	return mw.next.Create(ctx, user)
}

func (mw loggingMiddleware) GetByUsernamePassword(ctx context.Context, username string, password string) (model.User, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetByUsernamePassword", "took", time.Since(begin), "err", nil)
	}(time.Now())
	return mw.next.GetByUsernamePassword(ctx, username, password)
}

func (mw loggingMiddleware) GetList(ctx context.Context) (users []model.SecureUser, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetList", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetList(ctx)
}

type Middleware func(service.Service) service.Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next service.Service) service.Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}
