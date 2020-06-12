package rest

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	v1 "github.com/praslar/to-do-list-micro/pkg/api/v1"
	"github.com/praslar/to-do-list-micro/pkg/logger"
	"github.com/praslar/to-do-list-micro/pkg/protocol/rest/middleware"
)

// RunServer run gRPC service server to publish  ToDO services
func RunServer(ctx context.Context, opts []grpc.DialOption, httpPort, grpcPort string) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start server
	gwmux := runtime.NewServeMux()

	if err := v1.RegisterToDoServiceHandlerFromEndpoint(ctx, gwmux, "localhost:"+grpcPort, opts); err != nil {
		logger.Log.Fatal("failed to start HTTP gateway", zap.String("reason", err.Error()))
		return err
	}

	srv := &http.Server{
		Addr: "localhost:" + httpPort,
		Handler: middleware.AddRequestID(
			middleware.AddLogger(logger.Log, gwmux)),
	}

	// Gracefull shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logger.Log.Info("starting HTTP/REST gateway...")
			// signal is a Ctrl+C, handle it
		}
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	// start gRPC server
	log.Println("starting gRPC server on ", srv.Addr)
	return srv.ListenAndServe()
}
