// Copyright (c) 2020 TomoChain

package services

import (
	"context"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/tomochain/tomochain-rosetta-gateway/common"
	"github.com/tomochain/tomochain-rosetta-gateway/configuration"
)

// AccountAPIService implements the server.AccountAPIServicer interface.
type AccountAPIService struct {
	config *configuration.Configuration
	client Client
}
// NewAccountAPIService returns a new *AccountAPIService.
func NewAccountAPIService(
	cfg *configuration.Configuration,
	client Client,
) *AccountAPIService {
	return &AccountAPIService{
		config: cfg,
		client: client,
	}
}

// AccountBalance implements the /account/balance endpoint.
func (s *AccountAPIService) AccountBalance(
	ctx context.Context,
	request *types.AccountBalanceRequest,
) (*types.AccountBalanceResponse, *types.Error) {
	if s.config.Mode != configuration.Online {
		return nil, common.ErrUnavailableOffline
	}
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}
	resp, err := s.client.Balance(ctx, request.AccountIdentifier, request.BlockIdentifier)
	if err != nil {
		return nil, common.ErrUnableToGetAccount
	}
	return resp, nil
}

// AccountCoins implements /account/coins.
func (s *AccountAPIService) AccountCoins(context.Context, *types.AccountCoinsRequest) (*types.AccountCoinsResponse, *types.Error) {
	// TomoChain blockchain doesn't support coin identifier
	return nil, common.ErrNotImplemented
}