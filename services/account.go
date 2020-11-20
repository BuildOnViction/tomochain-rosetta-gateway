// Copyright (c) 2020 TomoChain

package services

import (
	"context"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/tomochain/tomochain-rosetta-gateway/common"
	"github.com/tomochain/tomochain-rosetta-gateway/config"
	tc "github.com/tomochain/tomochain-rosetta-gateway/tomochain-client"
)

type accountAPIService struct {
	client tc.TomoChainClient
}

// NewAccountAPIService creates a new instance of an AccountAPIService.
func NewAccountAPIService(client tc.TomoChainClient) server.AccountAPIServicer {
	return &accountAPIService{
		client: client,
	}
}

// AccountBalance implements the /account/balance endpoint.
func (s *accountAPIService) AccountBalance(
	ctx context.Context,
	request *types.AccountBalanceRequest,
) (*types.AccountBalanceResponse, *types.Error) {
	if s.client.GetConfig().Server.Mode != config.ServerModeOnline {
		return nil, common.ErrUnavailableOffline
	}
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}
	resp, err := s.client.GetAccount(ctx, request.AccountIdentifier.Address)
	if err != nil {
		return nil, common.ErrUnableToGetAccount
	}
	return resp, nil
}
