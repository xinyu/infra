package service

import (
	"context"

	"github.com/xinyu/infra/inventory/host"
)

type MiddlewareService func(Service) Service

func HostMiddleware(hostInfo host.Host) MiddlewareService {
	return func(next Service) Service {
		return &hostMiddleware{
			next:   next,
			hostInfo: hostInfo,
		}
	}
}

type hostMiddleware struct {
	next   Service
	hostInfo host.Host
}

func (mw hostMiddleware) PostServiceInfo(ctx context.Context, h ServiceInfo) (err error) {
	if _, err := mw.hostInfo.GetHostInfo(ctx, h.HostID); err != nil {
		return host.ErrNotFoundID
	}

	return mw.next.PostServiceInfo(ctx, h)
}

func (mw hostMiddleware) PutServiceInfo(ctx context.Context, id string, h ServiceInfo) (err error) {
	if _, err := mw.hostInfo.GetHostInfo(ctx, h.HostID); err != nil {
		return host.ErrNotFoundID
	}

	return mw.next.PutServiceInfo(ctx, id, h)
}

func (mw hostMiddleware) GetServiceInfo(ctx context.Context, id string) (h ServiceInfo, err error) {
	
	return mw.next.GetServiceInfo(ctx, id)
}

func (mw hostMiddleware) DeleteServiceInfo(ctx context.Context, id string) (err error) {
	return mw.next.DeleteServiceInfo(ctx, id)
}
