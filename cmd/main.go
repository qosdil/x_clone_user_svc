package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	gokitGrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	app "github.com/qosdil/x_clone_user_svc"
	configs "github.com/qosdil/x_clone_user_svc/configs"
	grpcSvc "github.com/qosdil/x_clone_user_svc/grpc/service"
	"github.com/qosdil/x_clone_user_svc/repository/databases"
	"github.com/qosdil/x_clone_user_svc/service"
	transport "github.com/qosdil/x_clone_user_svc/transports"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	// Load environment variables
	configs.LoadEnv()

	// Set up MongoDB connection
	mongoURI := configs.GetEnv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database(configs.GetEnv("DB_NAME"))
	repo, err := databases.NewMongoRepository(db)
	if err != nil {
		panic(err)
	}

	var (
		httpAddr = flag.String("http.addr", ":"+configs.GetEnv("HTTP_PORT"), "HTTP listen address")
		grpcAddr = flag.String("grpc-addr", ":"+configs.GetEnv("GRPC_PORT"), "gRPC listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s service.Service
	{
		s = service.NewService(repo)
		s = app.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = transport.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	// The gRPC listener mounts the Gokit gRPC server that we created
	grpcListener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
		os.Exit(1)
	}

	var (
		baseServer = grpc.NewServer(grpc.UnaryInterceptor(gokitGrpc.Interceptor))
		grpcServer = transport.NewGRPCServer(app.MakeServerEndpoints(s), logger)
	)

	go func() {
		defer grpcListener.Close()
		logger.Log("transport", "gRPC", "addr", *grpcAddr)
		grpcSvc.RegisterServiceServer(baseServer, grpcServer)
		errs <- baseServer.Serve(grpcListener)
	}()

	logger.Log("exit", <-errs)
}
