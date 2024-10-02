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

func (mw loggingMiddleware) GetList(ctx context.Context) (users []UserResponse, err error) {
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
