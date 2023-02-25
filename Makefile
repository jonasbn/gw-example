all: check_go_fmt deps proto lint test build build-client

ifneq ($(OS),Windows_NT)
    OS := $(shell sh -c 'uname -s 2>/dev/null')
endif

ifeq ($(OS),Linux)
    LD_FLAGS = -ldflags="-s -w"
endif

.PHONY: build
build:
	CGO_ENABLED=0 go build $(LD_FLAGS) -o bin/service

.PHONY: build-client
build-client:
	CGO_ENABLED=0 go build $(LD_FLAGS) -o ./bin/client ./client/

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
	@go mod init
	@go mod tidy

.PHONY: update-deps
update-deps:
	@go mod tidy

.PHONY: deps
deps:
	@go mod download

.PHONY: lint
lint: golint buflint

.PHONY: golint
golint:
	command -v golangci-lint || (cd /usr/local ; wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest)
	golangci-lint run --disable-all \
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

.PHONY: buflint
buflint:
	command -v buf || go install github.com/bufbuild/buf/cmd/buf@latest
	buf lint

.PHONY: test
test:
	go test -v -cover -race -tags=unit ./...

.PHONY: check_go_fmt
check_go_fmt:
	@if [ -n "$$(gofmt -d $$(find . -name '*.go'))" ]; then \
		>&2 echo "The .go sources aren't formatted. Please format them with 'go fmt'."; \
		exit 1; \
	fi

generate: generate/proto
generate/proto:
	command -v buf || go install github.com/bufbuild/buf/cmd/buf@latest
	buf generate
