package x_clone_user_svc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/qosdil/x_clone_user_svc/model"
	"github.com/qosdil/x_clone_user_svc/service"
)

type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetByUsernameRequest struct {
	Username string `json:"username"`
}

type listResponse struct {
	Users []model.SecureUser `json:"users"`
	Err   error              `json:"err"`
}

type Response struct {
	User model.User `json:"user"`
	Err  error      `json:"err"`
}

type SecureResponse struct {
	User model.SecureUser `json:"user"`
	Err  error            `json:"err"`
}

type Endpoints struct {
	CreateEndpoint        endpoint.Endpoint
	GetByUsernameEndpoint endpoint.Endpoint
	ListEndpoint          endpoint.Endpoint
}

func MakeCreateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateRequest)
		u, e := s.Create(ctx, model.User{Username: req.Username, Password: req.Password})
		return SecureResponse{User: model.SecureUser{
			ID:        u.ID,
			Username:  u.Username,
			CreatedAt: u.CreatedAt,
		}, Err: e}, e
	}
}

func MakeGetByUsernamePasswordEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetByUsernameRequest)
		u, e := s.GetByUsername(ctx, req.Username)
		return Response{User: u, Err: e}, nil
	}
}

func MakeListEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		users, err := s.GetList(ctx)
		var respUsers []model.SecureUser
		for _, user := range users {
			respUsers = append(respUsers, model.SecureUser{
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
		CreateEndpoint:        MakeCreateEndpoint(s),
		GetByUsernameEndpoint: MakeGetByUsernamePasswordEndpoint(s),
		ListEndpoint:          MakeListEndpoint(s),
	}
}
