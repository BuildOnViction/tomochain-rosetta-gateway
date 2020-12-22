// Copyright (c) 2020 TomoChain

package services

import (
	"context"
	"fmt"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/tomochain/tomochain-rosetta-gateway/common"
	"github.com/tomochain/tomochain-rosetta-gateway/config"
	tc "github.com/tomochain/tomochain-rosetta-gateway/tomochain-client"
	"github.com/tomochain/tomochain/params"
	"strconv"
)

type networkAPIService struct {
	client tc.TomoChainClient
}

// NewNetworkAPIService creates a new instance of a NetworkAPIService.
func NewNetworkAPIService(client tc.TomoChainClient) server.NetworkAPIServicer {
	return &networkAPIService{
		client: client,
	}
}

// NetworkList implements the /network/list endpoint.
func (s *networkAPIService) NetworkList(
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

// NetworkStatus implements the /network/status endpoint.
func (s *networkAPIService) NetworkStatus(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkStatusResponse, *types.Error) {
	if s.client.GetConfig().Server.Mode != config.ServerModeOnline {
		return nil, common.ErrUnavailableOffline
	}
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		fmt.Println(terr)
		return nil, terr
	}

	blk, err := s.client.GetBlockByNumber(ctx, nil) // nil means: get latest block
	if err != nil {
		fmt.Println(err)
		return nil, common.ErrUnableToGetNodeStatus
	}
	genesisblk, err := s.client.GetGenesisBlock(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, common.ErrUnableToGetNodeStatus
	}

	resp := &types.NetworkStatusResponse{
		CurrentBlockIdentifier: blk.BlockIdentifier,
		CurrentBlockTimestamp:  blk.Timestamp, // ms
		GenesisBlockIdentifier: genesisblk.BlockIdentifier,
	}

	return resp, nil
}

// NetworkOptions implements the /network/options endpoint.
func (s *networkAPIService) NetworkOptions(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkOptionsResponse, *types.Error) {
	return &types.NetworkOptionsResponse{
		Version: &types.Version{
			RosettaVersion: s.client.GetConfig().Server.RosettaVersion,
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
		},
	}, nil
}

// ValidateNetworkIdentifier validates the network identifier.
func ValidateNetworkIdentifier(ctx context.Context, client tc.TomoChainClient, ni *types.NetworkIdentifier) *types.Error {
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
