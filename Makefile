
# Copyright (c) 2020 TomoChain

export GO111MODULE=on

# Go parameters
GOCMD=go
GOLINT=golint
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GO_FILES := $(shell find $(shell go list -f '{{.Dir}}' $(GO_PACKAGES)) -name \*.go)
BUILD_TARGET=tomochain-rosetta-gateway
BIN_DIRECTORY=bin

default: build
all: clean build

build:
	$(GOBUILD) -o ./$(BIN_DIRECTORY)/$(BUILD_TARGET)

clean:
	@echo "Cleaning..."
	@rm -rf ./$(BIN_DIRECTORY)

update-tracer:
	curl https://raw.githubusercontent.com/tomochain/tomochain/master/eth/tracers/internal/tracers/call_tracer.js -o tomochain-client/call_tracer.js

gofmt:
	$(GOFMT) -s -w $(GO_FILES)