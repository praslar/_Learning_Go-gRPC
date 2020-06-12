PROJECT_NAME=gogodo
BUILD_VERSION=$(shell cat VERSION)
GO_BUILD_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on
GO_FILES=$(shell go list ./... | grep -v /vendor/)
DOCKER_IMAGE=$(PROJECT_NAME):$(BUILD_VERSION)

.SILENT:

all: mod_tidy fmt vet test install 

build:
	$(GO_BUILD_ENV) go build -v -o $(PROJECT_NAME)-$(BUILD_VERSION).bin ./cmd/server/main.go

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
	protoc --proto_path=api/proto/v1 --proto_path=third-party --go_out=plugins=grpc:. todo-service.proto

docker_prebuild: vet build
	mv $(PROJECT_NAME)-$(BUILD_VERSION).bin deployment/docker/$(PROJECT_NAME).bin; \
	cp -R configs/ssl/server.crt deployment/docker/;
	cp -R configs/ssl/server.pem deployment/docker/;
docker_build:
	cd deployment/docker; \
	docker build -t $(DOCKER_IMAGE) .;

docker_postbuild:
	cd deployment/docker; \
	rm -rf $(PROJECT_NAME).bin 2> /dev/null;\

docker: docker_prebuild docker_build docker_postbuild

