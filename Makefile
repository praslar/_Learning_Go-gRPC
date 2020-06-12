PROJECT_NAME=gogodo
BUILD_VERSION=$(shell cat VERSION)
GO_BUILD_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on
GO_FILES=$(shell go list ./... | grep -v /vendor/)
DOCKER_IMAGE=$(PROJECT_NAME):$(BUILD_VERSION)
DOCKER_IMAGE_CLIENT=$(PROJECT_NAME)-client:$(BUILD_VERSION)

.SILENT:

all: mod_tidy fmt vet test install 

build:
	$(GO_BUILD_ENV) go build -v -o $(PROJECT_NAME)-$(BUILD_VERSION).bin ./cmd/server/main.go

build-client:
	$(GO_BUILD_ENV) go build -v -o $(PROJECT_NAME)-client-$(BUILD_VERSION).bin ./cmd/client-grpc/main.go

install:
	$(GO_BUILD_ENV) go install

vet:
	$(GO_BUILD_ENV) go vet $(GO_FILES)

fmt:
	$(GO_BUILD_ENV) go fmt $(GO_FILES)

mod_tidy:
	$(GO_BUILD_ENV) go mod tidy

test:
	$(GO_BUILD_ENV) go test $(GO_FILES) -cover -v

compose_prod: docker
	cd deployment/docker && BUILD_VERSION=$(BUILD_VERSION) docker-compose up

protoc:
	protoc --proto_path=api/proto/v1 --proto_path=third-party --go_out=plugins=grpc:pkg/api/v1 todo-service.proto
	protoc --proto_path=api/proto/v1 --proto_path=third-party --grpc-gateway_out=logtostderr=true:pkg/api/v1 todo-service.proto
	protoc --proto_path=api/proto/v1 --proto_path=third-party --swagger_out=logtostderr=true:api/swagger/v1 todo-service.proto

docker_prebuild: vet build build-client
	mkdir -p deployment/docker/client
	mv $(PROJECT_NAME)-client-$(BUILD_VERSION).bin deployment/docker/client/$(PROJECT_NAME)-client.bin; \
	mv $(PROJECT_NAME)-$(BUILD_VERSION).bin deployment/docker/$(PROJECT_NAME).bin; \
	cp -R configs/ssl/ca.crt deployment/docker/client;
	cp -R configs/ssl/server.crt deployment/docker/;
	cp -R configs/ssl/server.pem deployment/docker/;
	
docker_build:
	cd deployment/docker; \
	docker build -t $(DOCKER_IMAGE) .;

	cd deployment/docker/client; \
	docker build -t $(DOCKER_IMAGE_CLIENT) .;

docker_postbuild:
	cd deployment/docker; \
	rm -rf $(PROJECT_NAME).bin 2> /dev/null;\

	cd client; \
	rm -rf $(PROJECT_NAME)-client.bin 2> /dev/null;\

docker: docker_prebuild docker_build docker_postbuild

