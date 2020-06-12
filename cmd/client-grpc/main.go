package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"

	v1 "github.com/praslar/to-do-list-micro/pkg/api/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

const (
	apiVerison = "v1"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	address  = flag.String("server", "localhost:50051", "gRPC server in format host:port")
)

func main() {
	flag.Parse()

	// Get configuration
	// Check start server with or without tls
	opts := grpc.WithInsecure()

	if *tls {
		if *certFile == "" {
			*certFile = testdata.Path("server.pem")
		}
		cred, err := credentials.NewClientTLSFromFile(*certFile, "")
		if err != nil {
			log.Fatalln("Error generate cred in clients, ", err)
			return
		}
		opts = grpc.WithTransportCredentials(cred)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, opts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// gRPC Deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(t)
	pfx := t.Format(time.RFC3339Nano)

	//Call Create Service
	connectServer := v1.NewToDoServiceClient(conn)

	reqCreate := v1.CreateRequest{
		Api: apiVerison,

		Todo: &v1.ToDo{
			Title:       "First Thing to do(" + pfx + ")",
			Description: "Learn Go and prepare a CV (" + pfx + ")",
			Reminder:    reminder,
		},
	}

	resCreate, err := connectServer.Create(ctx, &reqCreate)
	if err != nil {
		log.Fatalln("Fail to create: ", err)
	}

	log.Printf("Create suscessful with to_do_ID: %v", resCreate)
}
