GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
GOLINT=golint
BINARY_NAME=photolum

all: build

build:
	export GO111MODULE=on && \
	export GOOS=linux && \
	export GOARCH=amd64 && \
	export CGO_ENABLED=0 && \
	$(GOBUILD) -o $(BINARY_NAME) cmd/main.go
vet:
	# go vet
	$(GOVET) cmd
lint:
	# golint
	$(GOLINT) ./...
perfect: lint vet