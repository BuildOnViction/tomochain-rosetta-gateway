FROM golang:1.13-alpine as build
RUN apk add --no-cache make gcc musl-dev linux-headers git curl
ENV GO111MODULE=on

# download tomochain core source
ARG TOMOCHAIN_CORE_VERSION="v2.2.5"
RUN git clone --tag $TOMOCHAIN_CORE_VERSION https://github.com/tomochain/tomochain /usr/local/tomochain-source && \
cd /usr/local/tomochain-source && \
make clean && make tomo && chmod +x /usr/local/tomochain-source/build/bin/tomo && \
cp ./build/bin/tomo /usr/local/bin && \
cd /usr/local && rm -rf tomochain-source 

# download genesis
ARG GENESIS_URL="https://raw.githubusercontent.com/tomochain/tomochain/master/genesis/mainnet.json"
ARG GENESIS_PATH="~/genesis.json"
RUN curl -L $GENESIS_URL -o $GENESIS_PATH

# download tomochain-rosetta-gateway source
ARG TOMOCHAIN_ROSETTA_GATEWAY_VERSION="0.0.2"
RUN git clone --tag $TOMOCHAIN_ROSETTA_GATEWAY_VERSION https://github.com/tomochain/tomochain-rosetta-gateway /usr/local/tomochain-rosetta-gateway-source && \
cd /usr/local/tomochain-rosetta-gateway-source && \
make clean && make all && \
cp ./build/bin /usr/local/bin && \
cd /usr/local && rm -rf tomochain-rosetta-gateway-source

FROM alpine:latest
RUN apk add --no-cache ca-certificates

# https://www.rosetta-api.org/docs/standard_storage_location.html
VOLUME /data
WORKDIR /app

ARG P2P_PORT=30303
ARG CHAIN_ID=88
ARG DEFAULT_GAS_PRICE=2500
ARG DEFAULT_GAS_LIMIT=84000000
ARG LOG_LEVEL=3
ARG NODE_NAME="tomochain-rosetta"
ARG DEFAULT_BOOTNODES="enode://97f0ca95a653e3c44d5df2674e19e9324ea4bf4d47a46b1d8560f3ed4ea328f725acec3fcfcb37eb11706cf07da669e9688b091f1543f89b2425700a68bc8876@3.212.20.0:30301,enode://b72927f349f3a27b789d0ca615ffe3526f361665b496c80e7cc19dace78bd94785fdadc270054ab727dbb172d9e3113694600dd31b2558dd77ad85a869032dea@188.166.207.189:30301,enode://c8f2f0643527d4efffb8cb10ef9b6da4310c5ac9f2e988a7f85363e81d42f1793f64a9aa127dbaff56b1e8011f90fe9ff57fa02a36f73220da5ff81d8b8df351@104.248.98.60:30301"
ARG DATA_DIR=/app/tomochain-data
RUN cd /app && mkdir tomochain-data && cd tomochain-data && mkdir tomox


# Init TomoChain Client
RUN tomo init $GENESIS_PATH --datadir $DATA_DIR

# Run TomoChain node
RUN tomo  \
--gcmode "archive" \
--announce-txs   \
--store-reward  \
--tomox.dbengine "leveldb" \
--verbosity $LOG_LEVEL   \
--datadir $DATA_DIR   \
--tomox.datadir $DATA_DIR/tomox \
--keystore $KEYSTORE   \
--identity $NODE_NAME   \
--password $PASSWORD   \
--networkid $CHAIN_ID   \
--port $P2P_PORT   \
--gasprice $DEFAULT_GAS_PRICE \
--unlock $COINBASE \
--bootnodes $DEFAULT_BOOTNODES \
--syncmode full \
--targetgaslimit $DEFAULT_GAS_LIMIT \

# expose TomoChain client P2P ports
EXPOSE $P2P_PORT/tcp
EXPOSE $P2P_PORT/udp

