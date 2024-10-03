package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	app "x_clone_user_svc"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

var (
	ErrAlreadyExists   = errors.New("already exists")
	ErrBadRouting      = errors.New("inconsistent mapping between route and handler (programmer error)")
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrNotFound        = errors.New("not found")
)

type errorer interface {
	error() error
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

func decodeGetListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return nil, nil
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

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func MakeHTTPHandler(s app.Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := app.MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}
	pathPrefix := "/users"
	v1Path := "/v1" + pathPrefix
	r.Methods("GET").Path(v1Path).Handler(httptransport.NewServer(
		e.ListEndpoint,
		decodeGetListRequest,
		encodeResponse,
		options...,
	))
	return r
}
