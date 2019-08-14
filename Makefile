all: check_go_fmt deps proto lint test build build-client

ifneq ($(OS),Windows_NT)
    OS := $(shell sh -c 'uname -s 2>/dev/null')
endif

ifeq ($(OS),Linux)
    LD_FLAGS = -ldflags="-s -w"
endif

.PHONY: build
build:
	GO111MODULE=on CGO_ENABLED=0 go build $(LD_FLAGS) -o bin/service

.PHONY: build-client
build-client:
	GO111MODULE=on CGO_ENABLED=0 go build $(LD_FLAGS) -o ./bin/client ./client/

.PHONY: docker
docker:
	docker build --no-cache -t indiependente/gw-example-server .

.PHONY: docker-client
docker-client:
	docker build --no-cache -t indiependente/gw-example-client . -f client/Dockerfile

.PHONY: docker_clean
docker_clean:
	docker rm indiependente/gw-example-server | true
	docker rm indiependente/gw-example-client | true

.PHONY: deps-init
deps-init:
	rm -f go.mod go.sum
	@GO111MODULE=on go mod init
	@GO111MODULE=on go mod tidy

.PHONY: update-deps
update-deps:
	@GO111MODULE=on go mod tidy

.PHONY: deps
deps:
	@GO111MODULE=on go mod download

.PHONY: lint
lint:
	command -v golangci-lint || (cd /usr/local ; wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest)
	GO111MODULE=on golangci-lint run --disable-all \
	--deadline=10m \
	--skip-files \.*_mock\.*\.go \
	-E errcheck \
	-E govet \
	-E unused \
	-E gocyclo \
	-E golint \
	-E varcheck \
	-E structcheck \
	-E maligned \
	-E ineffassign \
	-E interfacer \
	-E unconvert \
	-E goconst \
	-E gosimple \
	-E staticcheck \
	-E gosec

.PHONY: proto
proto:
	@protoc -I/usr/local/include -I. \
  -I$$GOPATH/src \
  -I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:./rpc/ proto/service.proto

	@protoc -I/usr/local/include -I. \
  -I$$GOPATH/src \
  -I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:./rpc proto/service.proto

	@protoc -I/usr/local/include -I. \
  -I$$GOPATH/src \
  -I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --swagger_out=logtostderr=true:./swagger \
  proto/service.proto
	@mv rpc/proto/*.go rpc/
	@rm -rf rpc/proto
	@mv swagger/proto/service.swagger.json swagger/service.swagger.json
	@rm -rf swagger/proto

.PHONY: test
test:
	GO111MODULE=on go test -v -cover -race -tags=unit ./...

.PHONY: check_go_fmt
check_go_fmt:
	@if [ -n "$$(gofmt -d $$(find . -name '*.go'))" ]; then \
		>&2 echo "The .go sources aren't formatted. Please format them with 'go fmt'."; \
		exit 1; \
	fi
