// Copyright (c) 2020 TomoChain

package services

import (
	"context"
	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/tomochain/tomochain-rosetta-gateway/common"
	"github.com/tomochain/tomochain-rosetta-gateway/configuration"
	"github.com/tomochain/tomochain-rosetta-gateway/tomochain"
	"github.com/tomochain/tomochain/params"
	"strconv"
)

// NetworkAPIService implements the server.NetworkAPIServicer interface.
type NetworkAPIService struct {
	config *configuration.Configuration
	client Client
}

// NewNetworkAPIService creates a new instance of a NetworkAPIService.
func NewNetworkAPIService(
	cfg *configuration.Configuration,
	client Client,
) *NetworkAPIService {
	return &NetworkAPIService{
		config: cfg,
		client: client,
	}
}

// NetworkList implements the /network/list endpoint.
func (s *NetworkAPIService) NetworkList(
	ctx context.Context,
	request *types.MetadataRequest,
) (*types.NetworkListResponse, *types.Error) {
	return &types.NetworkListResponse{
		NetworkIdentifiers: []*types.NetworkIdentifier{
			{
				Blockchain: common.TomoChainBlockchain,
				Network:    strconv.FormatUint(common.TomoChainMainnetNetWorkId, 10),
			},
			{
				Blockchain: common.TomoChainBlockchain,
				Network:    strconv.FormatUint(common.TomoChainTestnetNetWorkId, 10),
			},
			{
				Blockchain: common.TomoChainBlockchain,
				Network:    strconv.FormatUint(common.TomoChainDevnetNetWorkId, 10),
			},
		},
	}, nil
}

// NetworkOptions implements the /network/options endpoint.
func (s *NetworkAPIService) NetworkOptions(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkOptionsResponse, *types.Error) {
	return &types.NetworkOptionsResponse{
		Version: &types.Version{
			RosettaVersion: types.RosettaAPIVersion,
			MiddlewareVersion: &configuration.MiddlewareVersion,
			NodeVersion:    params.Version,
		},
		Allow: &types.Allow{
			OperationStatuses: []*types.OperationStatus{
				{
					Status:     common.SUCCESS,
					Successful: true,
				},
				{
					Status:     common.FAIL,
					Successful: false,
				},
			},
			OperationTypes:          common.SupportedOperationTypes(),
			Errors:                  common.ErrorList,
			HistoricalBalanceLookup: common.HistoricalBalanceSupported,
			CallMethods:             tomochain.CallMethods,
		},
	}, nil
}

// NetworkStatus implements the /network/status endpoint.
func (s *NetworkAPIService) NetworkStatus(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkStatusResponse, *types.Error) {
	if s.config.Mode != configuration.Online {
		return nil, common.ErrUnavailableOffline
	}

	currentBlock, currentTime, syncStatus, peers, err := s.client.Status(ctx)
	if err != nil {
		return nil, common.ErrTomo
	}

	if currentTime < asserter.MinUnixEpoch {
		return nil, common.ErrTomoNotReady
	}

	return &types.NetworkStatusResponse{
		CurrentBlockIdentifier: currentBlock,
		CurrentBlockTimestamp:  currentTime,
		GenesisBlockIdentifier: s.config.GenesisBlockIdentifier,
		SyncStatus:             syncStatus,
		Peers:                  peers,
	}, nil
}

// ValidateNetworkIdentifier validates the network identifier.
func ValidateNetworkIdentifier(ctx context.Context, client Client, ni *types.NetworkIdentifier) *types.Error {
	if ni != nil {
		if ni.Blockchain != common.TomoChainBlockchain {
			return common.ErrInvalidBlockchain
		}
		if ni.SubNetworkIdentifier != nil {
			return common.ErrInvalidSubnetwork
		}
		id, err := strconv.Atoi(ni.Network)
		if err != nil {
			return common.ErrInvalidNetwork
		}
		if chainId, err := client.GetChainID(ctx); err != nil || chainId.Uint64() != uint64(id) {
			return common.ErrInvalidNetwork
		}
	} else {
		return common.ErrMissingNID
	}
	return nil
}
