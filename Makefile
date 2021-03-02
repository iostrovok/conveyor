CURDIR := $(shell pwd)
GOBIN := $(CURDIR)/bin/
ENV:=GOBIN=$(GOBIN)
DIR:=FILE_DIR=$(CURDIR)/testfiles TEST_SOURCE_PATH=$(CURDIR)
GODEBUG:=GODEBUG=gocacheverify=1
LOCDIR:=$(PWD)
LOADENV:=GO111MODULE=on GONOSUMDB="*" GOPROXY=direct $(ENV) CURDIR=$(CURDIR)

#################### Protobuf section
PROTO_VERSION=3.13.0
# PROTO_ZIP=protoc-$(PROTO_VERSION)-linux-x86_64.zip
PROTO_ZIP=protoc-$(PROTO_VERSION)-osx-x86_64.zip


##
## List of commands:
##

## default:
all: mod deps fmt lint test

all-deps: mod deps

deps:
	@echo "======================================================================"
	@echo 'MAKE: deps...'
	@mkdir -p $(GOBIN)
	$(LOADENV) go get -u golang.org/x/lint/golint@latest
	$(LOADENV) go get -u  github.com/golang/mock/mockgen@latest
	$(LOADENV) go get -u github.com/dizzyfool/genna@latest
	$(LOADENV) go get -u github.com/go-pg/migrations/v7@latest
	$(LOADENV) go get -u golang.org/x/tools/cmd/godoc@latest

test-workers:
	@echo "======================================================================"
	@echo "Run race test for workers"
	cd $(LOCDIR)/workers && $(DIR) $(GODEBUG) go test -cover -race ./

test-item:
	@echo "======================================================================"
	@echo "Run race test for item"
	cd $(LOCDIR)/item && $(DIR) $(GODEBUG) go test -cover -race ./

test-stack:
	@echo "======================================================================"
	@echo "Run race test for queues/stack"
	cd $(LOCDIR)/queues/stack && $(DIR) $(GODEBUG) go test -cover -race ./

test-std:
	@echo "======================================================================"
	@echo "Run race test for queues/std"
	cd $(LOCDIR)/queues/std && $(DIR) $(GODEBUG) go test -cover -race ./

tests-priorityqueue:
	@echo "======================================================================"
	@echo "Run race test for queues/priorityqueue"
	cd $(LOCDIR)/queues/priorityqueue && $(DIR) $(GODEBUG) go test -cover -race ./

test: test-std test-stack tests-priorityqueue test-workers
	@echo "======================================================================"
	@echo "----"
	@echo "Run race test for ./item/..."
	cd $(LOCDIR)/item/ && $(DIR) $(GODEBUG) go test -cover -race ./
	@echo "----"
	@echo "Run race test for ./slavenode/..."
	cd $(LOCDIR)/slavenode/ && $(DIR) $(GODEBUG) go test -cover -race ./
	@echo "----"
	@echo "Run race test for ./faces/..."
	cd $(LOCDIR)/faces/ && $(DIR) $(GODEBUG) go test -cover -race ./
	@echo "----"
	@echo "Run race test for ./tracer/..."
	cd $(LOCDIR)/tracer/ && $(DIR) $(GODEBUG) go test -cover -race ./
	@echo "----"
	@echo "Run race test for ./"
	cd $(LOCDIR)/ && $(DIR) $(GODEBUG) go test -cover -race ./

lint:
	@echo "======================================================================"
	@echo "Run golint..."
	$(GOBIN)golint ./manager/...
	$(GOBIN)golint ./item/...
	$(GOBIN)golint ./faces/...
	$(GOBIN)golint ./status/...
	$(GOBIN)golint ./worker/...
	$(GOBIN)golint ./*.go

fmt:
	@echo "======================================================================"
	@echo "Run go fmt..."
	@go fmt ./manager/...
	@go fmt ./item/...
	@go fmt ./faces/...
	@go fmt ./status/...
	@go fmt ./worker/...s

mod:
	@echo "======================================================================"
	@echo "Run MOD"
	$(LOADENV) go mod verify
	$(LOADENV) go mod tidy
	$(LOADENV) go mod vendor
	$(LOADENV) go mod download
	$(LOADENV) go mod verify


clean-cache:
	@echo "clean-cache started..."
	go clean -cache
	go clean -testcache
	@echo "clean-cache complete!"
	@echo "clean-cache complete!"


mock-gen:
	$(LOADENV) ./bin/mockgen -package mmock github.com/iostrovok/conveyor/faces IConveyor > ./faces/mmock/iconveyor_mock.go
	$(LOADENV) ./bin/mockgen -package mmock github.com/iostrovok/conveyor/faces IChan > ./faces/mmock/ichan_mock.go
	$(LOADENV) ./bin/mockgen -package mmock github.com/iostrovok/conveyor/faces IInput> ./faces/mmock/iinput_mock.go
	$(LOADENV) ./bin/mockgen -package mmock github.com/iostrovok/conveyor/faces IItem > ./faces/mmock/iitem_mock.go
	$(LOADENV) ./bin/mockgen -package mmock github.com/iostrovok/conveyor/faces IManager > ./faces/mmock/imanager_mock.go
	$(LOADENV) ./bin/mockgen -package mmock github.com/iostrovok/conveyor/faces IWorker > ./faces/mmock/iworker_mock.go
	$(LOADENV) ./bin/mockgen -package mmock github.com/iostrovok/conveyor/faces IHandler > ./faces/mmock/ihandler_mock.go
	$(LOADENV) ./bin/mockgen -package mmock github.com/iostrovok/conveyor/faces ITestObject > ./faces/mmock/itestobject_mock.go

#################### Protobuf section

# Build gRPC files for GO
install: clean deps install-proto mod

install-proto:
	@mkdir -p $(GOBIN)
	wget https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTO_VERSION)/$(PROTO_ZIP)
	unzip -o $(PROTO_ZIP)
	rm $(PROTO_ZIP)
	$(LOADENV) go get -u google.golang.org/grpc
	$(LOADENV) go get -u github.com/golang/protobuf
	$(LOADENV) go get -u github.com/golang/protobuf/protoc-gen-go
	$(LOADENV) go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
	$(LOADENV) go build github.com/golang/protobuf/protoc-gen-go
	$(LOADENV) go build github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

# Build gRPC files for Go
build-go:
	@rm -rf ./protobuf/go
	@rm -rf ./tmp
	@mkdir -p ./tmp
	@mkdir -p ./protobuf/go/nodes/
	@echo "Generating client sources for Go..."
	@./bin/protoc --plugin=protoc-gen-go -I. --go_out=plugins=grpc:tmp protobuf/proto/*.proto
	@mv -f ./tmp/go/* ./protobuf/go/
	@rm -rf tmp


# Documentation
build-docs:
	@rm -rf ./protobuf/docs
	@mkdir -p ./protobuf/docs
	@./bin/protoc --plugin=protoc-gen-doc=./protoc-gen-doc --doc_out=./protobuf/docs/ --doc_opt=html,index.html \
		protobuf/proto/*.proto

build: build-go build-docs docs-gen

clean:
	@echo "clean-cache started..."
	rm -rf ./vendor
	go clean -cache
	go clean -testcache
	@echo "clean-cache complete!"

docs-gen:
	@echo "======================================================================"
	@echo 'MAKE: docs...'
	@mkdir -p ./docs
	@$(LOADENV) go run ./console/docs_generator.go


test-example:
	@echo "======================================================================"
	@echo "Run race test for workers"
	cd $(LOCDIR)/example/example-test-simple && $(DIR) $(GODEBUG) go test -cover -race ./