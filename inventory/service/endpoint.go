package service

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type Endpoints struct {
	PostServiceInfoEndpoint   endpoint.Endpoint
	GetServiceInfoEndpoint    endpoint.Endpoint
	PutServiceInfoEndpoint   endpoint.Endpoint
	DeleteServiceInfoEndpoint    endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostServiceInfoEndpoint:   MakePostServiceInfoEndpoint(s),
		GetServiceInfoEndpoint:    MakeGetServiceInfoEndpoint(s),
		PutServiceInfoEndpoint:    MakePutServiceInfoEndpoint(s),
		DeleteServiceInfoEndpoint:    MakeDeleteServiceInfoEndpoint(s),
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
		PostServiceInfoEndpoint:   httptransport.NewClient("POST", tgt, encodePostServiceInfoRequest, decodePostServiceInfoResponse, options...).Endpoint(),
		GetServiceInfoEndpoint:    httptransport.NewClient("GET", tgt, encodeGetServiceInfoRequest, decodeGetServiceInfoResponse, options...).Endpoint(),
		PutServiceInfoEndpoint:    httptransport.NewClient("PUT", tgt, encodePutServiceInfoRequest, decodePutServiceInfoResponse, options...).Endpoint(),
		DeleteServiceInfoEndpoint:    httptransport.NewClient("DELETE", tgt, encodeDeleteServiceInfoRequest, decodeDeleteServiceInfoResponse, options...).Endpoint(),
	}, nil
}

func (e Endpoints) PostServiceInfo(ctx context.Context, h ServiceInfo) error {
	request := postServiceInfoRequest{ServiceInfo: h}
	response, err := e.PostServiceInfoEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(postServiceInfoResponse)
	return resp.Err
}

func (e Endpoints) GetServiceInfo(ctx context.Context, id string) (ServiceInfo, error) {
	request := getServiceInfoRequest{ID: id}
	response, err := e.GetServiceInfoEndpoint(ctx, request)
	if err != nil {
		return ServiceInfo{}, err
	}
	resp := response.(getServiceInfoResponse)
	return resp.ServiceInfo, resp.Err
}

func (e Endpoints) PutServiceInfo(ctx context.Context, id string, h ServiceInfo) error {
	request := putServiceInfoRequest{ID: id, ServiceInfo: h}
	response, err := e.PutServiceInfoEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(putServiceInfoResponse)
	return resp.Err
}

func (e Endpoints) DeleteServiceInfo(ctx context.Context, id string) error {
	request := deleteServiceInfoRequest{ID: id}
	response, err := e.DeleteServiceInfoEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteServiceInfoResponse)
	return resp.Err
}

func MakePostServiceInfoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postServiceInfoRequest)
		e := s.PostServiceInfo(ctx, req.ServiceInfo)
		return postServiceInfoResponse{Err: e}, nil
	}
}

func MakeGetServiceInfoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getServiceInfoRequest)
		h, e := s.GetServiceInfo(ctx, req.ID)
		return getServiceInfoResponse{ServiceInfo: h, Err: e}, nil
	}
}

func MakePutServiceInfoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(putServiceInfoRequest)
		e := s.PutServiceInfo(ctx, req.ID, req.ServiceInfo)
		return putServiceInfoResponse{Err: e}, nil
	}
}

func MakeDeleteServiceInfoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteServiceInfoRequest)
		e := s.DeleteServiceInfo(ctx, req.ID)
		return deleteServiceInfoResponse{Err: e}, nil
	}
}

type postServiceInfoRequest struct {
	ServiceInfo ServiceInfo
}

type postServiceInfoResponse struct {
	Err error `json:"err,omitempty"`
}

func (r postServiceInfoResponse) error() error { return r.Err }

type getServiceInfoRequest struct {
	ID string
}

type getServiceInfoResponse struct {
	ServiceInfo ServiceInfo `json:"serviceinfo,omitempty"`
	Err     error   `json:"err,omitempty"`
}

func (r getServiceInfoResponse) error() error { return r.Err }

type putServiceInfoRequest struct {
	ID      string
	ServiceInfo ServiceInfo
}

type putServiceInfoResponse struct {
	Err error `json:"err,omitempty"`
}

func (r putServiceInfoResponse) error() error { return nil }

type deleteServiceInfoRequest struct {
	ID string
}

type deleteServiceInfoResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteServiceInfoResponse) error() error { return r.Err }
