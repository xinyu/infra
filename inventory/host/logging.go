package host

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type Middleware func(Host) Host

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Host) Host {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Host
	logger log.Logger
}

func (mw loggingMiddleware) PostHostInfo(ctx context.Context, h HostInfo) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostHostInfo", "id", h.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostHostInfo(ctx, h)
}

func (mw loggingMiddleware) GetHostInfo(ctx context.Context, id string) (h HostInfo, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetHostInfo", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetHostInfo(ctx, id)
}

func (mw loggingMiddleware) PutHostInfo(ctx context.Context, id string, h HostInfo) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PutHostInfo", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PutHostInfo(ctx, id, h)
}

func (mw loggingMiddleware) DeleteHostInfo(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteHostInfo", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteHostInfo(ctx, id)
}

