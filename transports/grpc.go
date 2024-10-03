package transport

import (
	"context"
	app "x_clone_user_svc"
	grpcSvc "x_clone_user_svc/grpc/service"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
)

type GrpcServer struct {
	create                grpctransport.Handler
	getByUsernamePassword grpctransport.Handler
}

func (s *GrpcServer) Create(ctx context.Context, req *grpcSvc.Request) (*grpcSvc.Response, error) {
	_, rep, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*grpcSvc.Response), nil
}

func (s *GrpcServer) GetByUsernamePassword(ctx context.Context, req *grpcSvc.Request) (*grpcSvc.Response, error) {
	_, rep, err := s.getByUsernamePassword.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*grpcSvc.Response), nil
}

func decodeGrpcRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*grpcSvc.Request)
	return app.Request{
		Username: req.Username, Password: req.Password,
	}, nil
}

func encodeGrpcResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(app.Response)
	return &grpcSvc.Response{Id: resp.User.ID, Username: resp.User.Username, CreatedAt: resp.User.CreatedAt}, nil
}

func NewGRPCServer(endpoints app.Endpoints, logger log.Logger) grpcSvc.ServiceServer {
	return &GrpcServer{
		create: grpctransport.NewServer(
			endpoints.CreateEndpoint,
			decodeGrpcRequest,
			encodeGrpcResponse,
		),
		getByUsernamePassword: grpctransport.NewServer(
			endpoints.GetByUsernamePasswordEndpoint,
			decodeGrpcRequest,
			encodeGrpcResponse,
		),
	}
}
