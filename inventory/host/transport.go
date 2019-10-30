package host

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

func MakeHTTPHandler(s Host, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/host/v1/hostinfo/").Handler(httptransport.NewServer(
		e.PostHostInfoEndpoint,
		decodePostHostInfoRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/host/v1/hostinfo/{id}").Handler(httptransport.NewServer(
		e.GetHostInfoEndpoint,
		decodeGetHostInfoRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/host/v1/hostinfo/{id}").Handler(httptransport.NewServer(
		e.PutHostInfoEndpoint,
		decodePutHostInfoRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/host/v1/hostinfo/{id}").Handler(httptransport.NewServer(
		e.DeleteHostInfoEndpoint,
		decodeDeleteHostInfoRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodePostHostInfoRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postHostInfoRequest
	if e := json.NewDecoder(r.Body).Decode(&req.HostInfo); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetHostInfoRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getHostInfoRequest{ID: id}, nil
}

func decodePutHostInfoRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var hostinfo HostInfo
	if err := json.NewDecoder(r.Body).Decode(&hostinfo); err != nil {
		return nil, err
	}
	return putHostInfoRequest{
		ID:      id,
		HostInfo: hostinfo,
	}, nil
}

func decodeDeleteHostInfoRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteHostInfoRequest{ID: id}, nil
}

func encodePostHostInfoRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/host/v1/hostinfo/"
	return encodeRequest(ctx, req, request)
}

func encodeGetHostInfoRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(getHostInfoRequest)
	hostID := url.QueryEscape(r.ID)
	req.URL.Path = "/host/v1/hostinfo/" + hostID
	return encodeRequest(ctx, req, request)
}

func encodePutHostInfoRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PUT").Path("/hostinfo/{id}")
	r := request.(putHostInfoRequest)
	hostID := url.QueryEscape(r.ID)
	req.URL.Path = "/hostinfo/" + hostID
	return encodeRequest(ctx, req, request)
}

func encodeDeleteHostInfoRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("DELETE").Path("/hostinfo/{id}")
	r := request.(deleteHostInfoRequest)
	hostID := url.QueryEscape(r.ID)
	req.URL.Path = "/hostinfo/" + hostID
	return encodeRequest(ctx, req, request)
}

func decodePostHostInfoResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postHostInfoResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetHostInfoResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getHostInfoResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePutHostInfoResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response putHostInfoResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteHostInfoResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response deleteHostInfoResponse
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
