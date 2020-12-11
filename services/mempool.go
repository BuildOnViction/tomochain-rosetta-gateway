// Copyright (c) 2020 TomoChain

package services

import (
	"context"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/tomochain/tomochain-rosetta-gateway/common"
	tc "github.com/tomochain/tomochain-rosetta-gateway/tomochain-client"
)

type mempoolAPIService struct {
	client tc.TomoChainClient
}

// NewMempoolAPIService creates a new instance of an MempoolAPIService.
func NewMempoolAPIService(client tc.TomoChainClient) server.MempoolAPIServicer {
	return &mempoolAPIService{
		client: client,
	}
}

// Get all Transaction Identifiers in the mempool
func (m *mempoolAPIService) Mempool(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.MempoolResponse, *types.Error) {
	return nil, common.ErrNotImplemented
}

// Get a transaction in the mempool by its Transaction Identifier
func (m *mempoolAPIService) MempoolTransaction(
	ctx context.Context,
	request *types.MempoolTransactionRequest,
) (*types.MempoolTransactionResponse, *types.Error) {
	return nil, common.ErrNotImplemented
}
