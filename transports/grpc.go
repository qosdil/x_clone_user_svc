package transport

import (
	"context"
	app "x_clone_user_srv"
	grpcSvc "x_clone_user_srv/grpc/service"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
)

type GrpcServer struct {
	create grpctransport.Handler
}

func (s *GrpcServer) Create(ctx context.Context, req *grpcSvc.CreateRequest) (*grpcSvc.CreateResponse, error) {
	_, rep, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*grpcSvc.CreateResponse), nil
}

func decodeGRPCCreateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*grpcSvc.CreateRequest)
	return app.CreateRequest{
		User: app.User{Username: req.Username, Password: req.Password},
	}, nil
}

func encodeGRPCCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(app.CreateResponse)
	return &grpcSvc.CreateResponse{ID: resp.User.ID, Username: resp.User.Username, CreatedAt: resp.User.CreatedAt}, nil
}

func NewGRPCServer(endpoints app.Endpoints, logger log.Logger) grpcSvc.ServiceServer {
	return &GrpcServer{
		create: grpctransport.NewServer(
			endpoints.CreateEndpoint,
			decodeGRPCCreateRequest,
			encodeGRPCCreateResponse,
		),
	}
}
