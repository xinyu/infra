package host

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type Endpoints struct {
	PostHostInfoEndpoint   endpoint.Endpoint
	GetHostInfoEndpoint    endpoint.Endpoint
	PutHostInfoEndpoint   endpoint.Endpoint
	DeleteHostInfoEndpoint    endpoint.Endpoint
}

func MakeServerEndpoints(h Host) Endpoints {
	return Endpoints{
		PostHostInfoEndpoint:   MakePostHostInfoEndpoint(h),
		GetHostInfoEndpoint:    MakeGetHostInfoEndpoint(h),
		PutHostInfoEndpoint:    MakePutHostInfoEndpoint(h),
		DeleteHostInfoEndpoint:    MakeDeleteHostInfoEndpoint(h),
	}
}

func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	return Endpoints{
		PostHostInfoEndpoint:   httptransport.NewClient("POST", tgt, encodePostHostInfoRequest, decodePostHostInfoResponse, options...).Endpoint(),
		GetHostInfoEndpoint:    httptransport.NewClient("GET", tgt, encodeGetHostInfoRequest, decodeGetHostInfoResponse, options...).Endpoint(),
		PutHostInfoEndpoint:    httptransport.NewClient("PUT", tgt, encodePutHostInfoRequest, decodePutHostInfoResponse, options...).Endpoint(),
		DeleteHostInfoEndpoint:    httptransport.NewClient("DELETE", tgt, encodeDeleteHostInfoRequest, decodeDeleteHostInfoResponse, options...).Endpoint(),
	}, nil
}

func (e Endpoints) PostHostInfo(ctx context.Context, h HostInfo) error {
	request := postHostInfoRequest{HostInfo: h}
	response, err := e.PostHostInfoEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(postHostInfoResponse)
	return resp.Err
}

func (e Endpoints) GetHostInfo(ctx context.Context, id string) (HostInfo, error) {
	request := getHostInfoRequest{ID: id}
	response, err := e.GetHostInfoEndpoint(ctx, request)
	if err != nil {
		return HostInfo{}, err
	}
	resp := response.(getHostInfoResponse)
	return resp.HostInfo, resp.Err
}

func (e Endpoints) PutHostInfo(ctx context.Context, id string, h HostInfo) error {
	request := putHostInfoRequest{ID: id, HostInfo: h}
	response, err := e.PutHostInfoEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(putHostInfoResponse)
	return resp.Err
}

func (e Endpoints) DeleteHostInfo(ctx context.Context, id string) error {
	request := deleteHostInfoRequest{ID: id}
	response, err := e.DeleteHostInfoEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteHostInfoResponse)
	return resp.Err
}

func MakePostHostInfoEndpoint(s Host) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postHostInfoRequest)
		e := s.PostHostInfo(ctx, req.HostInfo)
		return postHostInfoResponse{Err: e}, nil
	}
}

func MakeGetHostInfoEndpoint(s Host) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getHostInfoRequest)
		h, e := s.GetHostInfo(ctx, req.ID)
		return getHostInfoResponse{HostInfo: h, Err: e}, nil
	}
}

func MakePutHostInfoEndpoint(s Host) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(putHostInfoRequest)
		e := s.PutHostInfo(ctx, req.ID, req.HostInfo)
		return putHostInfoResponse{Err: e}, nil
	}
}

func MakeDeleteHostInfoEndpoint(s Host) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteHostInfoRequest)
		e := s.DeleteHostInfo(ctx, req.ID)
		return deleteHostInfoResponse{Err: e}, nil
	}
}

type postHostInfoRequest struct {
	HostInfo HostInfo
}

type postHostInfoResponse struct {
	Err error `json:"err,omitempty"`
}

func (r postHostInfoResponse) error() error { return r.Err }

type getHostInfoRequest struct {
	ID string
}

type getHostInfoResponse struct {
	HostInfo HostInfo `json:"hostinfo,omitempty"`
	Err     error   `json:"err,omitempty"`
}

func (r getHostInfoResponse) error() error { return r.Err }

type putHostInfoRequest struct {
	ID      string
	HostInfo HostInfo
}

type putHostInfoResponse struct {
	Err error `json:"err,omitempty"`
}

func (r putHostInfoResponse) error() error { return nil }

type deleteHostInfoRequest struct {
	ID string
}

type deleteHostInfoResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteHostInfoResponse) error() error { return r.Err }
