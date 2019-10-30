package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/service/v1/serviceinfo/").Handler(httptransport.NewServer(
		e.PostServiceInfoEndpoint,
		decodePostServiceInfoRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/service/v1/serviceinfo/{id}").Handler(httptransport.NewServer(
		e.GetServiceInfoEndpoint,
		decodeGetServiceInfoRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/service/v1/serviceinfo/{id}").Handler(httptransport.NewServer(
		e.PutServiceInfoEndpoint,
		decodePutServiceInfoRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/service/v1/serviceinfo/{id}").Handler(httptransport.NewServer(
		e.DeleteServiceInfoEndpoint,
		decodeDeleteServiceInfoRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodePostServiceInfoRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postServiceInfoRequest
	if e := json.NewDecoder(r.Body).Decode(&req.ServiceInfo); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetServiceInfoRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getServiceInfoRequest{ID: id}, nil
}

func decodePutServiceInfoRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var serviceinfo ServiceInfo
	if err := json.NewDecoder(r.Body).Decode(&serviceinfo); err != nil {
		return nil, err
	}
	return putServiceInfoRequest{
		ID:      id,
		ServiceInfo: serviceinfo,
	}, nil
}

func decodeDeleteServiceInfoRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteServiceInfoRequest{ID: id}, nil
}

func encodePostServiceInfoRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/service/v1/serviceinfo/"
	return encodeRequest(ctx, req, request)
}

func encodeGetServiceInfoRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(getServiceInfoRequest)
	serviceID := url.QueryEscape(r.ID)
	req.URL.Path = "/service/v1/serviceinfo/" + serviceID
	return encodeRequest(ctx, req, request)
}

func encodePutServiceInfoRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PUT").Path("/serviceinfo/{id}")
	r := request.(putServiceInfoRequest)
	serviceID := url.QueryEscape(r.ID)
	req.URL.Path = "/serviceinfo/" + serviceID
	return encodeRequest(ctx, req, request)
}

func encodeDeleteServiceInfoRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("DELETE").Path("/serviceinfo/{id}")
	r := request.(deleteServiceInfoRequest)
	serviceID := url.QueryEscape(r.ID)
	req.URL.Path = "/serviceinfo/" + serviceID
	return encodeRequest(ctx, req, request)
}

func decodePostServiceInfoResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postServiceInfoResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetServiceInfoResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getServiceInfoResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePutServiceInfoResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response putServiceInfoResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteServiceInfoResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response deleteServiceInfoResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
