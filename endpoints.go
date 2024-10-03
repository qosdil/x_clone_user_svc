package x_clone_user_srv

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type listResponse struct {
	Users []UserSecureResponse `json:"users"`
	Err   error                `json:"err"`
}

type Endpoints struct {
	CreateEndpoint endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
}

func MakeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Request)
		u, e := s.Create(ctx, User{Username: req.Username, Password: req.Password})
		return CreateResponse{User: u, Err: e}, nil
	}
}

func MakeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		p, e := s.GetList(ctx)
		return listResponse{Users: p, Err: e}, nil
	}
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateEndpoint: MakeCreateEndpoint(s),
		ListEndpoint:   MakeListEndpoint(s),
	}
}
