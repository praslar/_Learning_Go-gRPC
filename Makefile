PROJECT_NAME=to-do-list
BUILD_VERSION=$(shell cat VERSION)
GO_BUILD_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on

.SILENT:

all: mod_tidy fmt vet test install 

build:
	$(GO_BUILD_ENV) go build -v -o $(PROJECT_NAME)-$(BUILD_VERSION).bin .

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

protoc:
	protoc --proto_path=api/proto/v1 --proto_path=third-party --go_out=plugins=grpc:. todo-service.proto
