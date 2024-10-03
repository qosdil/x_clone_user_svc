package x_clone_user_srv

import (
	"context"
	"time"

	"github.com/go-kit/log"
)

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) Create(ctx context.Context, user User) (UserSecureResponse, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Create", "took", time.Since(begin), "err", nil)
	}(time.Now())
	return mw.next.Create(ctx, user)
}

func (mw loggingMiddleware) GetByUsernamePassword(ctx context.Context, username string, password string) (UserSecureResponse, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetByUsernamePassword", "took", time.Since(begin), "err", nil)
	}(time.Now())
	return mw.next.GetByUsernamePassword(ctx, username, password)
}

func (mw loggingMiddleware) GetList(ctx context.Context) (users []UserSecureResponse, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetList", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetList(ctx)
}

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}
