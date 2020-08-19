// Copyright (c) 2020 TomoChain

package services

import (
	"context"
	tc "github.com/tomochain/tomochain-rosetta-gateway/tomochain-client"
	"math/big"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type blockAPIService struct {
	client tc.TomoChainClient
}

// NewBlockAPIService creates a new instance of an AccountAPIService.
func NewBlockAPIService(client tc.TomoChainClient) server.BlockAPIServicer {
	return &blockAPIService{
		client: client,
	}
}

func (s *blockAPIService) Block(
	ctx context.Context,
	request *types.BlockRequest,
) (*types.BlockResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}
	var blockNumber int64

	if request.BlockIdentifier != nil {
		if request.BlockIdentifier.Index != nil {
			blockNumber = *request.BlockIdentifier.Index
		} else if request.BlockIdentifier.Hash != nil {
			return nil, ErrMustQueryByIndex
		}
	}

	tblk, err := s.client.GetBlock(ctx, big.NewInt(blockNumber))
	if err != nil {
		return nil, ErrUnableToGetBlk
	}

	resp := &types.BlockResponse{
		Block: tblk,
	}

	return resp, nil
}

// BlockTransaction implements the /block/transaction endpoint.
// Note: we don't implement this, since we already return all transactions
// in the /block endpoint reponse above.
func (s *blockAPIService) BlockTransaction(
	ctx context.Context,
	request *types.BlockTransactionRequest,
) (*types.BlockTransactionResponse, *types.Error) {
	return nil, ErrNotImplemented
}
