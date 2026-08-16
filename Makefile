BINARY=hn
BUILD_DIR=bin

GIT_SHA=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION=$(shell go version | awk '{print $$3}')

LDFLAGS=-ldflags "\
				-X github.com/sanverite/hn/internal/version.GitSHA=$(GIT_SHA) \
				-X github.com/sanverite/hn/internal/version.BuildTime=$(BUILD_TIME) \
				-X github.com/sanverite/hn/internal/version.GoVersion=$(GO_VERSION)"

.PHONY: build test lint race clean

build:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) ./cmd/hn

test:
	go test ./...

race:
	go test -race ./...

lint:
	golangci-lint run

clean:
	rm -rf $(BUILD_DIR)

