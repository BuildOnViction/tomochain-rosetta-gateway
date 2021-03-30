# Compile golang
FROM ubuntu:18.04 as golang-builder

RUN mkdir -p /app \
  && chown -R nobody:nogroup /app
WORKDIR /app

RUN apt-get update && apt-get install -y curl make gcc g++ git
ENV GOLANG_VERSION 1.15.5
ENV GOLANG_DOWNLOAD_SHA256 9a58494e8da722c3aef248c9227b0e9c528c7318309827780f16220998180a0d
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
  && echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
  && tar -C /usr/local -xzf golang.tar.gz \
  && rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
ENV GO111MODULE=on

# Compile TomoChain Client software
FROM golang-builder as tomochain-builder

ARG TOMOCHAIN_CORE_VERSION="fix-new-evm"
RUN rm -rf tomochain-source
RUN git clone --branch $TOMOCHAIN_CORE_VERSION https://github.com/tomochain/tomochain.git tomochain-source
RUN cd tomochain-source && \
make clean && make tomo && chmod +x ./build/bin/tomo && \
mv ./build/bin/tomo /app/tomo && \
cp ./genesis/mainnet.json /app/genesis.json && \
cd .. && rm -rf tomochain-source


# Compile tomochain-rosetta
FROM golang-builder as rosetta-builder

# Use native remote build context to build in any directory
ARG TOMOCHAIN_ROSETTA_GATEWAY_VERSION="master"
RUN mkdir /app/tomochain
RUN cd /app
RUN rm -rf tomochain-rosetta-gateway-source
RUN git clone --branch $TOMOCHAIN_ROSETTA_GATEWAY_VERSION https://github.com/tomochain/tomochain-rosetta-gateway.git tomochain-rosetta-gateway-source
RUN cd tomochain-rosetta-gateway-source && \
go build -o tomochain-rosetta . && \
mv ./tomochain-rosetta /app/tomochain-rosetta && \
mv ./tomochain/call_tracer.js /app/tomochain/call_tracer.js && \
mv ./tomochain/tomochain.toml /app/tomochain/tomochain.toml && \
cd .. && rm -rf tomochain-rosetta-gateway-source



## Build Final Image
FROM ubuntu:18.04

RUN mkdir -p /app \
  && chown -R nobody:nogroup /app \
  && mkdir -p /data \
  && chown -R nobody:nogroup /data

WORKDIR /app

# Copy binary from tomochain-builder
COPY --from=tomochain-builder /app/tomo /app/tomo
# Copy genesis from tomochain-builder
COPY --from=tomochain-builder /app/genesis.json /app/genesis.json

# Copy binary from rosetta-builder
COPY --from=rosetta-builder /app/tomochain /app/tomochain
COPY --from=rosetta-builder /app/tomochain-rosetta /app/tomochain-rosetta

# Set permissions for everything added to /app
RUN chmod -R 755 /app/*

CMD ["/app/tomochain-rosetta", "run"]

