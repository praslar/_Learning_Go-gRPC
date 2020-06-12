package cmd

import (
	"context"
	"flag"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"

	v1 "github.com/praslar/to-do-list-micro/pkg/app/v1"
	db "github.com/praslar/to-do-list-micro/pkg/database"
	logger "github.com/praslar/to-do-list-micro/pkg/logger"
	server "github.com/praslar/to-do-list-micro/pkg/protocol/grpc"

	"github.com/praslar/to-do-list-micro/pkg/protocol/rest"
	"github.com/sirupsen/logrus"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.String("port", "50051", "The server port")
	httpPort = flag.String("http_port", "8080", "The server rest port")
)

type Config struct {
	LogLevel int
	// LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00
	LogTimeFormat string
}

func RunServer() error {
	// Get configuration
	var cfg Config
	flag.IntVar(&cfg.LogLevel, "log-level", 0, "Global log level")
	flag.StringVar(&cfg.LogTimeFormat, "log-time-format", "",
		"Print time format for logger e.g. 2006-01-02T15:04:05Z07:00")
	flag.Parse()
	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}
	// Check start server with or without tls
	var opts []grpc.ServerOption
	restOtps := []grpc.DialOption{grpc.WithInsecure()}
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
		restOtps = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
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
	go func() {
		_ = rest.RunServer(ctx, restOtps, *httpPort, *port)
	}()
	return server.RunServer(ctx, opts, api, *port)
}
