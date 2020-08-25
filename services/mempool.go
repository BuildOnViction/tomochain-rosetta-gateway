// Copyright (c) 2020 TomoChain

package services

import (
	"context"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/tomochain/tomochain-rosetta-gateway/common"
	tc "github.com/tomochain/tomochain-rosetta-gateway/tomochain-client"
	tomochaincommon "github.com/tomochain/tomochain/common"
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
	terr := common.ValidateNetworkIdentifier(ctx, m.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}

	pending, err := m.client.GetMempool(ctx)
	if err != nil {
		return nil, common.ErrUnableToGetTxns
	}
	res := &types.MempoolResponse{}
	tis := []*types.TransactionIdentifier{}
	for _, hash := range pending {
		tis = append(tis, &types.TransactionIdentifier{Hash: hash.String()})
	}
	res.TransactionIdentifiers = tis
	return res, nil
}

// Get a transaction in the mempool by its Transaction Identifier
func (m *mempoolAPIService) MempoolTransaction(
	ctx context.Context,
	request *types.MempoolTransactionRequest,
) (*types.MempoolTransactionResponse, *types.Error) {
	terr := common.ValidateNetworkIdentifier(ctx, m.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}

	tx, err := m.client.GetMempoolTransaction(ctx, tomochaincommon.HexToHash(request.TransactionIdentifier.Hash))
	if err != nil {
		return nil, common.ErrUnableToGetTxns
	}
	return &types.MempoolTransactionResponse{
		Transaction: tx,
	}, nil
}
