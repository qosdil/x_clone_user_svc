package transport

import (
	"context"

	app "github.com/qosdil/x_clone_user_svc"
	grpcSvc "github.com/qosdil/x_clone_user_svc/grpc/service"
	"github.com/qosdil/x_clone_user_svc/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
)

type GrpcServer struct {
	create        grpctransport.Handler
	getByUsername grpctransport.Handler
}

func (s *GrpcServer) Create(ctx context.Context, req *grpcSvc.CreateRequest) (*grpcSvc.SecureResponse, error) {
	_, rep, err := s.create.ServeGRPC(ctx, req)
	if err == model.ErrCodeUsernameNotAvailable {
		return nil, status.Error(codes.AlreadyExists, model.ErrCodeUsernameNotAvailable.Error())
	}
	if err != nil {
		return nil, err
	}
	return rep.(*grpcSvc.SecureResponse), nil
}

func (s *GrpcServer) GetByUsername(ctx context.Context, req *grpcSvc.GetByUsernameRequest) (*grpcSvc.Response, error) {
	_, rep, err := s.getByUsername.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*grpcSvc.Response), nil
}

func decodeGrpcCreateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*grpcSvc.CreateRequest)
	return app.CreateRequest{
		Username: req.Username, Password: req.Password,
	}, nil
}

func decodeGrpcGetByUsernameRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*grpcSvc.GetByUsernameRequest)
	return app.GetByUsernameRequest{Username: req.Username}, nil
}

func encodeGrpcGetByUsernameResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(app.Response)
	return &grpcSvc.Response{
		Id: resp.User.ID, Username: resp.User.Username,
		Password: resp.User.Password, CreatedAt: resp.User.CreatedAt}, nil
}

func encodeGrpcCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(app.SecureResponse)
	return &grpcSvc.SecureResponse{Id: resp.User.ID, Username: resp.User.Username, CreatedAt: resp.User.CreatedAt}, nil
}

func NewGRPCServer(endpoints app.Endpoints, logger log.Logger) grpcSvc.ServiceServer {
	return &GrpcServer{
		create: grpctransport.NewServer(
			endpoints.CreateEndpoint,
			decodeGrpcCreateRequest,
			encodeGrpcCreateResponse,
		),
		getByUsername: grpctransport.NewServer(
			endpoints.GetByUsernameEndpoint,
			decodeGrpcGetByUsernameRequest,
			encodeGrpcGetByUsernameResponse,
		),
	}
}
