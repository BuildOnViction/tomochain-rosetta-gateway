
# Copyright (c) 2020 TomoChain

export GO111MODULE=on

# Go parameters
GOCMD=go
GOLINT=golint
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BIN_DIRECTORY=bin

.PHONY: deps build run lint run-mainnet-online run-mainnet-offline run-testnet-online \
	run-testnet-offline check-comments \
	build-local fmt update-tracer \
	update-bootstrap-balances

default: build-local

build:
	docker build -t tomochain-rosetta:latest https://github.com/tomochain/tomochain-rosetta-gateway.git

build-local:
	docker build -t tomochain-rosetta:latest .

build-release:
	# make sure to always set version with vX.X.X
	docker build -t tomochain-rosetta:$(version) .;
	docker save tomochain-rosetta:$(version) | gzip > tomochain-rosetta-$(version).tar.gz;

run-mainnet-online:
	cp tomochain/mainnet.toml tomochain/tomochain.toml && docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -v "${PWD}/tomochain-data:/data" -e "MODE=ONLINE" -e "NETWORK=MAINNET" -e "PORT=8080" -p 8080:8080 -p 30303:30303 tomochain-rosetta:latest

run-mainnet-offline:
	cp tomochain/mainnet.toml tomochain/tomochain.toml && docker run -d --rm -e "MODE=OFFLINE" -e "NETWORK=MAINNET" -e "PORT=8081" -p 8081:8081 tomochain-rosetta:latest

run-mainnet-remote:
	cp tomochain/mainnet.toml tomochain/tomochain.toml && docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -e "MODE=ONLINE" -e "NETWORK=MAINNET" -e "PORT=8080" -e "TOMO=$(tomo)" -p 8080:8080 -p 30303:30303 tomochain-rosetta:latest

run-testnet-online:
	cp tomochain/testnet.toml tomochain/tomochain.toml && docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -v "${PWD}/tomochain-data:/data" -e "MODE=ONLINE" -e "NETWORK=TESTNET" -e "PORT=8080" -p 8080:8080 -p 30303:30303 tomochain-rosetta:latest

run-testnet-offline:
	docker run -d --rm -e "MODE=OFFLINE" -e "NETWORK=TESTNET" -e "PORT=8081" -p 8081:8081 tomochain-rosetta:latest

run-testnet-remote:
	cp tomochain/testnet.toml tomochain/tomochain.toml && docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -e "MODE=ONLINE" -e "NETWORK=TESTNET" -e "PORT=8080" -e "TOMO=$(tomo)" -p 8080:8080 -p 30303:30303 tomochain-rosetta:latest

run-devnet-online:
	cp tomochain/devnet.toml tomochain/tomochain.toml && docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -v "${PWD}/tomochain-data:/data" -e "MODE=ONLINE" -e "NETWORK=DEVNET" -e "PORT=8080" -p 8080:8080 -p 30303:30303 tomochain-rosetta:latest

run-devnet-offline:
	cp tomochain/devnet.toml tomochain/tomochain.toml && docker run -d --rm -e "MODE=OFFLINE" -e "NETWORK=DEVNET" -e "PORT=8081" -p 8081:8081 tomochain-rosetta:latest

run-devnet-remote:
	cp tomochain/devnet.toml tomochain/tomochain.toml && docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -e "MODE=ONLINE" -e "NETWORK=DEVNET" -e "PORT=8080" -e "TOMO=$(tomo)" -p 8080:8080 -p 30303:30303 tomochain-rosetta:latest


check-comments:
	${GOLINT_CMD} -set_exit_status ${GO_FOLDERS} .

lint: | check-comments
	golangci-lint run --timeout 2m0s -v -E ${LINT_SETTINGS},gomnd


clean:
	@echo "Cleaning..."
	@rm -rf ./$(BIN_DIRECTORY)

deps:
	go get ./...

update-tracer:
	curl https://raw.githubusercontent.com/tomochain/tomochain/master/eth/tracers/internal/tracers/call_tracer.js -o tomochain/call_tracer.js
update-bootstrap-balances:
	go run main.go utils:generate-bootstrap tomochain/genesis_files/mainnet.json rosetta-cli-conf/mainnet/bootstrap_balances.json;
	go run main.go utils:generate-bootstrap tomochain/genesis_files/testnet.json rosetta-cli-conf/testnet/bootstrap_balances.json;
gofmt:
	$(GOFMT) -s -w $(GO_FILES)
