# Not Support Windows

.PHONY: init gentool api openapi ts wire build docker all help

ifeq ($(OS),Windows_NT)
    IS_WINDOWS:=1
endif

CURRENT_DIR := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

SRCS_MK := $(foreach dir, app, $(wildcard $(dir)/*/*/Makefile))

# initialize develop environment
init: plugin cli

# install protoc plugin
plugin:
	# go
	@go install go.newcapec.cn/ncttools/nmskit/cmd/protoc-gen-go-faster@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install go.newcapec.cn/ncttools/nmskit/cmd/protoc-gen-go-grpc@latest
	@go install go.newcapec.cn/ncttools/nmskit/cmd/protoc-gen-go-grpc-impl@latest
	@go install go.newcapec.cn/ncttools/nmskit/cmd/protoc-gen-go-http/v2@latest
	@go install go.newcapec.cn/ncttools/nmskit/cmd/protoc-gen-go-errors/v2@latest
	@go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	@go install github.com/envoyproxy/protoc-gen-validate@latest
	# Typescript
	@npm install -g ts-proto

# install cli tools
cli:
	@go install go.newcapec.cn/ncttools/nmskit/cmd/NmsKit@latest
	@go install github.com/google/gnostic@latest
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@go install gorm.io/gen/tools/gentool@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# generate gorm code
gentool:
	@gentool -c gorm_gen.yaml

# generate protobuf api go code
api:
	@cd api && \
	buf generate

# generate OpenAPI v3 docs.
openapi:
	@cd api && \
	buf generate --template buf.openapi.gen.yaml

# generate typescript.
ts:
	@cd api && \
	buf generate --template buf.typescript.gen.yaml

# generate wire code
wire:
	@go mod tidy
	@cd cmd/start && \
	wire

# build all service applications
build:
	go mod tidy
	go env -w GOOS=linux
	go build -o ShopService ./cmd
	go env -w GOOS=darwin

docker:
	docker build -t shop-service .
	docker save -o ./shop-service.tar shop-service:latest

# generate & build all service applications
all:
	$(foreach dir, $(dir $(realpath $(SRCS_MK))),\
      cd $(dir);\
      make app;\
    )

# show help
help:
	@echo ""
	@echo "Usage:"
	@echo " make [target]"
	@echo ""
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
