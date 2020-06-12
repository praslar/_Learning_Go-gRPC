package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"

	v1 "github.com/praslar/to-do-list-micro/pkg/api/v1"
	"github.com/praslar/to-do-list-micro/pkg/logger"
	"github.com/praslar/to-do-list-micro/pkg/protocol/grpc/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// RunServer run gRPC service server to publish  ToDO services
func RunServer(ctx context.Context, opts []grpc.ServerOption, api v1.ToDoServiceServer, port string) error {
	listen, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		return err
	}

	opts = middleware.AddLogging(logger.Log, opts)
	// Register new grpc Server
	// With creds for SSL
	// This creds with be complie from use input --certca --...
	server := grpc.NewServer(opts...)
	// Register refelction service for Evans CLI
	reflection.Register(server)
	// Register service
	v1.RegisterToDoServiceServer(server, api)

	// Gracefull shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// signal is a Ctrl+C, handle it
			logger.Log.Warn("shutting down gRPC server...")
			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	logger.Log.Info("starting gRPC server on " + listen.Addr().String())
	return server.Serve(listen)
}
