package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) PostServiceInfo(ctx context.Context, h ServiceInfo) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostServiceInfo", "id", h.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostServiceInfo(ctx, h)
}

func (mw loggingMiddleware) GetServiceInfo(ctx context.Context, id string) (h ServiceInfo, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetServiceInfo", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetServiceInfo(ctx, id)
}

func (mw loggingMiddleware) PutServiceInfo(ctx context.Context, id string, h ServiceInfo) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PutServiceInfo", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PutServiceInfo(ctx, id, h)
}

func (mw loggingMiddleware) DeleteServiceInfo(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteServiceInfo", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteServiceInfo(ctx, id)
}

