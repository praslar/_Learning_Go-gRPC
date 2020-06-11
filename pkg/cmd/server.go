package cmd

import (
	"context"
	"flag"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"

	v1 "github.com/praslar/to-do-list-micro/pkg/app/v1"
	db "github.com/praslar/to-do-list-micro/pkg/database"
	server "github.com/praslar/to-do-list-micro/pkg/protocol/grpc"
	"github.com/sirupsen/logrus"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.String("port", "50051", "The server port")
)

func RunServer() error {
	flag.Parse()
	// Get configuration
	// Check start server with or without tls
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = testdata.Path("server.pem")
		}
		if *keyFile == "" {
			*keyFile = testdata.Path("server.key")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			logrus.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	// Start new mongodb session
	s, err := db.DialDefaultMongoDB()
	if err != nil {
		logrus.Errorf("Fail to connect database, %v", err)
		return err
	}
	repo := v1.NewMongoRepository(s)
	// Register gogodo apis
	api := v1.NewToDoServiceServer(repo)
	ctx := context.Background()

	return server.RunServer(ctx, opts, api, *port)
}
