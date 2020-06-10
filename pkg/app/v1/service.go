package v1

import (
	"context"

	v1 "github.com/praslar/to-do-list-micro/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

type (
	repository interface {
	}
	toDoServiceServer struct {
		repo repository
	}
)

// NewToDoServiceServer creates ToDO service
func NewToDoServiceServer(repo repository) v1.ToDoServiceServer {
	return &toDoServiceServer{
		repo: repo,
	}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *toDoServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

func (s *toDoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateRespone, error) {
	return &v1.CreateRespone{
		Api: apiVersion,
		Id:  "",
	}, nil
}

// Read todo task
func (s *toDoServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {

	return &v1.ReadResponse{
		Api:  apiVersion,
		Todo: nil,
	}, nil

}

// Update todo task
func (s *toDoServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	return &v1.UpdateResponse{
		Api:     apiVersion,
		Updated: 0,
	}, nil
}

// Delete todo task
func (s *toDoServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	return &v1.DeleteResponse{
		Api:     apiVersion,
		Deleted: 0,
	}, nil
}

// Read all todo tasks
func (s *toDoServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	return &v1.ReadAllResponse{
		Api:   apiVersion,
		ToDos: nil,
	}, nil
}
