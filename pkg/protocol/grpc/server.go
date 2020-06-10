package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	v1 "github.com/praslar/to-do-list-micro/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// RunServer run gRPC service server to publish  ToDO services
func RunServer(ctx context.Context, creds credentials.TransportCredentials, api v1.ToDoServiceServer, port string) error {

	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	// Register new grpc Server
	// With creds for SSL
	server := grpc.NewServer(grpc.Creds(creds))
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
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}
