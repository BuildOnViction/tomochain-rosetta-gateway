// Copyright (c) 2020 TomoChain

package services

import (
	"context"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	tc "github.com/tomochain/tomochain-rosetta-gateway/tomochain-client"
	"github.com/tomochain/tomochain/common"
	"math/big"
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

	var (
		block *types.Block
		err error
	)
	if request.BlockIdentifier != nil {
		if request.BlockIdentifier.Hash != nil {
			block, err = s.client.GetBlockByHash(ctx, common.HexToHash(*(request.BlockIdentifier.Hash)))
		} else if request.BlockIdentifier.Index != nil {
			block, err = s.client.GetBlockByNumber(ctx, big.NewInt(*(request.BlockIdentifier.Index)))
		}
	}
	if err != nil || block == nil {
		return nil, ErrUnableToGetBlk
	}


	resp := &types.BlockResponse{
		Block: block,
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
