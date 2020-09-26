GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=photolum

all: build

build:
	export GO111MODULE=on && \
	export GOOS=linux && \
	export GOARCH=amd64 && \
	export CGO_ENABLED=1 && \
	$(GOBUILD) -o $(BINARY_NAME) cmd/main.go