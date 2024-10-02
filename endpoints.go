package x_clone_user_srv

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	ListEndpoint endpoint.Endpoint
}

type listResponse struct {
	Users []UserResponse `json:"users"`
	Err   error          `json:"err"`
}

func MakeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		p, e := s.GetList(ctx)
		return listResponse{Users: p, Err: e}, nil
	}
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		ListEndpoint: MakeListEndpoint(s),
	}
}
