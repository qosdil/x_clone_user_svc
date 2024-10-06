package x_clone_user_svc

import (
	"context"
	"x_clone_user_svc/model"
	"x_clone_user_svc/service"

	"github.com/go-kit/kit/endpoint"
)

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	User UserSecureResponse `json:"user"`
	Err  error              `json:"err"`
}

type UserNotSecureResponse struct {
	User model.User `json:"user"`
	Err  error      `json:"err"`
}

type UserSecureResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt uint32 `json:"created_at"`
}

type listResponse struct {
	Users []UserSecureResponse `json:"users"`
	Err   error                `json:"err"`
}

type Endpoints struct {
	CreateEndpoint                endpoint.Endpoint
	GetByUsernamePasswordEndpoint endpoint.Endpoint
	ListEndpoint                  endpoint.Endpoint
}

func MakeCreateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Request)
		u, e := s.Create(ctx, model.User{Username: req.Username, Password: req.Password})
		return Response{User: UserSecureResponse{
			ID:        u.ID,
			Username:  u.Username,
			CreatedAt: u.CreatedAt,
		}, Err: e}, nil
	}
}

func MakeGetByUsernamePasswordEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Request)
		u, e := s.GetByUsernamePassword(ctx, req.Username, req.Password)
		return UserNotSecureResponse{User: u, Err: e}, nil
	}
}

func MakeListEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		users, err := s.GetList(ctx)
		var respUsers []UserSecureResponse
		for _, user := range users {
			respUsers = append(respUsers, UserSecureResponse{
				ID:        user.ID,
				Username:  user.Username,
				CreatedAt: user.CreatedAt,
			})
		}
		return listResponse{Users: respUsers, Err: err}, nil
	}
}

func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		CreateEndpoint:                MakeCreateEndpoint(s),
		GetByUsernamePasswordEndpoint: MakeGetByUsernamePasswordEndpoint(s),
		ListEndpoint:                  MakeListEndpoint(s),
	}
}
